package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Message structure for Mattermost JSON creation.
type Message struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
}

// NewMessageFromRequest w
func NewMessageFromRequest(b *Bridge, r *http.Request) (*Message, error) {

	channelOverride := r.URL.Query().Get("channel")

	hook, err := NewWebhookfromJSON(r.Body)
	if err != nil {
		return nil, err
	}

	m := &Message{
		Text:     hook.String(),
		Channel:  channelOverride,
		Username: b.UsernameOverride,
		IconURL:  b.IconURL,
	}
	fmt.Println(m)
	return m, nil
}

func (m *Message) toJSON() ([]byte, error) {
	return json.Marshal(m)
}
