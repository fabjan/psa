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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// N.B this is very similar to the Discordian code,
// but I'm not sure yet how much they will diverge so
// the webhooking will not be abstracted yet

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
	req, err := s.makeRequest(m)
	if err != nil {
		return fmt.Errorf("failed Slack marshalling: %w", err)
	}
	err = handleResponse(http.DefaultClient.Do(req))
	if err != nil {
		return fmt.Errorf("failed Slack announce: %w", err)
	}
	return nil
}

func (*Slacker) String() string {
	return "Slack Announcer"
}

func (s *Slacker) makeRequest(msg string) (*http.Request, error) {
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(map[string]string{"text": msg})
	if err != nil {
		return nil, fmt.Errorf("failed Slack payload marshalling: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.url.String(), &payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}
