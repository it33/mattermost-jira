package jira

import "encoding/json"

// Message structure for Mattermost JSON creation.
type Message struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
}

// NewMessageFromWebhook w
func NewMessageFromWebhook(w *Webhook, b *Bridge, channel string) *Message {
	m := &Message{
		Text:     w.String(),
		Channel:  channel,
		Username: b.UsernameOverride,
		IconURL:  b.IconURL,
	}

	return m
}

func (m *Message) toJSON() ([]byte, error) {
	return json.Marshal(m)
}
