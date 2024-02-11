package game

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Phase struct {
	Type  GamePhase
	Begin time.Time
	Exp   time.Time
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=GamePhase
type GamePhase byte

const (
	GamePhaseUninitialized GamePhase = iota
	GamePhaseLobby
	GamePhaseHero
	GamePhaseProcessing
	GamePhasePlayback
	GamePhaseDone
)

func NewPhase(t GamePhase) Phase {
	return Phase{
		Type:  t,
		Begin: time.Now(),
	}
}

func (p Phase) WithExp(exp time.Time) Phase {
	return Phase{
		Type:  p.Type,
		Begin: p.Begin,
		Exp:   exp,
	}
}

func (m Phase) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func PhaseFromBytes(bs []byte) Phase {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var m Phase
	if err := dec.Decode(&m); err != nil {
		panic(err)
	}
	return m
}
