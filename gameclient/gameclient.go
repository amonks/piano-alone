package gameclient

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/vdom"
)

type GameClient struct {
	send        chan<- *game.Message
	fingerprint string
	state       *game.State
	myScore     *smf.SMF
	myRendition *smf.SMF
	vdom        *vdom.VDOM
}

func New(fingerprint string, root js.Value) *GameClient {
	return &GameClient{
		fingerprint: fingerprint,
		vdom:        vdom.New(root),
	}
}

func (c *GameClient) Start(send chan<- *game.Message, recv <-chan *game.Message) error {
	c.send = send
	send <- &game.Message{
		Type:   game.MessageTypeJoin,
		Player: c.fingerprint,
		Data:   []byte(c.fingerprint),
	}
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second / 10):
				c.vdom.Render(c.Render())
			}
		}
	}()
	for m := range recv {
		log.Printf("incoming message: '%s'", m.Type)
		if err := c.handleMessage(m); err != nil {
			return fmt.Errorf("error handling message '%s': %w", m.Type, err)
		}
		c.vdom.Render(c.Render())
	}
	cancel()
	log.Printf("done")
	return nil
}

func (c *GameClient) handleMessage(m *game.Message) error {
	switch m.Type {
	case game.MessageTypeInitialState:
		c.state = game.StateFromBytes(m.Data)
	case game.MessageTypeBroadcastPhase:
		msg := game.PhaseChangeMessageFromBytes(m.Data)
		c.state.Phase = msg.Phase
		c.state.PhaseExp = msg.Exp
	case game.MessageTypeBroadcastConnectedPlayer:
		player := game.PlayerFromBytes(m.Data)
		c.state.Players[player.Fingerprint] = player
	case game.MessageTypeBroadcastDisconnectedPlayer:
		c.state.Players[string(m.Data)].ConnectionState = game.ConnectionStateDisconnected
	case game.MessageTypeAssignment:
		r := bytes.NewReader(m.Data)
		score, err := smf.ReadFrom(r)
		if err != nil {
			return fmt.Errorf("error reading assignment smf: %w", err)
		}
		c.myScore = score
	case game.MessageTypeBroadcastCombinedTrack:
		r := bytes.NewReader(m.Data)
		rendition, err := smf.ReadFrom(r)
		if err != nil {
			return fmt.Errorf("error reading combined smf: %w", err)
		}
		c.state.Rendition = rendition
	default:
		log.Printf("not handling message (type: '%s')", m.Type.String())
	}
	return nil
}

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
