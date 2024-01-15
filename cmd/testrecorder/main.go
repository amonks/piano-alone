package main

import (
	"bufio"
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

func run() error {
	c := make(chan midi.Message)

	r := recorder.New()
	go r.Record(120, c)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		c <- midi.NoteOn(1, 100, 100)
	}
	close(c)
        time.Sleep(time.Second)

	bs, err := r.Bytes()
	if err != nil {
		return err
	}

	if err := os.WriteFile("recorder.mid", bs, 0666); err != nil {
		return err
	}

	return nil
}
