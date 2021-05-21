package configure

import (
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
		cfg.DiscordHookURL = mustGetURL(discordWH)
		log.Printf("Discord webhook detected\n")
	}

	if slackWH != "" {
		cfg.SlackHookURL = mustGetURL(slackWH)
		log.Printf("Discord webhook detected\n")
	}

	if msgTmpl == "" {
		msgTmpl = "📣 {{.Message}}"
	}
	if !strings.Contains(msgTmpl, ".Message") {
		// this probably won't catch all issues
		log.Fatalf("$PSA_MSG_TEMPLATE has no '.Message'")
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

func mustGetURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("can't parse URL %s: %v\n", rawURL, err)
	}
	return u
}
