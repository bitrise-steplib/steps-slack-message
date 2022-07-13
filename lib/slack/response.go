package slack

import (
	"fmt"
	"os/exec"
)

// SendMessageResponse is the response from Slack POST
type SendMessageResponse struct {
	/// The Thread Timestamp
	Timestamp string `json:"ts"`
}

/// Exports env using envman
func exportEnvVariable(variable string, value string) error {
	c := exec.Command("envman", "add", "--key", variable, "--value", value)
	err := c.Run()
	if err != nil {
		return fmt.Errorf("Failed to run envman %s", err)
	}
	return nil
}
