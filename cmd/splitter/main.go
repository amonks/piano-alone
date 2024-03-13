package main

import (
	"gitlab.com/gomidi/midi/v2/smf"
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
	dst := abstrack.FromSMF(f, 0).Select([]uint8{61, 56, 64})
	out := smf.New()
	out.Add(dst.ToSMF())
	if err := out.WriteFile("split.mid"); err != nil {
		return err
	}
	return nil
}
