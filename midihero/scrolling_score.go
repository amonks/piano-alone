//go:build wasm || js

package main

import (
	"context"
	"syscall/js"
)

type ScrollingScore struct{}

func (ss *ScrollingScore) Start(ctx context.Context) {
	var (
		done     = make(chan struct{})
		doc      = js.Global().Get("document")
		body     = doc.Get("body")
		canvasEl = doc.Call("querySelector", "canvas")
		width    = body.Get("clientWidth").Float()
		height   = body.Get("clientHeight").Float()
	)
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	c2d := canvasEl.Call("getContext", "2d")

	i := 0

	var render js.Func
	defer render.Release()
	c2d.Set("fillStyle", "black")
	c2d.Set("strokeStyle", "white")
	render = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if ctx.Err() != nil {
			done <- struct{}{}
			return nil
		}
		c2d.Call("fillRect", 0, 0, width, height)
		x := i % int(width)
		c2d.Call("beginPath")
		c2d.Call("moveTo", x, 0)
		c2d.Call("lineTo", x, height)
		c2d.Call("stroke")

		c2d.Call("beginPath")
		c2d.Call("moveTo", 30, 50)
		c2d.Call("lineTo", 150, 100)
		c2d.Call("stroke")

		i++
		js.Global().Call("requestAnimationFrame", render)
		return nil
	})
	js.Global().Call("requestAnimationFrame", render)
	<-done
}
