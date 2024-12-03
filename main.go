package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/hashicorp/go-retryablehttp"
)

// Input ...
type Input struct {
	Debug         bool            `env:"is_debug_mode,opt[yes,no]"`
	BuildAPIToken stepconf.Secret `env:"BITRISE_BUILD_API_TOKEN,required"`
	BuildURL      string          `env:"BITRISE_BUILD_URL,required"`

	// Message
	WebhookURL            stepconf.Secret `env:"webhook_url"`
	WebhookURLOnError     stepconf.Secret `env:"webhook_url_on_error"`
	APIToken              stepconf.Secret `env:"api_token"`
	IntegrationID         string          `env:"workspace_integration_id"`
	IntegrationIDOnError  string          `env:"workspace_integration_id_on_error"`
	Channel               string          `env:"channel"`
	ChannelOnError        string          `env:"channel_on_error"`
	Text                  string          `env:"text"`
	TextOnError           string          `env:"text_on_error"`
	IconEmoji             string          `env:"emoji"`
	IconEmojiOnError      string          `env:"emoji_on_error"`
	IconURL               string          `env:"icon_url"`
	IconURLOnError        string          `env:"icon_url_on_error"`
	LinkNames             bool            `env:"link_names,opt[yes,no]"`
	Username              string          `env:"from_username"`
	UsernameOnError       string          `env:"from_username_on_error"`
	ThreadTs              string          `env:"thread_ts"`
	ThreadTsOnError       string          `env:"thread_ts_on_error"`
	Ts                    string          `env:"ts"`
	TsOnError             string          `env:"ts_on_error"`
	ReplyBroadcast        bool            `env:"reply_broadcast,opt[yes,no]"`
	ReplyBroadcastOnError bool            `env:"reply_broadcast_on_error,opt[yes,no]"`

	// Attachment
	Color             string `env:"color,required"`
	ColorOnError      string `env:"color_on_error"`
	PreText           string `env:"pretext"`
	PreTextOnError    string `env:"pretext_on_error"`
	AuthorName        string `env:"author_name"`
	Title             string `env:"title"`
	TitleOnError      string `env:"title_on_error"`
	TitleLink         string `env:"title_link"`
	Message           string `env:"message"`
	MessageOnError    string `env:"message_on_error"`
	ImageURL          string `env:"image_url"`
	ImageURLOnError   string `env:"image_url_on_error"`
	ThumbURL          string `env:"thumb_url"`
	ThumbURLOnError   string `env:"thumb_url_on_error"`
	Footer            string `env:"footer"`
	FooterOnError     string `env:"footer_on_error"`
	FooterIcon        string `env:"footer_icon"`
	FooterIconOnError string `env:"footer_icon_on_error"`
	TimeStamp         bool   `env:"timestamp,opt[yes,no]"`
	Fields            string `env:"fields"`
	Buttons           string `env:"buttons"`

	// Status
	BuildStatus         string `env:"build_status"`
	PipelineBuildStatus string `env:"pipeline_build_status"`

	// Step Outputs
	ThreadTsOutputVariableName string `env:"output_thread_ts"`
}

type config struct {
	Debug bool `env:"is_debug_mode,opt[yes,no]"`

	// Message
	APIToken       stepconf.Secret `env:"api_token"`
	WebhookURL     string
	Channel        string
	Text           string
	IconEmoji      string
	IconURL        string
	Username       string
	ThreadTs       string
	Ts             string
	ReplyBroadcast bool
	LinkNames      bool `env:"link_names,opt[yes,no]"`

	// Attachment
	Color      string
	PreText    string
	Title      string
	Message    string
	ImageURL   string
	ThumbURL   string
	AuthorName string `env:"author_name"`
	TitleLink  string `env:"title_link"`
	Footer     string `env:"footer"`
	FooterIcon string `env:"footer_icon"`
	TimeStamp  bool   `env:"timestamp,opt[yes,no]"`
	Fields     string `env:"fields"`
	Buttons    string `env:"buttons"`

	// Step Outputs
	ThreadTsOutputVariableName string `env:"output_thread_ts"`
}

// ensureNewlines replaces all \n substrings with newline characters.
func ensureNewlines(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func newMessage(c config) Message {
	msg := Message{
		Channel: strings.TrimSpace(c.Channel),
		Text:    c.Text,
		Attachments: []Attachment{{
			Fallback:   ensureNewlines(c.Message),
			Color:      c.Color,
			PreText:    c.PreText,
			AuthorName: c.AuthorName,
			Title:      c.Title,
			TitleLink:  c.TitleLink,
			Text:       ensureNewlines(c.Message),
			Fields:     parseFields(c.Fields),
			ImageURL:   c.ImageURL,
			ThumbURL:   c.ThumbURL,
			Footer:     c.Footer,
			FooterIcon: c.FooterIcon,
			Buttons:    parseButtons(c.Buttons),
		}},
		IconEmoji:      c.IconEmoji,
		IconURL:        c.IconURL,
		LinkNames:      c.LinkNames,
		Username:       c.Username,
		ThreadTs:       c.ThreadTs,
		Ts:             c.Ts,
		ReplyBroadcast: c.ReplyBroadcast,
	}
	if c.TimeStamp {
		msg.Attachments[0].TimeStamp = int(time.Now().Unix())
	}
	return msg
}

func getWebhookURL(buildURL string, id string, token string) (string, error) {
	var webookData struct {
		WebhookURL string `json:"webhook_url"`
	}
	siURL := fmt.Sprintf("%s/integrations/slack/%s", buildURL, id)

	req, err := retryablehttp.NewRequest("GET", siURL, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Add("Build-Api-Token", token)
	client := retry.NewHTTPClient()

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if err = json.Unmarshal(body, &webookData); err != nil {
			return "", err
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("server error, status: %s\nresponse: %s", resp.Status, string(body))
	}
	return webookData.WebhookURL, nil
}

// postMessage sends a message to a channel.
func postMessage(conf config, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Debugf("Request to Slack: %s\n", b)

	url := strings.TrimSpace(conf.WebhookURL)
	ts := strings.TrimSpace(conf.Ts)

	if url == "" {
		if ts == "" {
			url = "https://slack.com/api/chat.postMessage"
		} else {
			url = "https://slack.com/api/chat.update"
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if string(conf.APIToken) != "" {
		req.Header.Add("Authorization", "Bearer "+string(conf.APIToken))
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("server error: %s, failed to read response: %s", resp.Status, err)
		}
		return fmt.Errorf("server error: %s, response: %s", resp.Status, body)
	}

	if err := exportOutputs(&conf, resp); err != nil {
		return fmt.Errorf("failed to export outputs: %s", err)
	}

	return nil
}

func validate(inp *Input) error {
	if inp.APIToken == "" && inp.WebhookURL == "" && inp.IntegrationID == "" {
		return fmt.Errorf("All of Integration ID, API Token and WebhookURL are empty. You need to provide one of them. If you want to use incoming webhooks provide the webhook url. If you want to use a bot to send a message provide the bot API token. If you want to use a configured workspace integration use its ID.")
	}

	if inp.IntegrationID != "" {
		if inp.APIToken != "" {
			log.Warnf("Both API Token and Integration ID are provided. Ignoring API Token.")
			inp.APIToken = ""
		}
		if inp.WebhookURL != "" {
			log.Warnf("Both WebhookURL and Integration ID are provided. Ignoring WebhookURL.")
			inp.WebhookURL = ""
		}
		return nil
	}

	if inp.APIToken != "" && inp.WebhookURL != "" {
		log.Warnf("Both API Token and WebhookURL are provided. Using the API Token")
		inp.WebhookURL = ""

	}
	return nil
}

func parseInputIntoConfig(inp *Input) (config, error) {
	pipelineSuccess := inp.PipelineBuildStatus == "" ||
		inp.PipelineBuildStatus == "succeeded" ||
		inp.PipelineBuildStatus == "succeeded_with_abort"
	success := pipelineSuccess && inp.BuildStatus == "0"

	// selectValue chooses the right value based on the result of the build.
	var selectValue = func(ifSuccess, ifFailed string) string {
		if success || ifFailed == "" {
			return ifSuccess
		}
		return ifFailed
	}
	var integrationID = selectValue(inp.IntegrationID, inp.IntegrationIDOnError)
	var webhookURL = selectValue(string(inp.WebhookURL), string(inp.WebhookURLOnError))
	if integrationID != "" {
		var err error
		webhookURL, err = getWebhookURL(inp.BuildURL, integrationID, string(inp.BuildAPIToken))
		if err != nil {
			return config{}, err
		}
	}

	var config = config{
		Debug:                      inp.Debug,
		APIToken:                   inp.APIToken,
		WebhookURL:                 webhookURL,
		Channel:                    selectValue(inp.Channel, inp.ChannelOnError),
		Text:                       selectValue(inp.Text, inp.TextOnError),
		IconEmoji:                  selectValue(inp.IconEmoji, inp.IconEmojiOnError),
		IconURL:                    selectValue(inp.IconURL, inp.IconURLOnError),
		Username:                   selectValue(inp.Username, inp.UsernameOnError),
		ThreadTs:                   selectValue(inp.ThreadTs, inp.ThreadTsOnError),
		ReplyBroadcast:             (success && inp.ReplyBroadcast) || (!success && inp.ReplyBroadcastOnError),
		LinkNames:                  inp.LinkNames,
		Color:                      selectValue(inp.Color, inp.ColorOnError),
		PreText:                    selectValue(inp.PreText, inp.PreTextOnError),
		Title:                      selectValue(inp.Title, inp.TitleOnError),
		Message:                    selectValue(inp.Message, inp.MessageOnError),
		ImageURL:                   selectValue(inp.ImageURL, inp.ImageURLOnError),
		ThumbURL:                   selectValue(inp.ThumbURL, inp.ThumbURLOnError),
		AuthorName:                 inp.AuthorName,
		TitleLink:                  inp.TitleLink,
		Footer:                     selectValue(inp.Footer, inp.FooterOnError),
		FooterIcon:                 selectValue(inp.FooterIcon, inp.FooterIconOnError),
		TimeStamp:                  inp.TimeStamp,
		Fields:                     inp.Fields,
		Buttons:                    inp.Buttons,
		ThreadTsOutputVariableName: inp.ThreadTsOutputVariableName,
		Ts:                         selectValue(inp.Ts, inp.TsOnError),
	}
	return config, nil

}

func main() {
	var input Input
	if err := stepconf.Parse(&input); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(input)
	log.SetEnableDebugLog(input.Debug)

	if err := validate(&input); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	config, err := parseInputIntoConfig(&input)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	msg := newMessage(config)
	if err := postMessage(config, msg); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nSlack message successfully sent! 🚀\n")
}
