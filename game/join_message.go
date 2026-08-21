package game

type JoinMessage struct {
	Fingerprint  string
	NoteCapacity int
}

func (jm JoinMessage) Bytes() []byte { return encode(jm) }

func JoinMessageFromBytes(bs []byte) (JoinMessage, error) {
	return decode[JoinMessage]("join message", bs)
}
