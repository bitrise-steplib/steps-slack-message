package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type RequestParams struct {
	// - required
	Text string `json:"text"`
	// - optional
	Channel   *string `json:"channel"`
	Username  *string `json:"username"`
	EmojiIcon *string `json:"icon_emoji"`
	IconURL   *string `json:"icon_url"`
}

func CreatePayloadParam() (string, error) {
	// - required
	reqParams := RequestParams{
		Text: os.Getenv("SLACK_MESSAGE_TEXT"),
	}
	// - optional
	reqChannel := os.Getenv("SLACK_CHANNEL")
	if reqChannel != "" {
		reqParams.Channel = &reqChannel
	}
	reqUsername := os.Getenv("SLACK_FROM_NAME")
	if reqUsername != "" {
		reqParams.Username = &reqUsername
	}
	reqEmojiIcon := os.Getenv("SLACK_ICON_EMOJI")
	if reqEmojiIcon != "" {
		reqParams.EmojiIcon = &reqEmojiIcon
	}
	reqIconURL := os.Getenv("SLACK_ICON_URL")
	if reqIconURL != "" {
		reqParams.IconURL = &reqIconURL
	}
	fmt.Printf("Parameters: %#v\n", reqParams)

	// JSON serialize the request params
	reqParamsJsonBytes, err := json.Marshal(reqParams)
	if err != nil {
		return "", nil
	}
	reqParamsJsonString := string(reqParamsJsonBytes)

	return reqParamsJsonString, nil
}

func main() {
	//
	// request URL
	requestURL := os.Getenv("SLACK_WEBHOOK_URL")
	fmt.Println("URL: ", requestURL)

	//
	// request parameters
	reqParamsJsonString, err := CreatePayloadParam()
	if err != nil {
		fmt.Println("Failed to create JSON payload: ", err)
		os.Exit(1)
	}
	fmt.Println("JSON payload: ", reqParamsJsonString)

	//
	// send request
	resp, err := http.PostForm(requestURL,
		url.Values{"payload": []string{reqParamsJsonString}})
	if err != nil {
		fmt.Printf("Failed to send the request", err)
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
