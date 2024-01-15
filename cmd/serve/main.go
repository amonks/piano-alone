package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	addr := "0.0.0.0:8000"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "err opening port", err)
		os.Exit(1)
	}
	fmt.Printf("Listening at %s\n", addr)
	log.Fatal(http.Serve(listener, http.FileServer(http.Dir("./website"))))
}
