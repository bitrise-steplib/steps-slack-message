package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

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
		log.Debugf("Not requesting any outputs")
		return nil
	}

	var response SendMessageResponse
	parseError := json.NewDecoder(resp.Body).Decode(&response)
	if parseError != nil {
		// here we want to fail, because the user is expecting an output
		return fmt.Errorf("Failed to parse response: %s", parseError)
	}

	if response.Timestamp == "" {
		return fmt.Errorf("Response does not contain a timestamp, cannot export output")
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
