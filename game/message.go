package game

import (
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
	MessageTypeStartTutorial
	MessageTypeCompleteTutorial

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

// Control reports whether this message type drives the performance:
// the three the conductor sends, which are the whole of what the game
// loop will act on from anyone. Everything else is a player reporting
// on itself or the server talking to a client.
func (t MessageType) Control() bool {
	switch t {
	case MessageTypeRestart, MessageTypeAdvancePhase, MessageTypeBeginPerformance:
		return true
	default:
		return false
	}
}

// ServerOnly reports whether this message type is the server's own
// account of something rather than a client's report about itself.
// The four presence types are synthesized where a controller socket
// opens and closes; a client that sends one is claiming a connection
// exists that the server has not seen.
func (t MessageType) ServerOnly() bool {
	switch t {
	case MessageTypeDisklavierConnected, MessageTypeDisklavierDisconnected,
		MessageTypeConductorConnected, MessageTypeConductorDisconnected:
		return true
	default:
		return false
	}
}

func (m *Message) String() string {
	if m.Type == MessageTypeBroadcastPhase {
		// A phase that will not decode is not worth failing a log line
		// over; the type alone still says what happened.
		if phase, err := PhaseFromBytes(m.Data); err == nil {
			return fmt.Sprintf("%s: %s [%s]", m.Player, m.Type, phase)
		}
	}
	return m.Type.String()
}

func (m *Message) Bytes() []byte { return encode(m) }

func MessageFromBytes(bs []byte) (*Message, error) {
	m, err := decode[Message]("message", bs)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
