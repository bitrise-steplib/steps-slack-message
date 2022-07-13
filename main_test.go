package main

import (
	"github.com/bitrise-io/go-utils/v2/env"
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
