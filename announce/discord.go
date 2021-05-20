package announce

import (
	"fmt"
	"net/http"
	"net/url"
)

// Discordian can announce to Discord
type Discordian struct {
	url *url.URL
}

// DiscordHook creates a Discordian using the given webhook for announcements
func DiscordHook(url *url.URL) *Discordian {
	return &Discordian{url}
}

// Announce sends the given message to the configured channel(s)
func (d *Discordian) Announce(m string) error {
	payload := url.Values{
		"content": []string{m},
	}
	resp, err := http.DefaultClient.PostForm(d.url.String(), payload)
	err = handleResponse(resp, err)
	if err != nil {
		return fmt.Errorf("failed Discord announce: %w", err)
	}
	return nil
}

func (*Discordian) String() string {
	return "Discord Announcer"
}
