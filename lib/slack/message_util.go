package slack

import (
	"github.com/bitrise-steplib/steps-slack-message/lib/step"
	"os"
	"strings"
	"time"
)

var BuildIsSuccessful = os.Getenv("BITRISE_BUILD_STATUS") == "0"

// selectValue chooses the right value based on the result of the build.
func selectValue(ifSuccess, ifFailed string) string {
	if BuildIsSuccessful || ifFailed == "" {
		return ifSuccess
	}
	return ifFailed
}

// selectBool chooses the right boolean value based on the result of the build.
func selectBool(ifSuccess, ifFailed bool) bool {
	if BuildIsSuccessful {
		return ifSuccess
	}
	return ifFailed
}

// ensureNewlines replaces all \n substrings with newline characters.
func ensureNewlines(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func parseFields(s string) (fs []Field) {
	for _, p := range pairs(s) {
		fs = append(fs, Field{Title: p[0], Value: p[1]})
	}
	return
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

func NewMessage(c step.Config) Message {
	msg := Message{
		Channel: strings.TrimSpace(selectValue(c.Channel, c.ChannelOnError)),
		Text:    selectValue(c.Text, c.TextOnError),
		Attachments: []Attachment{{
			Fallback:   ensureNewlines(selectValue(c.Message, c.MessageOnError)),
			Color:      selectValue(c.Color, c.ColorOnError),
			PreText:    selectValue(c.PreText, c.PreTextOnError),
			AuthorName: c.AuthorName,
			Title:      selectValue(c.Title, c.TitleOnError),
			TitleLink:  c.TitleLink,
			Text:       ensureNewlines(selectValue(c.Message, c.MessageOnError)),
			Fields:     parseFields(c.Fields),
			ImageURL:   selectValue(c.ImageURL, c.ImageURLOnError),
			ThumbURL:   selectValue(c.ThumbURL, c.ThumbURLOnError),
			Footer:     c.Footer,
			FooterIcon: c.FooterIcon,
			Buttons:    parseButtons(c.Buttons),
		}},
		IconEmoji:      selectValue(c.IconEmoji, c.IconEmojiOnError),
		IconURL:        selectValue(c.IconURL, c.IconURLOnError),
		LinkNames:      c.LinkNames,
		Username:       selectValue(c.Username, c.UsernameOnError),
		ThreadTs:       selectValue(c.ThreadTs, c.ThreadTsOnError),
		ReplyBroadcast: selectBool(c.ReplyBroadcast, c.ReplyBroadcastOnError),
	}
	if c.TimeStamp {
		msg.Attachments[0].TimeStamp = int(time.Now().Unix())
	}
	return msg
}
