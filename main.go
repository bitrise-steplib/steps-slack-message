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

// Config ...
type Config struct {
	Debug bool `env:"is_debug_mode,opt[yes,no]"`

	// Message
	WebhookURL            stepconf.Secret `env:"webhook_url"`
	WebhookURLOnError     stepconf.Secret `env:"webhook_url_on_error"`
	APIToken              stepconf.Secret `env:"api_token"`
	Channel               string          `env:"channel"`
	ChannelOnError        string          `env:"channel_on_error"`
	Text                  string          `env:"text"`
	TextOnError           string          `env:"text_on_error"`
	IconEmoji             string          `env:"emoji"`
	IconEmojiOnError      string          `env:"emoji_on_error"`
	IconURL               string          `env:"icon_url"`
	IconURLOnError        string          `env:"icon_url_on_error"`
	LinkNames             bool            `env:"link_names,opt[yes,no]"`
	Username              string          `env:"from_username"`
	UsernameOnError       string          `env:"from_username_on_error"`
	ThreadTs              string          `env:"thread_ts"`
	ThreadTsOnError       string          `env:"thread_ts_on_error"`
	ReplyBroadcast        bool            `env:"reply_broadcast,opt[yes,no]"`
	ReplyBroadcastOnError bool            `env:"reply_broadcast_on_error,opt[yes,no]"`

	// Attachment
	Color           string `env:"color,required"`
	ColorOnError    string `env:"color_on_error"`
	PreText         string `env:"pretext"`
	PreTextOnError  string `env:"pretext_on_error"`
	AuthorName      string `env:"author_name"`
	Title           string `env:"title"`
	TitleOnError    string `env:"title_on_error"`
	TitleLink       string `env:"title_link"`
	Message         string `env:"message"`
	MessageOnError  string `env:"message_on_error"`
	ImageURL        string `env:"image_url"`
	ImageURLOnError string `env:"image_url_on_error"`
	ThumbURL        string `env:"thumb_url"`
	ThumbURLOnError string `env:"thumb_url_on_error"`
	Footer          string `env:"footer"`
	FooterIcon      string `env:"footer_icon"`
	TimeStamp       bool   `env:"timestamp,opt[yes,no]"`
	Fields          string `env:"fields"`
	Buttons         string `env:"buttons"`

	// Step Outputs
	ThreadTsOutputVariableName string `env:"output_thread_ts"`
}

// success is true if the build is successful, false otherwise.
var success = os.Getenv("BITRISE_BUILD_STATUS") == "0"

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

func newMessage(c Config) Message {
	msg := Message{
		Channel: strings.TrimSpace(selectValue(c.Channel, c.ChannelOnError)),
		Text:    selectValue(c.Text, c.TextOnError),
		Attachments: []Attachment{{
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

// postMessage sends a message to a channel.
func postMessage(conf Config, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Debugf("Request to Slack: %s\n", b)

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

	if err := processResponse(&conf, resp); err != nil {
		return err
	}

	return nil
}

func validate(conf *Config) error {
	if conf.APIToken == "" && conf.WebhookURL == "" {
		return fmt.Errorf("Both API Token and WebhookURL are empty. You need to provide one of them. If you want to use incoming webhooks provide the webhook url. If you want to use a bot to send a message provide the bot API token")
	}

	if conf.APIToken != "" && conf.WebhookURL != "" {
		log.Warnf("Both API Token and WebhookURL are provided. Using the API Token")
		conf.WebhookURL = ""

	}

	if conf.APIToken == "" && isRequestingOutput(conf) {
		return fmt.Errorf("Outputs can only be set when using the API Token.")
	}

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

	if err := validate(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	msg := newMessage(conf)
	if err := postMessage(conf, msg); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nSlack message successfully sent! 🚀\n")
}
