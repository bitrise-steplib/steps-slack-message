package main

import (
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-slack-message/lib/slack"
	"github.com/bitrise-steplib/steps-slack-message/lib/step"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	testRepository := TestRepository{
		Values: map[string]string{
			"is_debug_mode":            "yes",
			"link_names":               "yes",
			"reply_broadcast":          "yes",
			"reply_broadcast_on_error": "yes",
			"color":                    "orange",
			"timestamp":                "yes",
		},
	}
	config := step.Config{}

	type args struct {
		conf *step.Config
		repo env.Repository
	}

	defaults := args{&config, testRepository}
	tc := func(key string, value string) args {
		return args{
			&config,
			testRepository.Override(key, value),
		}
	}

	tests := []struct {
		wantErr bool
		args    args
		name    string
	}{
		{false, defaults, "Parse minimally valid config"},
		{true, tc("is_debug_mode", ""), "Invalid Debug Mode value"},
		{true, tc("link_names", ""), "Invalid LinkNames value"},
		{true, tc("reply_broadcast", ""), "Invalid ReplyBroadcast value"},
		{true, tc("reply_broadcast_on_error", ""), "Invalid ReplyBroadCastOnError value"},
		{true, tc("color", ""), "Invalid Color value"},
		{true, tc("timestamp", ""), "Invalid Timestampe value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseConfig(tt.args.conf, tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_enableDebugLog(t *testing.T) {
	testLogger := TestLogger{}
	zeroConf := step.Config{}
	debugOnConf := step.Config{Debug: true}

	type args struct {
		conf   *step.Config
		logger *TestLogger
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		expectDebugLog bool
	}{
		{"Should not enable debug log", args{&zeroConf, &testLogger}, false, false},
		{"Should enable debug log", args{&debugOnConf, &testLogger}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := enableDebugLog(tt.args.conf, tt.args.logger); (err != nil) != tt.wantErr || tt.expectDebugLog != tt.args.logger.IsDebugLogEnabled {
				t.Errorf("enableDebugLog() error = %v, wantErr %v, expectDebugLog %v", err, tt.wantErr, tt.expectDebugLog)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	testLogger := TestLogger{}

	type args struct {
		conf   *step.Config
		logger log.Logger
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		didWarn              bool
		expectedWebhookValue string
	}{
		{"No API token or Webhook URL", args{&step.Config{}, &testLogger}, true, false, ""},
		{"Has API token", args{&step.Config{APIToken: "token"}, &testLogger}, false, false, ""},
		{"Has Webhook", args{&step.Config{WebhookURL: "url"}, &testLogger}, false, false, "url"},
		{"Resets webhook when both are set", args{&step.Config{APIToken: "api", WebhookURL: "url"}, &testLogger}, false, true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.args.conf, tt.args.logger)
			didError := ((err != nil) != tt.wantErr) || (tt.didWarn != testLogger.DidWarn) || (string(tt.args.conf.WebhookURL) != tt.expectedWebhookValue)
			if didError {
				t.Errorf("validate() error = %v, didError %v", err, didError)
			}
		})
	}
}

func Test_createMessage(t *testing.T) {
	config := step.Config{
		Debug:                      false,
		WebhookURL:                 "webhook",
		WebhookURLOnError:          "webhookerr",
		APIToken:                   "",
		Channel:                    "channel",
		ChannelOnError:             "channelerr",
		Text:                       "text",
		TextOnError:                "texterr",
		IconEmoji:                  "emoji",
		IconEmojiOnError:           "emojierr",
		IconURL:                    "icon",
		IconURLOnError:             "iconerr",
		LinkNames:                  false,
		Username:                   "username",
		UsernameOnError:            "usernameerr",
		ThreadTs:                   "ts",
		ThreadTsOnError:            "tserr",
		ReplyBroadcast:             false,
		ReplyBroadcastOnError:      true,
		Color:                      "color",
		ColorOnError:               "colorerr",
		PreText:                    "pre",
		PreTextOnError:             "preerr",
		AuthorName:                 "",
		Title:                      "title",
		TitleOnError:               "titleerr",
		TitleLink:                  "",
		Message:                    "message",
		MessageOnError:             "messageerr",
		ImageURL:                   "image",
		ImageURLOnError:            "imageerr",
		ThumbURL:                   "thumb",
		ThumbURLOnError:            "thumberr",
		Footer:                     "",
		FooterIcon:                 "",
		TimeStamp:                  false,
		Fields:                     "title|value",
		Buttons:                    "text|url",
		ThreadTsOutputVariableName: "",
	}

	var message slack.Message

	type args struct {
		conf *step.Config
		msg  *slack.Message
	}
	tests := []struct {
		name               string
		args               args
		wantErr            bool
		timestamp          bool
		buildWasSuccessful bool
	}{
		{
			"Successful build message",
			args{&config, &message},
			false,
			false,
			true,
		},
		{
			"Failed build message",
			args{&config, &message},
			false,
			false,
			true,
		},
		{
			"Successful build message with timestamp",
			args{&config, &message},
			false,
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slack.BuildIsSuccessful = tt.buildWasSuccessful
			tt.args.conf.TimeStamp = tt.timestamp
			if err := createMessage(tt.args.conf, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("createMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.timestamp && message.Attachments[0].TimeStamp == 0 {
				t.Error("Failed to set timestamp")
			}
			if !tt.timestamp && message.Attachments[0].TimeStamp != 0 {
				t.Error("Set timestamp erroneously")
			}

			if tt.buildWasSuccessful {
				check1 := message.Channel == config.Channel && message.Text == config.Text && message.IconEmoji == config.IconEmoji
				check2 := message.IconURL == config.IconURL && message.Username == config.Username && message.ThreadTs == config.ThreadTs
				check3 := message.ReplyBroadcast == config.ReplyBroadcast

				att := message.Attachments[0]
				check4 := att.Fallback == config.Message && att.Color == config.Color && att.PreText == config.PreText
				check5 := att.Title == config.Title && att.Text == config.Message && att.ImageURL == config.ImageURL
				check6 := att.ThumbURL == config.ThumbURL && att.Fields[0].Title == "title" && att.Fields[0].Value == "value"
				if !(check1 && check2 && check3 && check4 && check5 && check6) {
					t.Error("Failed validation check while build was successful")
				}
			} else {
				check1 := message.Channel == config.ChannelOnError && message.Text == config.TextOnError && message.IconEmoji == config.IconEmojiOnError
				check2 := message.IconURL == config.IconURLOnError && message.Username == config.UsernameOnError && message.ThreadTs == config.ThreadTsOnError
				check3 := message.ReplyBroadcast == config.ReplyBroadcastOnError

				att := message.Attachments[0]
				check4 := att.Fallback == config.MessageOnError && att.Color == config.ColorOnError && att.PreText == config.PreTextOnError
				check5 := att.Title == config.TitleOnError && att.Text == config.MessageOnError && att.ImageURL == config.ImageURLOnError
				check6 := att.ThumbURL == config.ThumbURLOnError && att.Fields[0].Title == "title" && att.Fields[0].Value == "value"
				if !(check1 && check2 && check3 && check4 && check5 && check6) {
					t.Error("Failed validation check while build was successful")
				}
			}
		})
	}
}

type TestRepository struct {
	Values map[string]string
}

// Begin Repository

func (TestRepository) List() []string {
	//TODO implement me
	panic("implement me")
}

func (TestRepository) Unset(key string) error {
	//TODO implement me
	panic("implement me")
}

func (t TestRepository) Get(key string) string {
	return t.Values[key]
}

func (TestRepository) Set(key, value string) error {
	//TODO implement me
	panic("implement me")
}

// End

func (t TestRepository) Override(key string, value string) TestRepository {
	tmp := make(map[string]string)
	for k, v := range t.Values {
		tmp[k] = v
	}
	tmp[key] = value
	return TestRepository{
		tmp,
	}
}

// TestLogger
type TestLogger struct {
	IsDebugLogEnabled bool
	DidWarn           bool
}

func (TestLogger) Infof(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (t *TestLogger) Warnf(format string, v ...interface{}) {
	t.DidWarn = true
}

func (TestLogger) Printf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) Donef(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) Debugf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) Errorf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) TInfof(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (t *TestLogger) TWarnf(format string, v ...interface{}) {
	t.DidWarn = true
}

func (TestLogger) TPrintf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) TDonef(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) TDebugf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) TErrorf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (TestLogger) Println() {
	//TODO implement me
	panic("implement me")
}

func (t *TestLogger) EnableDebugLog(enable bool) {
	t.IsDebugLogEnabled = enable
}
