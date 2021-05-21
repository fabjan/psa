//	Copyright 2021 Fabian Bergström
//	
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//	
//			http://www.apache.org/licenses/LICENSE-2.0
//	
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.

package announce

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// N.B this is very similar to the Slacker code,
// but I'm not sure yet how much they will diverge so
// the webhooking will not be abstracted yet

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
	req, err := d.makeRequest(m)
	if err != nil {
		return fmt.Errorf("failed Discord marshalling: %w", err)
	}
	err = handleResponse(http.DefaultClient.Do(req))
	if err != nil {
		return fmt.Errorf("failed Discord announce: %w", err)
	}
	return nil
}

func (*Discordian) String() string {
	return "Discord Announcer"
}

func (d *Discordian) makeRequest(msg string) (*http.Request, error) {
	payload := url.Values{
		"content": []string{msg},
	}
	body := strings.NewReader(payload.Encode())
	req, err := http.NewRequest(http.MethodPost, d.url.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}
