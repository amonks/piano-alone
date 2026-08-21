package game

type Player struct {
	ConnectionState      ConnectionState
	Fingerprint          string
	NoteCapacity         int
	AssignedNotes        []uint8
	HasStartedTutorial   bool
	HasCompletedTutorial bool
	HasSubmitted         bool
}

func PlayerFromBytes(bs []byte) (*Player, error) {
	p, err := decode[Player]("player", bs)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Player) Bytes() []byte { return encode(p) }
