//go:build js && wasm

package gameclient

import (
	"fmt"
	"sort"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/vdom"
)

func (c *GameClient) Render() vdom.Element {
	return vdom.H("div",
		vdom.HK("section", "ui",
			vdom.H("h1", vdom.T("UI")),
			c.renderUI(),
		),
		vdom.HK("section", "state",
			vdom.H("h1", vdom.T("State")),
			vdom.H("dl",
				vdom.HK("dt", "phase", vdom.T("Phase")),
				vdom.HK("dd", "phase", c.renderPhase()),

				vdom.HK("dt", "players", vdom.T("Players")),
				vdom.HK("dd", "players", c.renderPlayerList()),

				vdom.HK("dt", "notes", vdom.T("Notes")),
				vdom.HK("dd", "notes", c.renderNotes()),
			),
		),
	)
}

func (c *GameClient) renderUI() vdom.Element {
	if c.state == nil {
		panic("nil state")
	}
	switch c.state.Phase.Type {
	case game.GamePhaseHero:
		if c.myRendition != nil {
			return vdom.H("span", vdom.T("already saved rendition"))
		}
		if c.myScore == nil {
			return vdom.H("span", vdom.T("no score"))
		}
		return vdom.C(c.renderCanvas()...)
	default:
		return vdom.H("strong")
	}
}

func (c *GameClient) renderCanvas() []vdom.SceneNode {
	var (
		out           = make([]vdom.SceneNode, len(c.myScore.NoteTracks))
		sectionHeight = 1.0 / float64(len(c.myScore.NoteTracks))

		screenStart = time.Since(c.state.Phase.Begin)
		screenWidth = time.Second * 10
		// screenEnd   = screenStart + screenWidth
	)
	for i, t := range c.myScore.NoteTracks {
		var (
			notes                       = []vdom.SceneNode{}
			ch, key, vel                uint8
			noteStart, noteEnd, noteDur time.Duration
		)
		for _, e := range t.Track.Events {
			// if e.Timestamp < screenStart || e.Timestamp > screenEnd {
			// 	continue
			// }
			if e.Message.GetNoteStart(&ch, &key, &vel) {
				noteStart = e.Timestamp - screenStart
			} else if e.Message.GetNoteEnd(&ch, &key) {
				noteEnd = e.Timestamp - screenStart
				noteDur = noteEnd - noteStart
				notes = append(notes, vdom.Rect(vdom.Bounds{
					X:      float64(noteStart) / float64(screenWidth),
					Y:      0.4,
					Width:  float64(noteDur) / float64(screenWidth),
					Height: 0.2,
				}))
			}
		}
		children := append(notes,
			vdom.Box(vdom.Bounds{0, 0, 1, 1}),
			vdom.Text("48px sans serif", 0, 1, fmt.Sprintf("%d", t.Note)),
		)
		out[i] = vdom.NewContainer(
			vdom.Bounds{X: 0, Y: float64(i) * sectionHeight, Width: 1.0, Height: sectionHeight},
			children...,
		)
	}
	return out
}

func (c *GameClient) renderPhase() *vdom.HTML {
	if c.state.Phase.Exp.IsZero() {
		return vdom.T(c.state.Phase.Type.String())
	}
	return vdom.T(
		"%s (%s)",
		c.state.Phase,
		time.Until(c.state.Phase.Exp).Round(time.Second),
	)
}

func (c *GameClient) renderPlayerList() vdom.Element {
	var playerList []string
	for f := range c.state.Players {
		playerList = append(playerList, f)
	}
	sort.Slice(playerList, func(a, b int) bool { return playerList[a] < playerList[b] })
	var lis []vdom.Element
	for _, f := range playerList {
		player := c.state.Players[f]
		id := player.Fingerprint[:6]
		li := vdom.HK("li", id,
			vdom.H("span", vdom.T(player.Pianist+" ")),
			vdom.H("code", vdom.T("(%s)", id)),
		)
		if player.Fingerprint == c.fingerprint {
			li = li.WithAttr("style", "color: green")
		}
		if player.ConnectionState == game.ConnectionStateDisconnected {
			li = li.WithAttr("style", "opacity: 0.5")
		}
		lis = append(lis, li)
	}
	return vdom.H("ul", lis...)
}

func (c *GameClient) renderNotes() vdom.Element {
	me := c.state.Players[c.fingerprint]
	if me == nil {
		return vdom.T("")
	}
	if len(me.Notes) == 0 {
		return vdom.H("span")
	}
	var lis []vdom.Element
	for _, n := range me.Notes {
		name := midi.Note(n).String()
		lis = append(lis, vdom.HK("li", name, vdom.T(name)))
	}
	return vdom.H("ul", lis...)
}
