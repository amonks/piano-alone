package proto

import (
	"bytes"
	"encoding/gob"
)

//go:generate stringer -type=GamePhase
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

type Message struct {
	Type   string
	Player string
	Data   []byte
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
