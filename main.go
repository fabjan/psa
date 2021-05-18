package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Announcer announces messages for the public
type Announcer interface {
	Announce(string) error
}

func main() {

	flagMsg := flag.String("msg", "", "the message to announce")
	flagVerbose := flag.Bool("v", false, "verbose logging")

	flag.Parse()

	anns := []Announcer{}

	dWH := os.Getenv("DISCORD_WEBHOOK")
	if dWH != "" {
		u, err := url.Parse(dWH)
		if err != nil {
			log.Fatalf("can't parse webhook URL %s: %v\n", dWH, err)
		}
		log.Printf("Discord web hook detected\n")
		anns = append(anns, &discordAnnouncer{url: u})
	}

	if *flagVerbose {
		log.Printf("%d announcers configured\n", len(anns))
	}

	for _, a := range anns {
		if *flagVerbose {
			log.Printf("announcing with %s\n", a)
		}
		a.Announce(*flagMsg)
	}
}

type discordAnnouncer struct {
	url *url.URL
}

func (da *discordAnnouncer) String() string {
	return "Discord Announcer"
}

func (da *discordAnnouncer) Announce(m string) error {
	payload := url.Values{
		"content": []string{m},
	}
	http.DefaultClient.PostForm(da.url.String(), payload)
	return nil
}
