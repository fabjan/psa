package announce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Slacker can announce to Slack
type Slacker struct {
	url *url.URL
}

// SlackHook creates a Slacker using the given webhook for announcements
func SlackHook(url *url.URL) *Slacker {
	return &Slacker{url}
}

// Announce sends the given message to the configured channel(s)
func (s *Slacker) Announce(m string) error {
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(map[string]string{"text": m})
	if err != nil {
		return fmt.Errorf("failed Slack payload marshalling: %w", err)
	}
	resp, err := http.DefaultClient.Post(s.url.String(), "application/json", &payload)
	err = handleResponse(resp, err)
	if err != nil {
		return fmt.Errorf("failed Slack announce: %w", err)
	}
	return nil
}

func (*Slacker) String() string {
	return "Slack Announcer"
}
