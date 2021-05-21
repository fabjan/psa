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

package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"

	"github.com/fabjan/psa/configure"
)

func main() {

	flagMsg := flag.String("message", "", "the message to announce")
	flagDryRun := flag.Bool("dryrun", false, "validate and log but don't send anything")
	flagVerbose := flag.Bool("v", false, "verbose logging")

	flag.Parse()

	cfg, err := configure.FromEnv()

	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	announcement := renderAnnouncement(cfg.MessageTemplate, *flagMsg)

	log.Printf("PSA: %s\n", announcement)

	announcers := cfg.Announcers()

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

func renderAnnouncement(tmpl *template.Template, msg string) string {
	var buf bytes.Buffer
	tmpl.Execute(&buf, struct{ Message string }{
		Message: msg,
	})
	return template.HTMLEscapeString(buf.String())
}
