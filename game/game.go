package game

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gitlab.com/gomidi/midi/v2/smf"
)

type State struct {
	Score     *smf.SMF
	Phase     Phase
	Players   map[string]*Player
	Rendition *smf.SMF

	DisklavierIsConnected bool
	ConductorIsConnected  bool
}

func init() {
	var mt smf.MetricTicks
	gob.Register(mt)
}

func NewState() *State {
	return &State{
		Phase:   NewPhase(GamePhaseLobby),
		Players: map[string]*Player{},
		Score:   nil,
	}
}

func (s *State) CountConnectedPlayers() int {
	n := 0
	for _, p := range s.Players {
		if p.ConnectionState == ConnectionStateConnected {
			n++
		}
	}
	return n
}

func (s *State) CountSubmittedTracks() int {
	n := 0
	for _, p := range s.Players {
		if p.HasSubmitted {
			n++
		}
	}
	return n
}

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
	HasSubmitted    bool
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

func NewMessage(messageType MessageType, player string, data []byte) *Message {
	return &Message{messageType, player, data}
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=MessageType
type MessageType byte

const (
	MessageTypeInvalid MessageType = iota

	MessageTypeExpireLobby
	MessageTypeExpireHero

	MessageTypeDisklavierConnected
	MessageTypeDisklavierDisconnected
	MessageTypeConductorConnected
	MessageTypeConductorDisconnected

	MessageTypeJoin
	MessageTypeLeave
	MessageTypeSubmitPartialTrack

	MessageTypeInitialState
	MessageTypeBroadcastConnectedPlayer
	MessageTypeBroadcastDisconnectedPlayer
	MessageTypeAssignment
	MessageTypeBroadcastPhase
	MessageTypeBroadcastSubmittedTrack
	MessageTypeBroadcastCombinedTrack
	MessageTypeBroadcastControllerModal

	MessageTypeRestart
	MessageTypeAdvancePhase
)

func (m *Message) String() string {
	switch m.Type {
	case MessageTypeBroadcastPhase:
		phase := PhaseFromBytes(m.Data)
		return fmt.Sprintf("%s: %s [%s]", m.Player, m.Type.String(), phase.Type.String())
	default:
		return m.Type.String()
	}
}

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
