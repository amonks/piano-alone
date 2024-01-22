package gameclient

import (
	"bytes"
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/proto"
)

type GameClient struct {
	fingerprint string
}

func New(fingerprint string) *GameClient {
	return &GameClient{fingerprint: fingerprint}
}

func (c *GameClient) Start(send chan<- *proto.Message, recv <-chan *proto.Message) error {
	send <- &proto.Message{
		Type:   "join",
		Player: c.fingerprint,
		Data:   []byte(c.fingerprint),
	}

	// splitMsg phase
	splitMsg := selectMessage("split", recv, msgHasType("split"))
	r := bytes.NewReader(splitMsg.Data)
	splitTrack, err := smf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("error constructing received smf: %w", err)
	}

	// hero phase
	selectMessage("hero", recv, msgIsPhase(proto.GamePhaseHero))
	playedTrack := splitTrack
	var buf bytes.Buffer
	playedTrack.WriteTo(&buf)
	send <- &proto.Message{
		Type:   "track",
		Player: c.fingerprint,
		Data:   buf.Bytes(),
	}

	// combine tracks
	selectMessage("joining", recv, msgIsPhase(proto.GamePhaseJoining))

	// playback phase
	selectMessage("playback", recv, msgIsPhase(proto.GamePhasePlayback))
	combined := selectMessage("combined", recv, msgHasType("combined"))
	log.Println("combined:", combined.Data)

	// done
	selectMessage("done", recv, msgIsPhase(proto.GamePhaseDone))

	return nil
}

func msgHasType(t string) func(*proto.Message) bool {
	return func(msg *proto.Message) bool {
		return msg.Type == t
	}
}

func msgIsPhase(p proto.GamePhase) func(*proto.Message) bool {
	return func(msg *proto.Message) bool {
		return msg.Type == "phase" && proto.GamePhase(msg.Data[0]) == p
	}
}

func selectMessage(label string, c <-chan *proto.Message, pred func(m *proto.Message) bool) *proto.Message {
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
