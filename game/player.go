package game

import (
	"bytes"
	"encoding/gob"
)

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
