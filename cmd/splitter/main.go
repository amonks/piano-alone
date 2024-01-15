package main

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	f, err := smf.ReadFile("example.mid")
	if err != nil {
		return err
	}
	dst := abstrack.New()
	var channel, key, velocity uint8
	for _, ev := range abstrack.FromSMF(f.Tracks[0]).Events {
		switch ev.Message.Type() {
		case midi.NoteOnMsg:
			ev.Message.GetNoteOn(&channel, &key, &velocity)
			if key != 61 && key != 56 && key != 64 {
				continue
			}
		case midi.NoteOffMsg:
			ev.Message.GetNoteOff(&channel, &key, &velocity)
			if key != 61 && key != 56 && key != 64 {
				continue
			}
		}
		dst.Events = append(dst.Events, ev)
	}
	out := smf.New()
	out.Add(dst.ToSMF())
	if err := out.WriteFile("split.mid"); err != nil {
		return err
	}
	return nil
}
