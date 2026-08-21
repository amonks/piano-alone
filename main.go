// Command piano-alone runs Life Online (Piano Telephone): a website
// where an audience, each given a handful of a score's keys, plays a
// piano piece together, and the collected parts are merged into one
// MIDI file for a disklavier to perform.
//
// This is the whole piece in one binary, with no authorization: every
// conductor verb is open to whoever can reach it, which is the right
// default for a laptop in a room and the wrong one for the public
// internet. Serving it publicly means putting a decision in front of
// it — see monks.co/piano-alone/server, whose Operation classifies a
// request for exactly that.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"monks.co/piano-alone/server"
	"monks.co/piano-alone/sigctx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var (
		addr   = flag.String("addr", "0.0.0.0:8080", "address to serve on")
		dbPath = flag.String("db", "piano-alone.db", "path to the SQLite database")
		client = flag.String("client", "", "path to the disklavier client binary served at /latest-client; empty means there is none")
	)
	flag.Parse()

	ctx := sigctx.New()
	s, err := server.New(ctx, server.Options{
		DBPath:     *dbPath,
		ClientPath: *client,
		Logger:     slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	if err != nil {
		return err
	}
	defer s.Close()

	return s.Go(ctx, *addr)
}
