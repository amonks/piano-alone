package gameclient

import (
	"bytes"
	"sort"
	"syscall/js"
	"time"

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
			),
		),
	)
}

func (c *GameClient) renderUI() vdom.Element {
	switch c.state.Phase {
	case game.GamePhaseHero:
		if c.myRendition != nil {
			return vdom.H("span", vdom.T("already saved rendition"))
		}
		if c.myScore == nil {
			return vdom.H("span", vdom.T("no score"))
		}
		return vdom.H("button", vdom.T("submit")).
			WithAttr("onclick", js.FuncOf(func(js.Value, []js.Value) any {
				c.myRendition = c.myScore

				var bs bytes.Buffer
				c.myRendition.WriteTo(&bs)
				c.send <- &game.Message{
					Type: game.MessageTypeSubmitPartialTrack,
					Data: bs.Bytes(),
				}
				return nil
			}))
	default:
		return vdom.H("strong")
	}
}

func (c *GameClient) renderPhase() *vdom.HTML {
	if c.state.PhaseExp.IsZero() {
		return vdom.T(c.state.Phase.String())
	}
	return vdom.T(
		"%s (%s)",
		c.state.Phase,
		time.Until(c.state.PhaseExp).Round(time.Second),
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
