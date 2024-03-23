package game

import (
	"bytes"
	"encoding/gob"
)

type Player struct {
	ConnectionState      ConnectionState
	Fingerprint          string
	NoteCapacity         int
	AssignedNotes        []uint8
	HasStartedTutorial   bool
	HasCompletedTutorial bool
	HasSubmitted         bool
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
