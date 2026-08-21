package game

import "time"

type Performance struct {
	Configuration *Configuration
	Date          time.Time
	IsFeatured    bool

	IsComplete  bool
	Rendition   []byte
	PlayerCount int
}

func PerformanceFromBytes(bs []byte) (*Performance, error) {
	p, err := decode[Performance]("performance", bs)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Performance) Bytes() []byte { return encode(s) }

func PerformancesFromBytes(bs []byte) ([]*Performance, error) {
	return decode[[]*Performance]("performances", bs)
}

func PerformancesToBytes(ps []*Performance) []byte { return encode(ps) }
