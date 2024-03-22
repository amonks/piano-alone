package game

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

	MessageTypeConfigurePerformance

	MessageTypeDisklavierConnected
	MessageTypeDisklavierDisconnected
	MessageTypeConductorConnected
	MessageTypeConductorDisconnected

	MessageTypeJoin
	MessageTypeLeave
	MessageTypeSubmitPartialTrack

	MessageTypeState
	MessageTypeBeginPerformance
	MessageTypeBroadcastConnectedPlayer
	MessageTypeBroadcastDisconnectedPlayer
	MessageTypeAssignment
	MessageTypeBroadcastPhase
	MessageTypeBroadcastSubmittedTrack

	MessageTypeBroadcastControllerModal
	MessageTypeSendRenditionToDisklavier

	MessageTypeRestart
	MessageTypeAdvancePhase
)

func (m *Message) String() string {
	switch m.Type {
	case MessageTypeBroadcastPhase:
		phase := PhaseFromBytes(m.Data)
		return fmt.Sprintf("%s: %s [%s]", m.Player, m.Type, phase)
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
