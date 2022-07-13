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
)

var logger = log.NewLogger()

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
	*msg = util.NewMessage(*conf)
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

func exportEnvironmentVariables(response *slack.SendMessageResponse, conf *step.Config) error {
	if string(conf.ThreadTsOutputVariableName) == "" {
		return nil
	}
	logger.Debugf("Exporting output: %s=%s\n", conf.ThreadTsOutputVariableName, response.Timestamp)

	/// Exports env using envman
	c := exec.Command("envman", "add", "--key", string(conf.ThreadTsOutputVariableName), "--value", response.Timestamp)
	err := c.Run()
	if err != nil {
		return fmt.Errorf("Failed to run envman %s", err)
	}
	return nil
}

func createSlackClient(conf *step.Config, client slack.SlackApi, logger *log.Logger) error {
	selector := util.SeedSelect[string](util.BuildIsSuccessful)(string(conf.WebhookURL), string(conf.WebhookURLOnError))
	client = &slack.SlackClient{
		Conf:               conf,
		Logger:             logger,
		WebhookUrlSelector: &selector,
	}
	return nil
}
