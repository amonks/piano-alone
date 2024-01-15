package main

import (
	"flag"
	"fmt"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/merger"
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

	if err := merger.Merge(smfA, smfB); err != nil {
		return fmt.Errorf("error merging: %w", err)
	}
	if err := smfA.WriteFile(*filenameOut); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}
