package main

import (
	"fmt"

	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/songs"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	f := songs.PreludeBergamasqueSMF
	counts := abstrack.FromSMF(f, 0).CountNotes()
	for _, key := range counts {
		fmt.Println(key.Key, key.Count)
	}
	fmt.Println(len(counts), "keys in total")
	return nil
}
