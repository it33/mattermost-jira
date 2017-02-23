package jira

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	// DefaultIconURL w
	DefaultIconURL = "https://raw.githubusercontent.com/csduarte/mattermost-jira/master/assets/logo.png"
	// DefaultUsername w
	DefaultUsername = "JIRA"
)

// Bridge sturcture will hold Jira Bridge data and settings
type Bridge struct {
	UsernameOverride string
	IconURL          string
}

// NewBridge generates a default bridge
func NewBridge() *Bridge {
	return &Bridge{
		UsernameOverride: DefaultUsername,
		IconURL:          DefaultIconURL,
	}
}

// Handler will return the handler for use any ServerMux
func (b *Bridge) Handler(w http.ResponseWriter, r *http.Request) {
	mattermostHookURL := r.URL.Query().Get("mattermost_hook_url")

	if len(mattermostHookURL) < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request\n"))
		return
	}

	message, err := NewMessageFromRequest(b, r)
	if err != nil {
		// error
	}

	data, err := message.toJSON()
	if err != nil {
		// error
	}

	req, _ := http.NewRequest("POST", mattermostHookURL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK\n"))
}
