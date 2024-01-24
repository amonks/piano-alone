package game

import (
	"bytes"
	"encoding/gob"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
)

type State struct {
	Score     *smf.SMF
	Phase     GamePhase
	PhaseExp  time.Time
	Players   map[string]*Player
	Rendition *smf.SMF
}

func init() {
	var mt smf.MetricTicks
	gob.Register(mt)
}

func NewState() *State {
	return &State{
		Phase:    GamePhaseUninitialized,
		PhaseExp: time.Time{},
		Players:  map[string]*Player{},
		Score:    nil,
	}
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=GamePhase
type GamePhase byte

const (
	GamePhaseUninitialized GamePhase = iota
	GamePhaseLobby
	GamePhaseSplitting
	GamePhaseHero
	GamePhaseJoining
	GamePhasePlayback
	GamePhaseDone
)

func StateFromBytes(bs []byte) *State {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var s State
	if err := dec.Decode(&s); err != nil {
		panic(err)
	}
	return &s
}

func (s *State) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type Player struct {
	ConnectionState ConnectionState
	Fingerprint     string
	Pianist         string
	Notes           []uint8
}

func PlayerFromBytes(bs []byte) *Player {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var p Player
	if err := dec.Decode(&p); err != nil {
		panic(err)
	}
	return &p
}

func (p *Player) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=ConnectionState
type ConnectionState byte

const (
	ConnectionStateUninitialized ConnectionState = iota
	ConnectionStateDisconnected
	ConnectionStateConnected
)

type Message struct {
	Type   MessageType
	Player string
	Data   []byte
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=MessageType
type MessageType byte

const (
	MessageTypeInvalid MessageType = iota

	MessageTypeExpireLobby
	MessageTypeExpireHero
	MessageTypeExpirePlayback

	MessageTypeJoin
	MessageTypeLeave
	MessageTypeSubmitPartialTrack

	MessageTypeInitialState
	MessageTypeBroadcastConnectedPlayer
	MessageTypeBroadcastDisconnectedPlayer
	MessageTypeAssignment
	MessageTypeBroadcastPhase
	MessageTypeBroadcastCombinedTrack
)

func (m *Message) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func MessageFromBytes(bs []byte) *Message {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var m Message
	if err := dec.Decode(&m); err != nil {
		panic(err)
	}
	return &m
}

func SelectMessage(c <-chan *Message, pred func(*Message) bool) *Message {
	for m := range c {
		if pred(m) {
			return m
		}
	}
	return nil
}

func SelectMessageByType(c <-chan *Message, t MessageType) *Message {
	return SelectMessage(c, func(m *Message) bool { return m.HasType(t) })
}

func SelectPhaseChangeMessage(c <-chan *Message, p GamePhase) *Message {
	return SelectMessage(c, func(m *Message) bool { return m.IsPhaseChangeTo(p) })
}

func (m *Message) HasType(t MessageType) bool {
	return m.Type == t
}

func (m *Message) IsPhaseChangeTo(p GamePhase) bool {
	return m.Type == MessageTypeBroadcastPhase && GamePhase(m.Data[0]) == p
}

type PhaseChangeMessage struct {
	Phase GamePhase
	Exp   time.Time
}

func (m *PhaseChangeMessage) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func PhaseChangeMessageFromBytes(bs []byte) *PhaseChangeMessage {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var m PhaseChangeMessage
	if err := dec.Decode(&m); err != nil {
		panic(err)
	}
	return &m
}
