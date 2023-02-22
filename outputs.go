package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/bitrise-io/go-utils/log"
)

// SendMessageResponse is the response from Slack POST
// https://api.slack.com/methods/chat.postMessage#examples
type SendMessageResponse struct {
	/// The status of the request. When `false`, check the Error for a reason
	Ok bool `json:"ok"`

	/// Describes an error that prevented the message from being sent
	Error string `json:"error"`

	/// The Thread Timestamp
	Timestamp string `json:"ts"`
}

// Check the response status and set any output variables if required
func processResponse(conf *Config, resp *http.Response) error {
	// if the request was made using a legacy webhook url, skip processing as there is no response.
	if conf.APIToken == "" {
		log.Debugf("Skipping response processing because legacy webhook urls do not return any content")
		return nil
	}

	var response SendMessageResponse
	parseError := json.NewDecoder(resp.Body).Decode(&response)
	if parseError != nil {
		// here we want to fail, because the user is expecting an output
		return fmt.Errorf("Failed to parse response: %s", parseError)
	}

	// if slack didn't return 'ok', fail with the error code
	if !response.Ok {
		return fmt.Errorf("Slack responded with error: %s", response.Error)
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

// Checks if we are requesting an output of anything
func isRequestingOutput(conf *Config) bool {
	return string(conf.ThreadTsOutputVariableName) != ""
}

// Exports env using envman
func exportEnvVariable(variable string, value string) error {
	c := exec.Command("envman", "add", "--key", variable, "--value", value)
	err := c.Run()
	if err != nil {
		return fmt.Errorf("Failed to run envman %s", err)
	}
	return nil
}
