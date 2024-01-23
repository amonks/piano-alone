//go:build wasm || js

// jsws is based on this code:
//
// https://github.com/nhooyr/websocket/blob/master/internal/wsjs/wsjs_js.go
//
// which comes with the following notice:
//
//     Copyright (c) 2023 Anmol Sethi <hi@nhooyr.io>
//
//     Permission to use, copy, modify, and distribute this software for any
//     purpose with or without fee is hereby granted, provided that the above
//     copyright notice and this permission notice appear in all copies.
//
//     THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
//     WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
//     MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
//     ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
//     WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
//     ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
//     OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package jsws

import (
	"fmt"
	"syscall/js"
)

type WebsocketClient struct {
	ws js.Value
}

func Open(url string) (wc *WebsocketClient, err error) {
	defer handleJSException(&err)

	ws := js.Global().Get("WebSocket").New(url)
	wc = &WebsocketClient{ws: ws}

	ready := make(chan struct{})
	wc.addEventListener("open", func(_ js.Value) {
		ready <- struct{}{}
	})
	<-ready

	ws.Set("binaryType", "arraybuffer")

	return wc, nil
}

func (wc *WebsocketClient) Close(code int, reason string) (err error) {
	defer handleJSException(&err)
	wc.ws.Call("close", code, reason)
	return nil
}

func (wc *WebsocketClient) Send(msg []byte) (err error) {
	defer handleJSException(&err)
	wc.ws.Call("send", uint8Array(msg))
	return nil
}

func (wc *WebsocketClient) OnMessage(handler func([]byte)) (func(), error) {
	return wc.addEventListener("message", func(val js.Value) {
		data := val.Get("data")
		if data.Type() == js.TypeString {
			handler([]byte(data.String()))
		} else {
			handler(extractArrayBuffer(data))
		}
	})
}

func (wc *WebsocketClient) OnError(handler func(error)) (func(), error) {
	return wc.addEventListener("error", func(val js.Value) {
		message := val.Get("message").String()
		handler(fmt.Errorf("websocket error: '%s'", message))
	})
}

func (wc *WebsocketClient) addEventListener(event string, handler func(js.Value)) (stop func(), err error) {
	defer handleJSException(&err)
	f := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		handler(args[0])
		return nil
	})
	wc.ws.Call("addEventListener", event, f)
	return func() {
		wc.ws.Call("removeEventListener", event, f)
		f.Release()
	}, nil
}

func handleJSException(err *error) {
	r := recover()
	if jsErr, ok := r.(js.Error); ok {
		*err = jsErr
		return
	}
	if r != nil {
		panic(r)
	}
}
func extractArrayBuffer(arrayBuffer js.Value) []byte {
	uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)
	dst := make([]byte, uint8Array.Length())
	js.CopyBytesToGo(dst, uint8Array)
	return dst
}

func uint8Array(src []byte) js.Value {
	uint8Array := js.Global().Get("Uint8Array").New(len(src))
	js.CopyBytesToJS(uint8Array, src)
	return uint8Array
}
