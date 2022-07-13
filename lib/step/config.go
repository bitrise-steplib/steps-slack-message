package step

import "github.com/bitrise-io/go-steputils/v2/stepconf"

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
