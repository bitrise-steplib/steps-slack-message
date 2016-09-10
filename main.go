package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

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
func CreatePayloadParam(isBuildFailedMode bool) (string, error) {
	// - required
	reqParams := RequestParams{
		Text: os.Getenv("message"),
	}
	if isBuildFailedMode {
		failedMsg := os.Getenv("message_on_error")
		if failedMsg == "" {
			fmt.Println(" (i) Build failed but no message_on_error defined, using default.")
		} else {
			reqParams.Text = failedMsg
		}
	}

	// - optional
	reqChannel := os.Getenv("channel")
	if reqChannel != "" {
		reqParams.Channel = &reqChannel
	}
	reqUsername := os.Getenv("from_username")
	if reqUsername != "" {
		reqParams.Username = &reqUsername
	}
	if isBuildFailedMode {
		failedUsername := os.Getenv("from_username_on_error")
		if failedUsername == "" {
			fmt.Println(" (i) Build failed but no from_username_on_error defined, using default.")
		} else {
			reqParams.Username = &failedUsername
		}
	}

	reqEmojiIcon := os.Getenv("emoji")
	if reqEmojiIcon != "" {
		reqParams.EmojiIcon = &reqEmojiIcon
	}
	if isBuildFailedMode {
		failedEmojiIcon := os.Getenv("emoji_on_error")
		if failedEmojiIcon == "" {
			fmt.Println(" (i) Build failed but no emoji_on_error defined, using default.")
		} else {
			reqParams.EmojiIcon = &failedEmojiIcon
		}
	}

	reqIconURL := os.Getenv("icon_url")
	if reqIconURL != "" {
		reqParams.IconURL = &reqIconURL
	}
	if isBuildFailedMode {
		failedIconURL := os.Getenv("icon_url_on_error")
		if failedIconURL == "" {
			fmt.Println(" (i) Build failed but no icon_url_on_error defined, using default.")
		} else {
			reqParams.IconURL = &failedIconURL
		}
	}
	// if Icon URL defined ignore the emoji input
	if reqParams.IconURL != nil {
		reqParams.EmojiIcon = nil
	}

	fmt.Printf("Parameters: %#v\n", reqParams)

	// JSON serialize the request params
	reqParamsJSONBytes, err := json.Marshal(reqParams)
	if err != nil {
		return "", nil
	}
	reqParamsJSONString := string(reqParamsJSONBytes)

	return reqParamsJSONString, nil
}

func main() {
	//
	// request URL
	requestURL := os.Getenv("webhook_url")
	fmt.Println("URL: ", requestURL)

	isBuildFailedMode := (os.Getenv("STEPLIB_BUILD_STATUS") != "0")

	//
	// request parameters
	reqParamsJSONString, err := CreatePayloadParam(isBuildFailedMode)
	if err != nil {
		fmt.Println("Failed to create JSON payload: ", err)
		os.Exit(1)
	}
	fmt.Println("JSON payload: ", reqParamsJSONString)

	//
	// send request
	resp, err := http.PostForm(requestURL,
		url.Values{"payload": []string{reqParamsJSONString}})
	if err != nil {
		fmt.Printf("Failed to send the request: %s", err)
		os.Exit(1)
	}

	//
	// process the response
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	bodyStr := string(body)
	fmt.Println("Response: ", bodyStr)

	if resp.StatusCode != 200 {
		os.Exit(1)
	}

	os.Exit(0)
}
