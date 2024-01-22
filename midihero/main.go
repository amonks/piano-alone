//go:build js && wasm

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"

	"monks.co/piano-alone/gameclient"
	"monks.co/piano-alone/jsws"
	"monks.co/piano-alone/proto"
)

func main() {
	go func() {
		ss := &ScrollingScore{}
		ss.Start(context.Background())
	}()

	fingerprint := randomID()
	wc, err := jsws.Open("ws://localhost:8000/ws?fingerprint=" + fingerprint)
	if err != nil {
		panic(err)
	}

	inbox := make(chan *proto.Message)
	outbox := make(chan *proto.Message)

	if _, err := wc.OnMessage(func(bs []byte) {
		inbox <- proto.MessageFromBytes(bs)
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

	gc := gameclient.New(fingerprint)
	gc.Start(outbox, inbox)
}

func randomID() string {
	bs := make([]byte, 128)
	io.ReadFull(rand.Reader, bs)
	return hex.EncodeToString(bs)
}
