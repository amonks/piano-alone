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
	counts := abstrack.FromSMF(f, 0).CountNotes()
	for _, key := range counts {
		fmt.Println(key.Key, key.Count)
	}
	fmt.Println(len(counts), "keys in total")
	return nil
}
