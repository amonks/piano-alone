package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

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

	c := make(chan midi.Message)

	r := recorder.New()
	go r.Record(120, c)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		c <- midi.NoteOn(1, uint8(*note), 100)
	}
	close(c)
	time.Sleep(time.Second)

	bs, err := r.Bytes()
	if err != nil {
		return err
	}

	if err := os.WriteFile(*filename, bs, 0666); err != nil {
		return err
	}

	return nil
}
