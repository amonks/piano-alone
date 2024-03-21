//go:build js && wasm

package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"syscall/js"

	"monks.co/piano-alone/baseurl"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/jsws"
	"monks.co/piano-alone/storage"
)

func main() {
	var (
		doc  = js.Global().Get("document")
		root = doc.Call("getElementById", "root")

		baseURL = baseurl.Discover()

		inbox  = make(chan *game.Message)
		outbox = make(chan *game.Message)
	)

	fingerprint := storage.Session.Get("fingerprint")
	if fingerprint == "" {
		fingerprint = randomID()
		storage.Session.Set("fingerprint", fingerprint)
	}
	log.Printf("fingerprint: %s", fingerprint)

	wc, err := jsws.Open(baseURL.WS("/ws?fingerprint=" + fingerprint))
	if err != nil {
		panic(err)
	}

	if _, err := wc.OnMessage(func(bs []byte) {
		inbox <- game.MessageFromBytes(bs)
	}); err != nil {
		panic(err)
	}

	if _, err := wc.OnError(func(err error) {
		log.Printf("ws error: %s", err)
	}); err != nil {
		log.Printf("error creating ws error handler: %s", err)
	}

	go func() {
		for m := range outbox {
			bs := m.Bytes()
			if err := wc.Send(bs); err != nil {
				panic(err)
			}
		}
	}()

	gc := New(fingerprint, root)
	if err := gc.Start(outbox, inbox); err != nil {
		panic(err)
	}
}

func randomID() string {
	bs := make([]byte, 128)
	io.ReadFull(rand.Reader, bs)
	return hex.EncodeToString(bs)
}
