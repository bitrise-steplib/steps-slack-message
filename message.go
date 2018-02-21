package main

import (
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

	// AuthorLink is a valid URL that will hyperlink the AuthorName.
	// AuthorLink string `json:"author_link,omitempty"`

	// AuthorIcon is a valid URL that displays a small 16x16px image to the left of the AuthorName.
	// AuthorIcon string `json:"author_icon,omitempty"`

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
	Title string `json:"title"`

	// Value is the text value of the field.
	Value string `json:"value"`

	// Short is an optional flag indicating whether the value is short enough
	// to be displayed side-by-side with other values.
	Short bool `json:"short"`
}

// Button is just a link that looks like a button.
type Button struct {
	// Type is set to button to tell slack to render a button.
	Type string `json:"type"`

	// Text is the label for the button.
	Text string `json:"text"`

	// URL is the fully qualified http or https url to deliver users to.
	URL string `json:"url"`

	// Style is set to default so the buttons will use the UI's default text color.
	Style string `json:"style"`
}

func parseFields(s string) []Field {
	var fields []Field
	for _, line := range strings.Split(s, "\n") {
		a := strings.SplitN(line, "|", 2)
		if len(a) == 2 && a[0] != "" && a[1] != "" {
			fields = append(fields, Field{
				Title: a[0],
				Value: a[1],
				Short: len(a[1]) < 20,
			})
		}
	}
	return fields
}

func parseButtons(s string) []Button {
	var buttons []Button
	for _, line := range strings.Split(s, "\n") {
		a := strings.SplitN(line, "|", 2)
		if len(a) == 2 && a[0] != "" && a[1] != "" {
			buttons = append(buttons, Button{
				Type:  "button",
				Text:  a[0],
				URL:   a[1],
				Style: "default",
			})
		}
	}
	return buttons
}
