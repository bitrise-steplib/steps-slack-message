package main

import (
	"fmt"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-slack-message/lib/slack"
	"github.com/bitrise-steplib/steps-slack-message/lib/step"
	"github.com/bitrise-steplib/steps-slack-message/lib/util"
	"os"
	"os/exec"
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

func postMessage(api slack.SlackApi, msg *slack.Message, response *slack.SendMessageResponse) error {
	tmp, err := api.Post(msg)
	response = &tmp
	return err
}

func main() {
	var conf step.Config
	var msg slack.Message
	var slackApi slack.SlackApi
	var response slack.SendMessageResponse
	envRepo := env.NewRepository()

	err := run(
		func() error { return parseConfig(&conf, envRepo) },
		func() error { return enableDebugLog(&conf) },
		func() error { return validate(&conf) },
		func() error { return createSlackClient(&conf, slackApi, &logger) },
		func() error { return createMessage(&conf, &msg) },
		func() error { return postMessage(slackApi, &msg, &response) },
		func() error { return exportEnvironmentVariables(&response, &conf) },
	)
	if err != nil {
		logger.Errorf("Error: %s", err)
		os.Exit(1)
	}

	logger.Donef("\nSlack message successfully sent! 🚀\n")
}

func exportEnvironmentVariables(response *slack.SendMessageResponse, config *step.Config) error {
	/// Exports env using envman
	c := exec.Command("envman", "add", "--key", string(config.ThreadTsOutputVariableName), "--value", response.Timestamp)
	err := c.Run()
	if err != nil {
		return fmt.Errorf("Failed to run envman %s", err)
	}
	return nil
}

func createSlackClient(conf *step.Config, client slack.SlackApi, logger *log.Logger) error {
	selector := util.SeedSelect[string](success)(string(conf.WebhookURL), string(conf.WebhookURLOnError))
	client = &slack.SlackClient{
		Conf:               conf,
		Logger:             logger,
		WebhookUrlSelector: &selector,
	}
	return nil
}
