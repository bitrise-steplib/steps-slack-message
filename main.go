package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// success is true if the build is successful, false otherwise.
var success = os.Getenv("BITRISE_BUILD_STATUS") == "0"

// byStatus chooses the right value for fields based on the result of the build.
func byStatus(onSuccess, onError string) string {
	if success || onError == "" {
		return onSuccess
	}
	return onError
}

// ensureNewlines replaces all \n substrings with newline characters.
func ensureNewlines(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func newMessage(c Config) Message {
	msg := Message{
		Channel: byStatus(c.Channel, c.ChannelOnError),
		Text:    byStatus(c.Text, c.TextOnError),
		Attachments: []Attachment{{
			Fallback:   ensureNewlines(byStatus(c.Message, c.MessageOnError)),
			Color:      byStatus(c.Color, c.ColorOnError),
			PreText:    byStatus(c.PreText, c.PreTextOnError),
			AuthorName: c.AuthorName,
			Title:      byStatus(c.Title, c.TitleOnError),
			TitleLink:  c.TitleLink,
			Text:       ensureNewlines(byStatus(c.Message, c.MessageOnError)),
			Fields:     parseFields(c.Fields),
			ImageURL:   byStatus(c.ImageURL, c.ImageURLOnError),
			ThumbURL:   byStatus(c.ThumbURL, c.ThumbURLOnError),
			Footer:     c.Footer,
			FooterIcon: c.FooterIcon,
			Buttons:    parseButtons(c.Buttons),
		}},
		IconEmoji: byStatus(c.IconEmoji, c.IconEmojiOnError),
		IconURL:   byStatus(c.IconURL, c.IconURLOnError),
		LinkNames: c.LinkNames,
		Username:  byStatus(c.Username, c.UsernameOnError),
	}
	if c.TimeStamp {
		msg.Attachments[0].TimeStamp = int(time.Now().Unix())
	}
	return msg
}

// postMessage sends a message to a channel.
func postMessage(webhookURL string, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Debugf("Request to Slack: %s\n", b)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("Response from Slack: %s\n", body)
	return nil
}

func main() {
	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(conf)
	log.SetEnableDebugLog(conf.Debug)

	msg := newMessage(conf)
	if err := postMessage(conf.WebhookURL, msg); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nSlack message successfully sent! 🚀\n")
}
