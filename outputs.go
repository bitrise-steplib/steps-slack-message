package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/log"
)

// SendMessageResponse is the response from Slack POST
type SendMessageResponse struct {
	/// The Thread Timestamp
	Timestamp string `json:"ts"`
}

// / Export the output variables after a successful response
func exportOutputs(conf *config, resp *http.Response) error {

	if !isRequestingOutput(conf) {
		log.Debugf("not requesting any outputs")
		return nil
	}

	isWebhook := strings.TrimSpace(conf.WebhookURL) != ""

	// Slack webhooks do not return any useful response information
	if isWebhook {
		return fmt.Errorf("for output support, do not submit a WebHook URL")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %s", err)
	}

	log.Debugf("Response:\n%s", string(body))

	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// here we want to fail, because the user is expecting an output
		return fmt.Errorf("failed to parse response: %s", err)
	}

	if string(conf.ThreadTsOutputVariableName) != "" {
		log.Debugf("Exporting output: %s=%s\n", string(conf.ThreadTsOutputVariableName), response.Timestamp)
		err := exportEnvVariable(string(conf.ThreadTsOutputVariableName), response.Timestamp)
		if err != nil {
			return err
		}
	}

	return nil

}

// / Checks if we are requesting an output of anything
func isRequestingOutput(conf *config) bool {
	return string(conf.ThreadTsOutputVariableName) != ""
}

// / Exports env using envman
func exportEnvVariable(variable string, value string) error {
	c := exec.Command("envman", "add", "--key", variable, "--value", value)
	err := c.Run()
	if err != nil {
		return fmt.Errorf("Failed to run envman %s", err)
	}
	return nil
}
