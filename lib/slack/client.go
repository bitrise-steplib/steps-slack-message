package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-slack-message/lib/step"
	"github.com/bitrise-steplib/steps-slack-message/lib/util"
	"io/ioutil"
	"net/http"
	"strings"
)

type SlackClient struct {
	Conf               *step.Config
	Logger             *log.Logger
	WebhookUrlSelector *util.Select[string]
}

func (c SlackClient) Post(msg *Message) (resp SendMessageResponse, err error) {
	conf := c.Conf
	logger := *c.Logger

	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	logger.Debugf("Request to Slack: %s\n", b)

	url := strings.TrimSpace(c.WebhookUrlSelector.Get())
	if url == "" {
		url = "https://slack.com/api/chat.postMessage"
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if string(conf.APIToken) != "" {
		req.Header.Add("Authorization", "Bearer "+string(conf.APIToken))
	}

	client := &http.Client{}

	httpResp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to send the request: %s", err)
		return
	}
	defer func() {
		if cerr := httpResp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	if httpResp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			err = fmt.Errorf("server error: %s, failed to read response: %s", httpResp.Status, err)
			return
		}
		err = fmt.Errorf("server error: %s, response: %s", httpResp.Status, body)
		return
	}

	if resp, err = parseMessageResponse(IsWebhook(c.WebhookUrlSelector), conf, httpResp, logger); err != nil {
		err = fmt.Errorf("failed to export outputs: %s", err)
		return
	}

	return
}

func IsWebhook(url *util.Select[string]) bool {
	return strings.TrimSpace(url.Get()) != ""
}

/// Export the output variables after a successful response
func parseMessageResponse(isWebhook bool, conf *step.Config, resp *http.Response, logger log.Logger) (SendMessageResponse, error) {
	var response SendMessageResponse

	if !isRequestingOutput(conf) {
		logger.Debugf("Not requesting any outputs")
		return response, nil
	}

	// Slack webhooks do not return any useful response information
	if isWebhook {
		return response, fmt.Errorf("For output support, do not submit a WebHook URL")
	}

	parseError := json.NewDecoder(resp.Body).Decode(&response)
	if parseError != nil {
		// here we want to fail, because the user is expecting an output
		return response, fmt.Errorf("Failed to parse response: %s", parseError)
	}

	return response, nil
}

/// Checks if we are requesting an output of anything
func isRequestingOutput(conf *step.Config) bool {
	return string(conf.ThreadTsOutputVariableName) != ""
}
