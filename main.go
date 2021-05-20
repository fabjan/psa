package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/fabjan/psa/announce"
)

// Announcer announces messages for the public
type Announcer interface {
	Announce(string) error
}

func main() {

	flagMsg := flag.String("message", "", "the message to announce")
	flagDryRun := flag.Bool("dryrun", false, "validate and log but don't send anything")
	flagVerbose := flag.Bool("v", false, "verbose logging")

	flag.Parse()

	announcement := renderAnnouncement(*flagMsg)

	log.Printf("PSA: %s\n", announcement)

	announcers := configureAnnouncers()

	if *flagVerbose {
		log.Printf("%d announcers configured\n", len(announcers))
	}

	if *flagDryRun {
		log.Printf("dry run complete, exiting\n")
		return
	}

	numErrors := 0
	for _, a := range announcers {
		if *flagVerbose {
			log.Printf("[%s] announcing ...\n", a)
		}

		err := a.Announce(announcement)
		if err != nil {
			numErrors++
			log.Printf("[%s] annouce failed: %v\n", a, err)
		} else if *flagVerbose {
			log.Printf("[%s] OK!\n", a)
		}
	}

	if numErrors == 0 {
		log.Printf("all announcers succeeded\n")
	} else {
		log.Printf("%d announcers reported errors\n", numErrors)
	}
}

func mustGetURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("can't parse URL %s: %v\n", rawURL, err)
	}
	return u
}

func renderAnnouncement(m string) string {
	rawTmpl := os.Getenv("PSA_MSG_TEMPLATE")
	if rawTmpl == "" {
		rawTmpl = "📣 {{.Message}}"
	}
	if !strings.Contains(rawTmpl, ".Message") {
		// this probably won't catch all issues
		log.Fatalf("$PSA_MSG_TEMPLATE has no '.Message'")
	}

	tmpl := template.Must(template.New("message").Parse(rawTmpl))

	var buf bytes.Buffer
	tmpl.Execute(&buf, struct{ Message string }{
		Message: m,
	})

	return template.HTMLEscapeString(buf.String())
}

func configureAnnouncers() []Announcer {
	anns := []Announcer{}

	dWH := os.Getenv("PSA_DISCORD_WEBHOOK")
	if dWH != "" {
		u := mustGetURL(dWH)
		log.Printf("Discord webhook detected\n")
		anns = append(anns, announce.DiscordHook(u))
	}

	sWH := os.Getenv("PSA_SLACK_WEBHOOK")
	if sWH != "" {
		u := mustGetURL(sWH)
		log.Printf("Slack webhook detected\n")
		anns = append(anns, announce.SlackHook(u))
	}

	return anns
}
