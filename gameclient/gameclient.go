package gameclient

import (
	"bytes"
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/game"
)

type GameClient struct {
	fingerprint string
}

func New(fingerprint string) *GameClient {
	return &GameClient{fingerprint: fingerprint}
}

func (c *GameClient) Start(send chan<- *game.Message, recv <-chan *game.Message) error {
	send <- &game.Message{
		Type:   game.MessageTypeJoin,
		Player: c.fingerprint,
		Data:   []byte(c.fingerprint),
	}

	// splitMsg phase
	splitMsg := game.SelectMessageByType(recv, game.MessageTypeAssignment)
	r := bytes.NewReader(splitMsg.Data)
	splitTrack, err := smf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("error constructing received smf: %w", err)
	}

	// hero phase
	game.SelectPhaseChangeMessage(recv, game.GamePhaseHero)
	playedTrack := splitTrack
	var buf bytes.Buffer
	playedTrack.WriteTo(&buf)
	send <- &game.Message{
		Type:   game.MessageTypeSubmitPartialTrack,
		Player: c.fingerprint,
		Data:   buf.Bytes(),
	}

	// combine tracks
	selectMessage("joining", recv, func(m *game.Message) bool {
		return m.IsPhaseChangeTo(game.GamePhaseJoining)
	})

	// playback phase
	selectMessage("playback", recv, func(m *game.Message) bool {
		return m.IsPhaseChangeTo(game.GamePhasePlayback)
	})
	combined := selectMessage("combined", recv, func(m *game.Message) bool {
		return m.HasType(game.MessageTypeBroadcastCombinedTrack)
	})
	log.Println("combined:", combined.Data)

	// done
	selectMessage("done", recv, func(m *game.Message) bool {
		return m.IsPhaseChangeTo(game.GamePhaseDone)
	})

	return nil
}

func selectMessage(label string, c <-chan *game.Message, pred func(m *game.Message) bool) *game.Message {
	log.Printf("waiting for %s", label)
	for m := range c {
		if pred(m) {
			log.Printf("got %s", label)
			return m
		} else {
			log.Printf("ignore message (waiting for %s): %+v", label, m)
		}
	}
	return nil
}
