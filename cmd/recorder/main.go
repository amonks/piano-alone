package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"gitlab.com/gomidi/midi/v2"
	"monks.co/piano-alone/recorder"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Println("done")
}

var (
	note     = flag.Int("note", 100, "note")
	filename = flag.String("filename", "out.mid", "filename")
)

func run() error {
	flag.Parse()

	r := recorder.New(120)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		r.Record(recorder.Now(midi.NoteOn(1, uint8(*note), 100)))
	}
	r.Close()

	bs, err := r.Bytes()
	if err != nil {
		return err
	}

	if err := os.WriteFile(*filename, bs, 0666); err != nil {
		return err
	}

	return nil
}
