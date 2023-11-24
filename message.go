package main

import (
	"encoding/json"
	"strings"
)

// Message to post to a slack channel.
// See also: https://api.slack.com/methods/chat.postMessage
type Message struct {
	// Channel to send message to.
	//
	// Can be an encoded ID (eg. C024BE91L), or the channel's name (eg. #general).
	Channel string `json:"channel"`

	// Text of the message to send. Required, unless providing only attachments instead.
	Text string `json:"text,omitempty"`

	// Attachments is a list of structured attachments.
	Attachments []Attachment `json:"attachments,omitempty"`

	// IconEmoji is the emoji to use as the icon for the message. Overrides IconUrl.
	IconEmoji string `json:"icon_emoji,omitempty"`

	// IconURL is the URL to an image to use as the icon for the message.
	IconURL string `json:"icon_url,omitempty"`

	// LinkNames linkifies channel names and usernames.
	LinkNames bool `json:"link_names,omitempty"`

	// Username specifies the bot's username for the message.
	Username string `json:"username,omitempty"`

	// Provide another message's ts value to make this message a reply.
	ThreadTs string `json:"thread_ts,omitempty"`

	// Provide another message's ts value to make message to update
	Ts string `json:"ts,omitempty"`

	// Used in conjunction with thread_ts and indicates whether reply should be made visible to everyone in the channel or conversation.
	ReplyBroadcast bool `json:"reply_broadcast,omitempty"`
}

// Attachment adds more context to a slack chat message.
// See also: https://api.slack.com/docs/message-attachments
type Attachment struct {
	// Fallback is the plain-text summary of the attachment.
	//
	// This text will be used in clients that don't show formatted text (eg. IRC, mobile notifications)
	// and should not contain any markup.
	Fallback string `json:"fallback"`

	// Color is used to color the border along the left side of the attachment.
	//
	// Can either be one of good, warning, danger, or any hex color code (eg. #439FE0).
	Color string `json:"color"`

	// PreText is an optional text that appears above the attachment block.
	PreText string `json:"pretext,omitempty"`

	// AuthorName is a small text used to display the author's name.
	AuthorName string `json:"author_name,omitempty"`

	// Title is displayed as larger, bold text near the top of a attachment.
	Title string `json:"title,omitempty"`

	// TitleLink is a URL that will hyperlink the Title.
	TitleLink string `json:"title_link,omitempty"`

	// Text is the main text of the attachment, and can contain standard message markup.
	//
	// The content will automatically collapse if it contains 700+ characters or 5+ linebreaks,
	// and will display a "Show more..." link to expand the content.
	Text string `json:"text,omitempty"`

	// Fields is a list of fields to be displayed in a table inside the attachment.
	Fields []Field `json:"fields,omitempty"`

	// ImageURL is a URL to an image file that will be displayed inside the attachment.
	//
	// Supported formats: GIF, JPEG, PNG, and BMP.
	// Large images will be resized to a maximum width of 400px or a maximum height of 500px.
	ImageURL string `json:"image_url,omitempty"`

	// ThumbURL is a URL to an image file that will be displayed as a
	// thumbnail on the right side of a attachment.
	//
	// Supported formats: GIF, JPEG, PNG, and BMP.
	// The thumbnail's longest dimension will be scaled down to 75px.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Footer adds some brief text to help contextualize and identify an attachment.
	//
	// Limited to 300 characters.
	Footer string `json:"footer,omitempty"`

	// FooterIcon renders a small icon beside the footer text.
	//
	// It will be scaled down to 16px by 16px.
	FooterIcon string `json:"footer_icon,omitempty"`

	// TimeStamp is an integer value in epoch time to display and additional
	// timestamp value as part of the attachment's footer.
	TimeStamp int `json:"ts,omitempty"`

	// Buttons is a list of buttons attached to the message as link buttons.
	//
	// An attachment may contain 1 to 5 buttons.
	Buttons []Button `json:"actions,omitempty"`
}

// Field will be displayed in a table inside the attachment.
type Field struct {
	// Title is shown as a bold heading above the value text.
	Title string

	// Value is the text value of the field.
	Value string

	// Short is an optional flag indicating whether the value is short enough
	// to be displayed side-by-side with other values.
	Short bool
}

// MarshalJSON implements json.Marshaler.MarshalJSON.
func (f Field) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["title"] = f.Title
	m["value"] = f.Value
	m["short"] = len(f.Value) < 40
	return json.Marshal(m)
}

func parseFields(s string) (fs []Field) {
	for _, p := range pairs(s) {
		fs = append(fs, Field{Title: p[0], Value: ensureNewlines(p[1])})
	}
	return
}

// Button is just a link that looks like a button.
type Button struct {
	// Type is set to button to tell slack to render a button.
	Type string

	// Text is the label for the button.
	Text string

	// URL is the fully qualified http or https url to deliver users to.
	URL string

	// Style is set to default so the buttons will use the UI's default text color.
	Style string
}

// MarshalJSON implements json.Marshaler.MarshalJSON.
func (b Button) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	m["type"] = "button"
	m["text"] = b.Text
	m["url"] = b.URL
	m["style"] = "default"
	return json.Marshal(m)
}

func parseButtons(s string) (bs []Button) {
	for _, p := range pairs(s) {
		bs = append(bs, Button{Text: p[0], URL: p[1]})
	}
	return
}

// pairs slices every lines in s into two substrings separated by the first pipe
// character and returns a slice of those pairs.
func pairs(s string) [][2]string {
	var ps [][2]string
	for _, line := range strings.Split(s, "\n") {
		a := strings.SplitN(line, "|", 2)
		if len(a) == 2 && a[0] != "" && a[1] != "" {
			ps = append(ps, [2]string{a[0], a[1]})
		}
	}
	return ps
}
