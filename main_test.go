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
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Empty",
			args{
				&config,
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

type TestRepository struct {
	Values map[string]string
}

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
