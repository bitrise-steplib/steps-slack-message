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

// Stage A stage represents a discrete unit of work in the execution of this step.  Each stage should do one thing. Due
// to current imagination limitations, stages can only emit errors but may take references to objects that need to be
// created or updated during the execution of a stage.
type Stage func() error

// Execute stages in sequence.  Stop execution if an error is received from a stage.
func run(stages ...Stage) error {
	for _, stage := range stages {
		if err := stage(); err != nil {
			return err
		}
	}
	return nil
}

// Begin stages
func parseConfig(conf *step.Config, repo env.Repository) error {
	if err := stepconf.NewInputParser(repo).Parse(conf); err != nil {
		return err
	}
	stepconf.Print(conf)
	return nil
}

func enableDebugLog(conf *step.Config, logger log.Logger) error {
	logger.EnableDebugLog(conf.Debug)
	return nil
}

func validate(conf *step.Config, logger log.Logger) error {
	if conf.APIToken == "" && conf.WebhookURL == "" {
		return fmt.Errorf("Both API Token and WebhookURL are empty. You need to provide one of them. If you want to use incoming webhooks provide the webhook url. If you want to use a bot to send a message provide the bot API token")
	}

	if conf.APIToken != "" && conf.WebhookURL != "" {
		logger.Warnf("Both API Token and WebhookURL are provided. Using the API Token")
		conf.WebhookURL = ""
	}
	return nil
}

func createMessage(conf *step.Config, msg *slack.Message) error {
	*msg = slack.NewMessage(*conf)
	return nil
}

func postMessage(api slack.SlackApi, msg *slack.Message, response *slack.SendMessageResponse) error {
	tmp, err := api.Post(msg)
	*response = tmp
	return err
}

func exportEnvironmentVariables(response *slack.SendMessageResponse, conf *step.Config, logger log.Logger) error {
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
	selector := util.SeedSelect[string](slack.BuildIsSuccessful)(string(conf.WebhookURL), string(conf.WebhookURLOnError))
	client = &slack.SlackClient{
		Conf:               conf,
		Logger:             logger,
		WebhookUrlSelector: &selector,
	}
	return nil
}

// End stages

func main() {
	envRepo := env.NewRepository()
	logger := log.NewLogger()

	// Populated by stages
	var conf step.Config
	var msg slack.Message
	var slackApi slack.SlackApi
	var response slack.SendMessageResponse

	err := run(
		func() error { return parseConfig(&conf, envRepo) },
		func() error { return enableDebugLog(&conf, logger) },
		func() error { return validate(&conf, logger) },
		func() error { return createSlackClient(&conf, slackApi, &logger) },
		func() error { return createMessage(&conf, &msg) },
		func() error { return postMessage(slackApi, &msg, &response) },
		func() error { return exportEnvironmentVariables(&response, &conf, logger) },
	)
	if err != nil {
		logger.Errorf("Error: %s", err)
		os.Exit(1)
	}

	logger.Donef("\nSlack message successfully sent! 🚀\n")
}
