package game

import (
	"bytes"
	"encoding/gob"
)

type JoinMessage struct {
	Fingerprint  string
	NoteCapacity int
}

func (jm JoinMessage) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(jm); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func JoinMessageFromBytes(bs []byte) JoinMessage {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var jm JoinMessage
	if err := dec.Decode(&jm); err != nil {
		panic(err)
	}
	return jm
}
