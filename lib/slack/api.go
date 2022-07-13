package slack

type SlackApi interface {
	Post(msg *Message) (SendMessageResponse, error)
}
