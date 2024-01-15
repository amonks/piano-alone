package main

import (
	"fmt"
	"sort"

	"gitlab.com/gomidi/midi/v2/smf"
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
	keyset := map[uint8]int{}
	var channel, key, velocity uint8
	for _, ev := range f.Tracks[0] {
		if !ev.Message.GetNoteOn(&channel, &key, &velocity) {
			continue
		}
		keyset[key] += 1
	}
	var keys []uint8
	for key := range keyset {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keyset[keys[i]] < keyset[keys[j]] })
	for _, key := range keys {
		fmt.Println(key, keyset[key])
	}
	fmt.Println(len(keys), "keys in total")
	return nil
}
