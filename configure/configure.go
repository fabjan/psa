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

package configure

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/fabjan/psa/announce"
)

// AppConfig contains the information needed to setup PSA
type AppConfig struct {
	DiscordHookURL  *url.URL
	SlackHookURL    *url.URL
	MessageTemplate *template.Template
}

// FromEnv creates a configuration from environment variables
func FromEnv() (cfg AppConfig, err error) {
	dWH := os.Getenv("PSA_DISCORD_WEBHOOK")
	sWH := os.Getenv("PSA_SLACK_WEBHOOK")
	rawTmpl := os.Getenv("PSA_MSG_TEMPLATE")

	return FromStrings(dWH, sWH, rawTmpl)
}

// FromStrings creates a configuration based on the given values
func FromStrings(discordWH, slackWH, msgTmpl string) (cfg AppConfig, err error) {

	if discordWH != "" {
		cfg.DiscordHookURL, err = parseHTTPURL(discordWH)
		if err != nil {
			return cfg, err
		}
		log.Printf("Discord webhook configured\n")
	}

	if slackWH != "" {
		cfg.SlackHookURL, err = parseHTTPURL(slackWH)
		if err != nil {
			return cfg, err
		}
		log.Printf("Slack webhook configured\n")
	}

	if msgTmpl == "" {
		msgTmpl = "📣 {{.Message}}"
	}
	if !strings.Contains(msgTmpl, ".Message") {
		// this probably won't catch all issues
		return cfg, errors.New("configured message template has no '.Message'")
	}

	tmpl, err := template.New("message").Parse(msgTmpl)

	if err != nil {
		return cfg, err
	}

	cfg.MessageTemplate = tmpl

	return cfg, nil
}

// Announcers creates the announcers from the configuration
func (cfg *AppConfig) Announcers() []announce.Announcer {
	anns := []announce.Announcer{}

	if cfg.DiscordHookURL != nil {
		anns = append(anns, announce.DiscordHook(cfg.DiscordHookURL))
	}

	if cfg.SlackHookURL != nil {
		anns = append(anns, announce.SlackHook(cfg.SlackHookURL))
	}

	return anns
}

func parseHTTPURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("can't parse URL %s: %v", rawURL, err)
	}
	if !strings.HasPrefix(u.Scheme, "http") {
		return nil, fmt.Errorf("can't parse URL %s: %v", rawURL, err)
	}
	return u, nil
}
