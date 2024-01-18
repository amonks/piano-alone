package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"monks.co/piano-alone/server"
)

func main() {
	addr := "0.0.0.0:8000"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "err opening port", err)
		os.Exit(1)
	}
	fmt.Printf("Listening at %s\n", addr)
	s := server.New()
	log.Fatal(http.Serve(listener, s))
}
