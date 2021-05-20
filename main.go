package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"net/url"
	"os"

	"github.com/fabjan/psa/announce"
)

// Announcer announces messages for the public
type Announcer interface {
	Announce(string) error
}

func main() {

	flagMsg := flag.String("message", "", "the message to announce")
	flagMsgTmpl := flag.String("template", "📣 {{.Message}}", "template to apply when announcing")
	flagVerbose := flag.Bool("v", false, "verbose logging")

	flag.Parse()

	tpl := template.Must(template.New("message").Parse(*flagMsgTmpl))

	anns := []Announcer{}

	dWH := os.Getenv("DISCORD_WEBHOOK")
	if dWH != "" {
		u := mustGetURL(dWH)
		log.Printf("Discord web hook detected\n")
		anns = append(anns, announce.DiscordHook(u))
	}

	sWH := os.Getenv("SLACK_WEBHOOK")
	if sWH != "" {
		u := mustGetURL(sWH)
		log.Printf("Slack web hook detected\n")
		anns = append(anns, announce.SlackHook(u))
	}

	if *flagVerbose {
		log.Printf("%d announcers configured\n", len(anns))
	}

	log.Printf("PSA: %q\n", *flagMsg)

	numErrors := 0
	var buf bytes.Buffer
	for _, a := range anns {
		if *flagVerbose {
			log.Printf("[%s] announcing ...\n", a)
		}

		buf.Reset()
		tpl.Execute(&buf, struct{ Message string }{
			Message: *flagMsg,
		})

		msg := template.HTMLEscapeString(buf.String())

		err := a.Announce(msg)
		if err != nil {
			numErrors++
			log.Printf("[%s] annouce failed: %v\n", a, err)
		} else if *flagVerbose {
			log.Printf("[%s] OK!\n", a)
		}
	}

	if numErrors == 0 {
		log.Printf("all announcers succeeded\n")
	}
}

func mustGetURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("can't parse URL %s: %v\n", rawURL, err)
	}
	return u
}
