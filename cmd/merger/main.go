package main

import (
	"flag"
	"fmt"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Println("done")
}

var (
	filenameA   = flag.String("a", "", "filename a")
	filenameB   = flag.String("b", "", "filename b")
	filenameOut = flag.String("out", "merged.mid", "output filename")
)

func run() error {
	flag.Parse()

	smfA, err := smf.ReadFile(*filenameA)
	if err != nil {
		return fmt.Errorf("error reading '%s': %w", *filenameA, err)
	}
	smfB, err := smf.ReadFile(*filenameB)
	if err != nil {
		return fmt.Errorf("error reading '%s': %w", *filenameB, err)
	}

	mergedTrack := abstrack.FromSMF(smfA.Tracks[0]).Merge(abstrack.FromSMF(smfB.Tracks[0]))
	merged := smf.New()
	merged.Add(mergedTrack.ToSMF())
	if err := merged.WriteFile(*filenameOut); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}
