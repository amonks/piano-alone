//go:build js && wasm

package main

import (
	"syscall/js"
)

func main() {
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
	ctx := canvasEl.Call("getContext", "2d")

	js.Global().Get("console").Call("log", "start")

	i := 0

	var render js.Func
	defer render.Release()
	ctx.Set("fillStyle", "black")
	ctx.Set("strokeStyle", "white")
	render = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.Call("fillRect", 0, 0, width, height)
		x := i % int(width)
		js.Global().Get("console").Call("log", "frame", x)
		ctx.Call("beginPath")
		ctx.Call("moveTo", x, 0)
		ctx.Call("lineTo", x, height)
		ctx.Call("stroke")

		ctx.Call("beginPath")
		ctx.Call("moveTo", 30, 50)
		ctx.Call("lineTo", 150, 100)
		ctx.Call("stroke")

		i++

		js.Global().Call("requestAnimationFrame", render)

		return nil
	})
	js.Global().Call("requestAnimationFrame", render)

	js.Global().Get("console").Call("log", "wait")
	<-done
	js.Global().Get("console").Call("log", "done")
}
