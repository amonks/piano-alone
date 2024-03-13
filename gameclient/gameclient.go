//go:build js && wasm

package gameclient

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/vdom"
)

const (
	fps         = 30
	screenWidth = time.Second * 4
)

type GameClient struct {
	send        chan<- *game.Message
	fingerprint string
	state       *game.State
	myScore     *Score
	myRendition *smf.SMF
	vdom        *vdom.VDOM
}

func New(fingerprint string, root js.Value) *GameClient {
	return &GameClient{
		state:       game.NewState(),
		fingerprint: fingerprint,
		vdom:        vdom.New(root),
	}
}

func (c *GameClient) Start(send chan<- *game.Message, recv <-chan *game.Message) error {
	c.send = send
	send <- game.NewMessage(
		game.MessageTypeJoin,
		c.fingerprint,
		[]byte(c.fingerprint),
	)
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		for {
			select {
			case <-ctx.Done():
				cancel()
				return
			case <-time.After(time.Second / fps):
				c.vdom.Render(c.Render())
			}
		}
	}()
	for m := range recv {
		log.Printf("incoming message: '%s'", m.Type)
		if err := c.handleMessage(m); err != nil {
			return fmt.Errorf("error handling message '%s': %w", m.Type, err)
		}
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
		msg := game.PhaseFromBytes(m.Data)
		c.state.Phase = msg
	case game.MessageTypeBroadcastConnectedPlayer:
		player := game.PlayerFromBytes(m.Data)
		c.state.Players[player.Fingerprint] = player
	case game.MessageTypeBroadcastDisconnectedPlayer:
		c.state.Players[string(m.Data)].ConnectionState = game.ConnectionStateDisconnected
	case game.MessageTypeAssignment:
		me := c.state.Players[c.fingerprint]
		me.Notes = m.Data
		c.myScore = NewScore(abstrack.FromSMF(c.state.Score, 0).Select(me.Notes))
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
