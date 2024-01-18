//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"sync"

	"monks.co/piano-alone/jsws"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ss := &ScrollingScore{}
		ss.Start(context.Background())
		wg.Done()
	}()

	wc, err := jsws.Open("ws://localhost:8000/ws")
	if err != nil {
		panic(err)
	}

	if _, err := wc.OnMessage(func(msg []byte) {
		fmt.Println("got", string(msg))
	}); err != nil {
		panic(err)
	}

	if _, err := wc.OnError(func(err error) {
		panic(err)
	}); err != nil {
		panic(err)
	}

	fmt.Println("send")
	if err := wc.Send([]byte("hello")); err != nil {
		panic(err)
	}

	wg.Wait()
	fmt.Println("done waiting")
}
