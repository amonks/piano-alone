package id

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func Random128() string {
	bs := make([]byte, 128)
	io.ReadFull(rand.Reader, bs)
	return hex.EncodeToString(bs)
}
