package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-slack-message/lib/slack"
	"github.com/bitrise-steplib/steps-slack-message/lib/step"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// success is true if the build is successful, false otherwise.
var success = os.Getenv("BITRISE_BUILD_STATUS") == "0"
var logger = log.NewLogger()

// selectValue chooses the right value based on the result of the build.
func selectValue(ifSuccess, ifFailed string) string {
	if success || ifFailed == "" {
		return ifSuccess
	}
	return ifFailed
}

// selectBool chooses the right boolean value based on the result of the build.
func selectBool(ifSuccess, ifFailed bool) bool {
	if success {
		return ifSuccess
	}
	return ifFailed
}

// ensureNewlines replaces all \n substrings with newline characters.
func ensureNewlines(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func newMessage(c step.Config) slack.Message {
	msg := slack.Message{
		Channel: strings.TrimSpace(selectValue(c.Channel, c.ChannelOnError)),
		Text:    selectValue(c.Text, c.TextOnError),
		Attachments: []slack.Attachment{{
			Fallback:   ensureNewlines(selectValue(c.Message, c.MessageOnError)),
			Color:      selectValue(c.Color, c.ColorOnError),
			PreText:    selectValue(c.PreText, c.PreTextOnError),
			AuthorName: c.AuthorName,
			Title:      selectValue(c.Title, c.TitleOnError),
			TitleLink:  c.TitleLink,
			Text:       ensureNewlines(selectValue(c.Message, c.MessageOnError)),
			Fields:     parseFields(c.Fields),
			ImageURL:   selectValue(c.ImageURL, c.ImageURLOnError),
			ThumbURL:   selectValue(c.ThumbURL, c.ThumbURLOnError),
			Footer:     c.Footer,
			FooterIcon: c.FooterIcon,
			Buttons:    parseButtons(c.Buttons),
		}},
		IconEmoji:      selectValue(c.IconEmoji, c.IconEmojiOnError),
		IconURL:        selectValue(c.IconURL, c.IconURLOnError),
		LinkNames:      c.LinkNames,
		Username:       selectValue(c.Username, c.UsernameOnError),
		ThreadTs:       selectValue(c.ThreadTs, c.ThreadTsOnError),
		ReplyBroadcast: selectBool(c.ReplyBroadcast, c.ReplyBroadcastOnError),
	}
	if c.TimeStamp {
		msg.Attachments[0].TimeStamp = int(time.Now().Unix())
	}
	return msg
}

func parseFields(s string) (fs []slack.Field) {
	for _, p := range pairs(s) {
		fs = append(fs, slack.Field{Title: p[0], Value: p[1]})
	}
	return
}

func parseButtons(s string) (bs []slack.Button) {
	for _, p := range pairs(s) {
		bs = append(bs, slack.Button{Text: p[0], URL: p[1]})
	}
	return
}

// pairs slices every lines in s into two substrings separated by the first pipe
// character and returns a slice of those pairs.
func pairs(s string) [][2]string {
	var ps [][2]string
	for _, line := range strings.Split(s, "\n") {
		a := strings.SplitN(line, "|", 2)
		if len(a) == 2 && a[0] != "" && a[1] != "" {
			ps = append(ps, [2]string{a[0], a[1]})
		}
	}
	return ps
}

// postMessage sends a message to a channel.
func postMessage(conf step.Config, msg slack.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	logger.Debugf("Request to Slack: %s\n", b)

	url := strings.TrimSpace(selectValue(string(conf.WebhookURL), string(conf.WebhookURLOnError)))
	if url == "" {
		url = "https://slack.com/api/chat.postMessage"
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if string(conf.APIToken) != "" {
		req.Header.Add("Authorization", "Bearer "+string(conf.APIToken))
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("server error: %s, failed to read response: %s", resp.Status, err)
		}
		return fmt.Errorf("server error: %s, response: %s", resp.Status, body)
	}

	if err := exportOutputs(&conf, resp); err != nil {
		return fmt.Errorf("failed to export outputs: %s", err)
	}

	return nil
}

func validate(conf *step.Config) error {
	if conf.APIToken == "" && conf.WebhookURL == "" {
		return fmt.Errorf("Both API Token and WebhookURL are empty. You need to provide one of them. If you want to use incoming webhooks provide the webhook url. If you want to use a bot to send a message provide the bot API token")
	}

	if conf.APIToken != "" && conf.WebhookURL != "" {
		logger.Warnf("Both API Token and WebhookURL are provided. Using the API Token")
		conf.WebhookURL = ""

	}
	return nil
}

type Stage func() error

func run(stages ...Stage) error {
	for _, stage := range stages {
		if err := stage(); err != nil {
			return err
		}
	}
	return nil
}

func parseConfig(conf *step.Config, repo env.Repository) error {
	if err := stepconf.NewInputParser(repo).Parse(&conf); err != nil {
		return err
	}
	stepconf.Print(conf)
	return nil
}

func enableDebugLog(conf *step.Config) error {
	logger.EnableDebugLog(conf.Debug)
	return nil
}

func createMessage(conf *step.Config, msg *slack.Message) error {
	*msg = newMessage(*conf)
	return nil
}

func main() {
	var conf step.Config
	var msg slack.Message
	envRepo := env.NewRepository()

	err := run(
		func() error { return parseConfig(&conf, envRepo) },
		func() error { return enableDebugLog(&conf) },
		func() error { return validate(&conf) },
		func() error { return createMessage(&conf, &msg) },
		func() error { return postMessage(conf, msg) },
	)
	if err != nil {
		logger.Errorf("Error: %s", err)
		os.Exit(1)
	}

	logger.Donef("\nSlack message successfully sent! 🚀\n")
}
