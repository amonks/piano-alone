package main

import (
	"context"
	"log"
	"net/http"

	"monks.co/piano-alone/sigctx"
)

func main() {
	ctx := sigctx.New()
	handler := NewHandler()
	httpErrs := make(chan error)
	handlerErrs := make(chan error)

	addr := "0.0.0.0:8080"
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Printf("listening on '%s'", addr)

	go func() { handlerErrs <- handler.Start() }()
	go func() { httpErrs <- s.ListenAndServe() }()
	select {
	case <-ctx.Done():
		// interrupt: stop http, stop game
		log.Printf("canceled: %s; shutting down", ctx.Err())
		s.Shutdown(context.Background())
		log.Printf("http server stopped")
		// TODO: stop game
		log.Printf("game server stopped")
	case err := <-httpErrs:
		// http error: stop game
		log.Printf("http server error: %s; shutting down", err)
		// TODO: stop game
		log.Printf("game server stopped")
	case err := <-handlerErrs:
		// game error: stop http
		log.Printf("game error: %s; shutting down", err)
		s.Shutdown(context.Background())
		log.Printf("http server stopped")
	}
}
