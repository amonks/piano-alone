package game

type Configuration struct {
	PerformanceID string
	Title         string
	Composer      string
	Score         []byte
}

func ConfigurationFromBytes(bs []byte) (*Configuration, error) {
	c, err := decode[Configuration]("configuration", bs)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Configuration) Bytes() []byte { return encode(s) }
