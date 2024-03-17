//go:build js && wasm

package gameclient

import (
	"bytes"
	"fmt"
	"log"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/recorder"
)

const (
	fps            = 30
	screenDuration = time.Second * 4
)

type GameClient struct {
	loopback       chan<- message
	send           chan<- *game.Message
	fingerprint    string
	state          *game.State
	myScore        *Score
	myRendition    *smf.SMF
	pianoHandle    js.Value
	canvasNode     js.Value
	startPlayingAt time.Time
	recorder       *recorder.Recorder
}

func New(fingerprint string, root js.Value) *GameClient {
	return &GameClient{
		state:       game.NewState(),
		fingerprint: fingerprint,
		recorder:    recorder.New(120),
	}
}

func (c *GameClient) Start(send chan<- *game.Message, recv <-chan *game.Message) error {
	c.send = send

	loopback := make(chan message)
	c.loopback = loopback

	global := js.Global()
	doc := global.Get("document")
	pianoNode := doc.Call("getElementById", "piano")
	c.pianoHandle = global.Call("Piano", pianoNode)
	c.canvasNode = doc.Call("querySelector", "canvas")

	send <- game.NewMessage(
		game.MessageTypeJoin,
		c.fingerprint,
		[]byte(c.fingerprint),
	)
	for {
		var m message
		var ok bool
		select {
		case m, ok = <-recv:
		case m, ok = <-loopback:
		}
		if !ok {
			break
		}
		log.Printf("incoming message: '%s'", m)
		if err := c.handleMessage(m); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
	log.Printf("done")
	return nil
}

type message interface{}

type (
	msgHeroReady struct{}
	msgHeroDone  struct{}
	msgKey       = recorder.Event
)

func (c *GameClient) handleMessage(m message) error {
	switch m := m.(type) {

	case msgHeroReady:
		c.startPlayingAt = time.Now()
		c.animate()
		go func() {
			time.Sleep(c.myScore.NoteTracks[0].Track.Dur() + screenDuration)
			c.loopback <- msgHeroDone{}
		}()

	case msgHeroDone:
		c.recorder.Close()
		bs, err := c.recorder.Bytes()
		if err != nil {
			panic(err)
		}
		c.send <- game.NewMessage(game.MessageTypeSubmitPartialTrack, "", bs)

	case msgKey:
		c.recorder.Record(m)

	case *game.Message:
		switch m.Type {

		case game.MessageTypeInitialState:
			c.state = game.StateFromBytes(m.Data)

		case game.MessageTypeBroadcastPhase:
			msg := game.PhaseFromBytes(m.Data)
			c.state.Phase = msg
			switch msg.Type {
			case game.GamePhaseHero:
				notes := c.state.Players[c.fingerprint].Notes
				list := fmt.Sprintf("%d", notes[0])
				for _, n := range notes[1:] {
					list += fmt.Sprintf(",%d", n)
				}
				c.pianoHandle.Call("transition", list,
					js.FuncOf(func(_ js.Value, args []js.Value) any {
						eventType := args[0].String()
						noteno := uint8(args[1].Int())

						switch eventType {

						case "on":
							c.loopback <- recorder.Now(midi.NoteOn(1, noteno, 100))

						case "off":
							c.loopback <- recorder.Now(midi.NoteOff(1, noteno))

						case "ready":
							c.loopback <- msgHeroReady{}
						}
						return nil
					}),
				)
			}

		case game.MessageTypeBroadcastConnectedPlayer:
			player := game.PlayerFromBytes(m.Data)
			c.state.Players[player.Fingerprint] = player
		case game.MessageTypeBroadcastDisconnectedPlayer:
			c.state.Players[string(m.Data)].ConnectionState = game.ConnectionStateDisconnected
		case game.MessageTypeAssignment:
			me := c.state.Players[c.fingerprint]
			me.Notes = m.Data
			c.myScore = NewScore(abstrack.FromSMF(c.state.Score, 0).Select(me.Notes))
		case game.MessageTypeBroadcastSubmittedTrack:
			fingerprint := string(m.Data)
			if _, ok := c.state.Players[fingerprint]; !ok {
				c.state.Players[fingerprint] = &game.Player{}
			}
			c.state.Players[fingerprint].HasSubmitted = true
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
	}
	return nil
}
