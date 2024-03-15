package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"monks.co/piano-alone/baseurl"
)

var (
	fBaseURL = flag.String("baseURL", "https://piano.computer", "base server url")

	menu = []string{
		"Performance Status",
		"MIDI Configuration",
		"MIDI Output Test",
		"Message Log",
	}
)

func main() {
	flag.Parse()
	defer midi.CloseDriver()
	p := tea.NewProgram(model{baseURL: baseurl.From(*fBaseURL)}, tea.WithAltScreen())
	if m, err := p.Run(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(m.(model).output)
	}
}
