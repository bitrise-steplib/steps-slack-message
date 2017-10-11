package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"strings"

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
	Color               string
	ColorOnError        string
	ImageURL            string
	ImageURLOnError     string
	Emoji               string
	EmojiOnError        string
	IconURL             string
	IconURLOnError      string
	IsLinkNames         bool
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
		Color:               os.Getenv("color"),
		ColorOnError:        os.Getenv("color_on_error"),
		ImageURL:            os.Getenv("image_url"),
		ImageURLOnError:     os.Getenv("image_url_on_error"),
		IconURL:             os.Getenv("icon_url"),
		IconURLOnError:      os.Getenv("icon_url_on_error"),
		IsLinkNames:         os.Getenv("link_names") == "yes",
		//
		IsDebugMode: (os.Getenv("is_debug_mode") == "yes"),
		//
		IsBuildFailed: (os.Getenv("STEPLIB_BUILD_STATUS") != "0"),
	}
}

func (configs ConfigsModel) print() {
	fmt.Println("")
	fmt.Println(colorstring.Blue("Slack configs:"))
	fmt.Println(" - WebhookURL:", configs.WebhookURL)
	fmt.Println(" - Channel:", configs.Channel)
	fmt.Println(" - FromUsername:", configs.FromUsername)
	fmt.Println(" - FromUsernameOnError:", configs.FromUsernameOnError)
	fmt.Println(" - Message:", configs.Message)
	fmt.Println(" - MessageOnError:", configs.MessageOnError)
	fmt.Println(" - Color:", configs.Color)
	fmt.Println(" - ColorOnError:", configs.ColorOnError)
	fmt.Println(" - ImageURL:", configs.ImageURL)
	fmt.Println(" - ImageURLOnError:", configs.ImageURLOnError)
	fmt.Println(" - Emoji:", configs.Emoji)
	fmt.Println(" - EmojiOnError:", configs.EmojiOnError)
	fmt.Println(" - IconURL:", configs.IconURL)
	fmt.Println(" - IconURLOnError:", configs.IconURLOnError)
	fmt.Println(" - IsLinkNames:", configs.IsLinkNames)
	fmt.Println("")
	fmt.Println(colorstring.Blue("Other configs:"))
	fmt.Println(" - IsDebugMode:", configs.IsDebugMode)
	fmt.Println(" - IsBuildFailed:", configs.IsBuildFailed)
	fmt.Println("")
}

func (configs ConfigsModel) validate() error {
	// required
	if configs.WebhookURL == "" {
		return errors.New("No Webhook URL parameter specified")
	}
	if configs.Message == "" {
		return errors.New("No Message parameter specified")
	}
	if configs.Color == "" {
		return errors.New("No Color parameter specified")
	}

	return nil
}

// AttachmentItemModel ...
type AttachmentItemModel struct {
	Fallback string   `json:"fallback"`
	Text     string   `json:"text"`
	Color    string   `json:"color,omitempty"`
	ImageURL string   `json:"image_url"`
	MrkdwnIn []string `json:"mrkdwn_in,omitempty"`
}

// RequestParams ...
type RequestParams struct {
	// - required
	Text string `json:"text"`
	// OR use attachment instead of text, for better formatting
	Attachments []AttachmentItemModel `json:"attachments,omitempty"`
	// - optional
	Channel   *string `json:"channel"`
	Username  *string `json:"username"`
	EmojiIcon *string `json:"icon_emoji"`
	IconURL   *string `json:"icon_url"`
	LinkNames int     `json:"link_names"`
}

// ensureNewlineEscapeChar replaces the "\" + "n" char sequences with the "\n" newline char
func ensureNewlineEscapeChar(s string) string {
	return strings.Replace(s, "\\"+"n", "\n", -1)
}

// CreatePayloadParam ...
func CreatePayloadParam(configs ConfigsModel) (string, error) {
	// - required
	msgColor := configs.Color
	if configs.IsBuildFailed {
		if configs.ColorOnError == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no color_on_error defined, using default."))
		} else {
			msgColor = configs.ColorOnError
		}
	}
	msgText := configs.Message
	if configs.IsBuildFailed {
		if configs.MessageOnError == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no message_on_error defined, using default."))
		} else {
			msgText = configs.MessageOnError
		}
	}
	msgText = ensureNewlineEscapeChar(msgText)
	// - optional attachment params
	msgImage := configs.ImageURL
	if configs.IsBuildFailed {
		if configs.ImageURLOnError == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no image_url_on_error defined, using default."))
		} else {
			msgImage = configs.ImageURLOnError
		}
	}

	reqParams := RequestParams{
		Attachments: []AttachmentItemModel{
			{
				Text: msgText, Fallback: msgText,
				Color:    msgColor,
				ImageURL: msgImage,
				MrkdwnIn: []string{"text", "pretext", "fields"},
			},
		},
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
		if configs.FromUsernameOnError == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no from_username_on_error defined, using default."))
		} else {
			reqParams.Username = &configs.FromUsernameOnError
		}
	}

	reqEmojiIcon := configs.Emoji
	if reqEmojiIcon != "" {
		reqParams.EmojiIcon = &reqEmojiIcon
	}
	if configs.IsBuildFailed {
		if configs.EmojiOnError == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no emoji_on_error defined, using default."))
		} else {
			reqParams.EmojiIcon = &configs.EmojiOnError
		}
	}

	reqIconURL := configs.IconURL
	if reqIconURL != "" {
		reqParams.IconURL = &reqIconURL
	}
	if configs.IsBuildFailed {
		if configs.IconURLOnError == "" {
			fmt.Println(colorstring.Yellow(" (i) Build failed but no icon_url_on_error defined, using default."))
		} else {
			reqParams.IconURL = &configs.IconURLOnError
		}
	}
	// if Icon URL defined ignore the emoji input
	if reqParams.IconURL != nil {
		reqParams.EmojiIcon = nil
	}

	if configs.IsLinkNames {
		reqParams.LinkNames = 1
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

	if resp.StatusCode != 200 || bodyStr != "ok" {
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
