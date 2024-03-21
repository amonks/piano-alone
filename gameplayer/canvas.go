//go:build js && wasm

package main

import (
	"syscall/js"
	"time"

	"monks.co/piano-alone/canvas"
)

func (c *GameClient) animate(sceneGraphFunc func() ([]canvas.SceneNode, bool)) <-chan struct{} {
	done := make(chan struct{})
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(_ js.Value, _ []js.Value) any {
		sceneGraph, hasNext := sceneGraphFunc()
		if !hasNext {
			done <- struct{}{}
			return nil
		}
		canvas.Draw(canv, sceneGraph)
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)
	return done
}

func (c *GameClient) tutorialSceneGraph() ([]canvas.SceneNode, bool) {
	return renderScore(c.tutorialScore, c.startTutorialAt)
}

func (c *GameClient) performanceSceneGraph() ([]canvas.SceneNode, bool) {
	return renderScore(c.myScore, c.startPlayingAt)
}

func renderScore(score *Score, startAt time.Time) ([]canvas.SceneNode, bool) {
	var (
		noteTracks   = score.NoteTracks
		columns      = len(noteTracks)
		out          = make([]canvas.SceneNode, columns)
		sectionWidth = 1.0 / float64(columns)
		screenStart  = time.Since(startAt) - screenDuration
		screenEnd    = screenStart + screenDuration
	)
	sawAnyNotes := false
	for i, t := range noteTracks {
		var (
			notes                       = []canvas.SceneNode{}
			ch, key, vel                uint8
			noteStart, noteEnd, noteDur time.Duration
		)
		for _, e := range t.Track.Events {
			if e.Timestamp < screenStart {
				continue
			} else {
				sawAnyNotes = true
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
	if !sawAnyNotes {
		return nil, false
	}
	return out, true
}
