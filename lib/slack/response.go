package slack

// SendMessageResponse is the response from Slack POST
type SendMessageResponse struct {
	/// The Thread Timestamp
	Timestamp string `json:"ts"`
}
