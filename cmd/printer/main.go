package main

import (
	"fmt"

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
	t := abstrack.FromSMF(f, 0).Select([]uint8{61})
	var c, key, vel uint8
	var bpm float64
	j := 0
	fmt.Println(f.TimeFormat.String())
	for _, ev := range t.Events {
		if ev.Message.GetNoteStart(&c, &key, &vel) {
			fmt.Printf("%d: ON\n", j)
			j++
		} else if ev.Message.GetNoteEnd(&c, &key) {
			fmt.Printf("%d: OFF\n", j)
			j++
		} else if ev.Message.GetMetaTempo(&bpm) {
			fmt.Printf("%f\n", bpm)
		}
	}
	return nil
}
