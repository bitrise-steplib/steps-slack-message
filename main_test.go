package main

import (
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-steplib/steps-slack-message/lib/step"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	var testRepository TestRepository
	config := newTestConfig()

	type args struct {
		conf *step.Config
		repo env.Repository
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Empty",
			args{
				config,
				testRepository,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseConfig(tt.args.conf, tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func newTestConfig() *step.Config {
	return &step.Config{
		Debug:                      false,
		WebhookURL:                 "",
		WebhookURLOnError:          "",
		APIToken:                   "",
		Channel:                    "",
		ChannelOnError:             "",
		Text:                       "",
		TextOnError:                "",
		IconEmoji:                  "",
		IconEmojiOnError:           "",
		IconURL:                    "",
		IconURLOnError:             "",
		LinkNames:                  false,
		Username:                   "",
		UsernameOnError:            "",
		ThreadTs:                   "",
		ThreadTsOnError:            "",
		ReplyBroadcast:             false,
		ReplyBroadcastOnError:      false,
		Color:                      "",
		ColorOnError:               "",
		PreText:                    "",
		PreTextOnError:             "",
		AuthorName:                 "",
		Title:                      "",
		TitleOnError:               "",
		TitleLink:                  "",
		Message:                    "",
		MessageOnError:             "",
		ImageURL:                   "",
		ImageURLOnError:            "",
		ThumbURL:                   "",
		ThumbURLOnError:            "",
		Footer:                     "",
		FooterIcon:                 "",
		TimeStamp:                  false,
		Fields:                     "",
		Buttons:                    "",
		ThreadTsOutputVariableName: "",
	}
}

type TestRepository struct{}

func (TestRepository) List() []string {
	//TODO implement me
	panic("implement me")
}

func (TestRepository) Unset(key string) error {
	//TODO implement me
	panic("implement me")
}

func (TestRepository) Get(key string) string {
	//TODO implement me
	panic("implement me")
}

func (TestRepository) Set(key, value string) error {
	//TODO implement me
	panic("implement me")
}
