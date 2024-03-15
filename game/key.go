package game

import (
	"bytes"
	"encoding/gob"
)

type Key struct {
	Noteno   uint8
	IsNoteOn bool
}

func NewKeyMsg(noteno uint8, isNoteOn bool) *Message {
	return NewMessage(MessageTypeKey, "", (&Key{
		Noteno:   noteno,
		IsNoteOn: isNoteOn,
	}).Bytes())
}

func (m *Key) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func KeyFromBytes(bs []byte) *Key {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var m Key
	if err := dec.Decode(&m); err != nil {
		panic(err)
	}
	return &m
}
