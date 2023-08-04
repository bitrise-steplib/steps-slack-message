package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

type InputConfig struct {
	Debug bool `env:"is_debug_mode,opt[yes,no]"`

	// Message
	APIToken  stepconf.Secret `env:"api_token"`
	LinkNames bool            `env:"link_names,opt[yes,no]"`

	// Attachment
	AuthorName string `env:"author_name"`
	TitleLink  string `env:"title_link"`
	Footer     string `env:"footer"`
	FooterIcon string `env:"footer_icon"`
	TimeStamp  bool   `env:"timestamp,opt[yes,no]"`
	Fields     string `env:"fields"`
	Buttons    string `env:"buttons"`

	// Step Outputs
	ThreadTsOutputVariableName string `env:"output_thread_ts"`
}

// Input ...
type Input struct {
	InputConfig

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
	Title           string `env:"title"`
	TitleOnError    string `env:"title_on_error"`
	Message         string `env:"message"`
	MessageOnError  string `env:"message_on_error"`
	ImageURL        string `env:"image_url"`
	ImageURLOnError string `env:"image_url_on_error"`
	ThumbURL        string `env:"thumb_url"`
	ThumbURLOnError string `env:"thumb_url_on_error"`

	// Status
	BuildStatus         string `env:"build_status"`
	PipelineBuildStatus string `env:"pipeline_build_status"`

	// Step Outputs
	ThreadTsOutputVariableName string `env:"output_thread_ts"`
}

type Config struct {
	InputConfig

	// Message
	WebhookURL     string
	Channel        string
	Text           string
	IconEmoji      string
	IconURL        string
	Username       string
	ThreadTs       string
	ReplyBroadcast bool

	// Attachment
	Color    string
	PreText  string
	Title    string
	Message  string
	ImageURL string
	ThumbURL string
}

// ensureNewlines replaces all \n substrings with newline characters.
func ensureNewlines(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func newMessage(c Config) Message {
	msg := Message{
		Channel: strings.TrimSpace(c.Channel),
		Text:    c.Text,
		Attachments: []Attachment{{
			Fallback:   ensureNewlines(c.Message),
			Color:      c.Color,
			PreText:    c.PreText,
			AuthorName: c.AuthorName,
			Title:      c.Title,
			TitleLink:  c.TitleLink,
			Text:       ensureNewlines(c.Message),
			Fields:     parseFields(c.Fields),
			ImageURL:   c.ImageURL,
			ThumbURL:   c.ThumbURL,
			Footer:     c.Footer,
			FooterIcon: c.FooterIcon,
			Buttons:    parseButtons(c.Buttons),
		}},
		IconEmoji:      c.IconEmoji,
		IconURL:        c.IconURL,
		LinkNames:      c.LinkNames,
		Username:       c.Username,
		ThreadTs:       c.ThreadTs,
		ReplyBroadcast: c.ReplyBroadcast,
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

	url := strings.TrimSpace(conf.WebhookURL)
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
		body, err := io.ReadAll(resp.Body)
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

func validate(conf *Input) error {
	if conf.APIToken == "" && conf.WebhookURL == "" {
		return fmt.Errorf("Both API Token and WebhookURL are empty. You need to provide one of them. If you want to use incoming webhooks provide the webhook url. If you want to use a bot to send a message provide the bot API token")
	}

	if conf.APIToken != "" && conf.WebhookURL != "" {
		log.Warnf("Both API Token and WebhookURL are provided. Using the API Token")
		conf.WebhookURL = ""

	}
	return nil
}

func parseConfig(cfg *Input) Config {
	pipelineSuccess := cfg.PipelineBuildStatus == "" ||
		cfg.PipelineBuildStatus == "succeeded" ||
		cfg.PipelineBuildStatus == "succeeded_with_abort"
	success := pipelineSuccess && cfg.BuildStatus == "0"

	// selectValue chooses the right value based on the result of the build.
	var selectValue = func(ifSuccess, ifFailed string) string {
		if success || ifFailed == "" {
			return ifSuccess
		}
		return ifFailed
	}

	// selectBool chooses the right boolean value based on the result of the build.
	var selectBool = func(ifSuccess, ifFailed bool) bool {
		if success {
			return ifSuccess
		}
		return ifFailed
	}

	var input = Config{
		InputConfig:    cfg.InputConfig,
		WebhookURL:     selectValue(string(cfg.WebhookURL), string(cfg.WebhookURLOnError)),
		Channel:        selectValue(cfg.Channel, cfg.ChannelOnError),
		Text:           selectValue(cfg.Text, cfg.TextOnError),
		IconEmoji:      selectValue(cfg.IconEmoji, cfg.IconEmojiOnError),
		IconURL:        selectValue(cfg.IconURL, cfg.IconURLOnError),
		Username:       selectValue(cfg.Username, cfg.UsernameOnError),
		ThreadTs:       selectValue(cfg.ThreadTs, cfg.ThreadTsOnError),
		ReplyBroadcast: selectBool(cfg.ReplyBroadcast, cfg.ReplyBroadcastOnError),
		Color:          selectValue(cfg.Color, cfg.ColorOnError),
		PreText:        selectValue(cfg.PreText, cfg.PreTextOnError),
		Title:          selectValue(cfg.Title, cfg.TitleOnError),
		Message:        selectValue(cfg.Message, cfg.MessageOnError),
		ImageURL:       selectValue(cfg.ImageURL, cfg.ImageURLOnError),
		ThumbURL:       selectValue(cfg.ThumbURL, cfg.ThumbURLOnError),
	}
	return input

}

func main() {
	var conf Input
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

	input := parseConfig(&conf)

	msg := newMessage(input)
	if err := postMessage(input, msg); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nSlack message successfully sent! 🚀\n")
}
