package game

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Every value on the wire is gob, and every decode goes through this
// file. Both websockets and the schedule endpoint take bytes from the
// public internet, so a decoder that panics is a way to kill the
// process from off-site: the decoders return errors and the callers
// answer them. Encoding stays infallible — gob only fails on a type it
// cannot describe, which is a bug in this package rather than
// something a caller can provoke — so encode panics and the Bytes
// methods keep their signatures.

func decode[T any](what string, bs []byte) (T, error) {
	var v T
	if err := gob.NewDecoder(bytes.NewReader(bs)).Decode(&v); err != nil {
		return v, fmt.Errorf("decoding %s: %w", what, err)
	}
	return v, nil
}

func encode(v any) []byte {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
