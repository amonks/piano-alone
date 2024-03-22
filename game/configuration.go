package game

import (
	"bytes"
	"encoding/gob"
)

type Configuration struct {
	PerformanceID string `gorm:"primaryKey;column:id"`
	Title         string `gorm:"column:title"`
	Composer      string `gorm:"column:composer"`
	Score         []byte `gorm:"column:score"`
}

func ConfigurationFromBytes(bs []byte) *Configuration {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var s Configuration
	if err := dec.Decode(&s); err != nil {
		panic(err)
	}
	return &s
}

func (s *Configuration) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
