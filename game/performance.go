package game

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Performance struct {
	Configuration *Configuration `gorm:"embedded"`
	Date          time.Time      `gorm:"column:date"`
	IsFeatured    bool           `gorm:"column:is_featured"`

	IsComplete  bool   `gorm:"column:is_complete"`
	Rendition   []byte `gorm:"column:rendition"`
	PlayerCount int    `gorm:"column:player_count"`
}

func PerformanceFromBytes(bs []byte) *Performance {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var s Performance
	if err := dec.Decode(&s); err != nil {
		panic(err)
	}
	return &s
}

func (s *Performance) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func PerformancesFromBytes(bs []byte) []*Performance {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var s []*Performance
	if err := dec.Decode(&s); err != nil {
		panic(err)
	}
	return s
}

func PerformancesToBytes(ps []*Performance) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(ps); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
