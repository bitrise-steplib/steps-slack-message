package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/bitrise-io/go-utils/colorstring"
)

// ConfigsModel ...
type ConfigsModel struct {
	// Slack Inputs
	WebhookURL          string
	Channel             string
	FromUsername        string
	FromUsernameOnError string
	Message             string
	MessageOnError      string
	Emoji               string
	EmojiOnError        string
	IconURL             string
	IconURLOnError      string
	// Other Inputs
	IsDebugMode bool
	// Other configs
	IsBuildFailed bool
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		WebhookURL:          os.Getenv("webhook_url"),
		Channel:             os.Getenv("channel"),
		FromUsername:        os.Getenv("from_username"),
		FromUsernameOnError: os.Getenv("from_username_on_error"),
		Message:             os.Getenv("message"),
		MessageOnError:      os.Getenv("message_on_error"),
		Emoji:               os.Getenv("emoji"),
		EmojiOnError:        os.Getenv("emoji_on_error"),
		IconURL:             os.Getenv("icon_url"),
		IconURLOnError:      os.Getenv("icon_url_on_error"),
		//
		IsDebugMode: (os.Getenv("is_debug_mode") == "yes"),
		//
		IsBuildFailed: (os.Getenv("STEPLIB_BUILD_STATUS") != "0"),
	}
}

func (inputs ConfigsModel) print() {
	fmt.Println("")
	fmt.Println(colorstring.Blue("Slack configs:"))
	fmt.Println(" - WebhookURL:", inputs.WebhookURL)
	fmt.Println(" - Channel:", inputs.Channel)
	fmt.Println(" - FromUsername:", inputs.FromUsername)
	fmt.Println(" - FromUsernameOnError:", inputs.FromUsernameOnError)
	fmt.Println(" - Message:", inputs.Message)
	fmt.Println(" - MessageOnError:", inputs.MessageOnError)
	fmt.Println(" - Emoji:", inputs.Emoji)
	fmt.Println(" - EmojiOnError:", inputs.EmojiOnError)
	fmt.Println(" - IconURL:", inputs.IconURL)
	fmt.Println(" - IconURLOnError:", inputs.IconURLOnError)
	fmt.Println("")
	fmt.Println(colorstring.Blue("Other configs:"))
	fmt.Println(" - IsDebugMode:", inputs.IsDebugMode)
	fmt.Println(" - IsBuildFailed:", inputs.IsBuildFailed)
	fmt.Println("")
}

func (inputs ConfigsModel) validate() error {
	// required
	if inputs.WebhookURL == "" {
		return fmt.Errorf("No Webhook URL parameter specified!")
	}
	if inputs.Message == "" {
		return fmt.Errorf("No Message parameter specified!")
	}
	return nil
}

// RequestParams ...
type RequestParams struct {
	// - required
	Text string `json:"text"`
	// - optional
	Channel   *string `json:"channel"`
	Username  *string `json:"username"`
	EmojiIcon *string `json:"icon_emoji"`
	IconURL   *string `json:"icon_url"`
}

// CreatePayloadParam ...
func CreatePayloadParam(configs ConfigsModel) (string, error) {
	// - required
	reqParams := RequestParams{
		Text: configs.Message,
	}
	if configs.IsBuildFailed {
		failedMsg := configs.MessageOnError
		if failedMsg == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no message_on_error defined, using default."))
		} else {
			reqParams.Text = failedMsg
		}
	}

	// - optional
	reqChannel := configs.Channel
	if reqChannel != "" {
		reqParams.Channel = &reqChannel
	}
	reqUsername := configs.FromUsername
	if reqUsername != "" {
		reqParams.Username = &reqUsername
	}
	if configs.IsBuildFailed {
		failedUsername := configs.FromUsernameOnError
		if failedUsername == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no from_username_on_error defined, using default."))
		} else {
			reqParams.Username = &failedUsername
		}
	}

	reqEmojiIcon := configs.Emoji
	if reqEmojiIcon != "" {
		reqParams.EmojiIcon = &reqEmojiIcon
	}
	if configs.IsBuildFailed {
		failedEmojiIcon := configs.EmojiOnError
		if failedEmojiIcon == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no emoji_on_error defined, using default."))
		} else {
			reqParams.EmojiIcon = &failedEmojiIcon
		}
	}

	reqIconURL := configs.IconURL
	if reqIconURL != "" {
		reqParams.IconURL = &reqIconURL
	}
	if configs.IsBuildFailed {
		failedIconURL := configs.IconURLOnError
		if failedIconURL == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no icon_url_on_error defined, using default."))
		} else {
			reqParams.IconURL = &failedIconURL
		}
	}
	// if Icon URL defined ignore the emoji input
	if reqParams.IconURL != nil {
		reqParams.EmojiIcon = nil
	}

	if configs.IsDebugMode {
		fmt.Printf("Parameters: %#v\n", reqParams)
	}

	// JSON serialize the request params
	reqParamsJSONBytes, err := json.Marshal(reqParams)
	if err != nil {
		return "", nil
	}
	reqParamsJSONString := string(reqParamsJSONBytes)

	return reqParamsJSONString, nil
}

func main() {
	configs := createConfigsModelFromEnvs()
	configs.print()
	if err := configs.validate(); err != nil {
		fmt.Println()
		fmt.Println(colorstring.Red("Issue with input:"), err)
		fmt.Println()
		os.Exit(1)
	}

	//
	// request URL
	requestURL := configs.WebhookURL

	//
	// request parameters
	reqParamsJSONString, err := CreatePayloadParam(configs)
	if err != nil {
		fmt.Println(colorstring.Red("Failed to create JSON payload:"), err)
		os.Exit(1)
	}
	if configs.IsDebugMode {
		fmt.Println()
		fmt.Println("JSON payload: ", reqParamsJSONString)
	}

	//
	// send request
	resp, err := http.PostForm(requestURL,
		url.Values{"payload": []string{reqParamsJSONString}})
	if err != nil {
		fmt.Println(colorstring.Red("Failed to send the request:"), err)
		os.Exit(1)
	}

	//
	// process the response
	body, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println()
		fmt.Println(colorstring.Red("Request failed"))
		fmt.Println("Response from Slack: ", bodyStr)
		fmt.Println()
		os.Exit(1)
	}

	if configs.IsDebugMode {
		fmt.Println()
		fmt.Println("Response from Slack: ", bodyStr)
	}
	fmt.Println()
	fmt.Println(colorstring.Green("Slack message successfully sent! 🚀"))
	fmt.Println()
	os.Exit(0)
}
