//go:build js && wasm

package main

import (
	"syscall/js"
	"time"

	"monks.co/piano-alone/canvas"
)

func (c *GameClient) animate() {
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(_ js.Value, _ []js.Value) any {
		canvas.Draw(c.canvasNode, c.sceneGraph())
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)
}

func (c *GameClient) sceneGraph() []canvas.SceneNode {
	var (
		out          = make([]canvas.SceneNode, len(c.myScore.NoteTracks))
		sectionWidth = 1.0 / float64(len(c.myScore.NoteTracks))

		screenStart = time.Since(c.startPlayingAt) - screenDuration
		screenEnd   = screenStart + screenDuration
	)
	for i, t := range c.myScore.NoteTracks {
		var (
			notes                       = []canvas.SceneNode{}
			ch, key, vel                uint8
			noteStart, noteEnd, noteDur time.Duration
		)
		for _, e := range t.Track.Events {
			if e.Timestamp < screenStart {
				continue
			}
			if e.Message.GetNoteStart(&ch, &key, &vel) {
				if e.Timestamp > screenEnd {
					break
				}
				noteStart = e.Timestamp - screenStart
			} else if e.Message.GetNoteEnd(&ch, &key) {
				noteEnd = e.Timestamp - screenStart
				noteDur = noteEnd - noteStart
				notes = append(notes, canvas.Rect(canvas.Bounds{
					X:      0,
					Y:      1 - float64(noteStart+noteDur)/float64(screenDuration),
					Width:  1,
					Height: float64(noteDur) / float64(screenDuration),
				}))
			}
		}
		out[i] = canvas.G(
			canvas.Bounds{X: float64(i) * sectionWidth, Y: 0, Width: sectionWidth, Height: 1.0},
			notes...,
		)
	}
	return out
}
