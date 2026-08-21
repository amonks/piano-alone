package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"monks.co/piano-alone/game"
	"monks.co/piano-alone/id"
	"monks.co/piano-alone/songs"
)

func main() {
	if err := start(); err != nil {
		panic(err)
	}
}

func start() error {
	p := &game.Performance{
		Configuration: &game.Configuration{
			PerformanceID: id.Random128(),
			Title:         "Test",
			Composer:      "Sergei Rachmaninoff",
			Score:         songs.PreludeOpus3No2Bytes,
		},
		Date:       mustParseTime("2006-01-02 15:04:05 MST", "2024-06-11 21:20:00 CDT"),
		IsFeatured: false,
	}

	r := bytes.NewReader(p.Bytes())
	if resp, err := http.Post("https://piano.computer/performances", "audio/midi", r); err != nil {
		return fmt.Errorf("http error: %w", err)
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rejected: %d", resp.StatusCode)
	}
	return nil
}

func mustParseTime(fmt string, t string) time.Time {
	if out, err := time.Parse(fmt, t); err != nil {
		panic(err)
	} else {
		return out
	}
}
