package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/songs"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	f := songs.PreludeBergamasqueSMF

	var tr smf.Track
	activeNotes := 0
	counter := 0
	var c, key, vel uint8
	for _, ev := range f.Tracks[0] {
		tr.Add(ev.Delta, ev.Message)
		if ev.Message.GetNoteStart(&c, &key, &vel) {
			activeNotes += 1
		} else if ev.Message.GetNoteEnd(&c, &key) {
			activeNotes -= 1
			if activeNotes == 0 {
				counter += 1
				if counter == 15 {
					fmt.Println("enough")
					break
				}
			}
		}
	}
	tr.Close(0)

	out := smf.New()
	out.TimeFormat = f.TimeFormat
	out.Add(tr)
	if err := out.WriteFile("excerpt.mid"); err != nil {
		return err
	}
	return nil
}
