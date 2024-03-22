package game

//go:generate go run golang.org/x/tools/cmd/stringer -type=Phase
type Phase byte

const (
	GamePhaseUninitialized Phase = iota
	GamePhaseLobby
	GamePhaseHero
	GamePhaseProcessing
	GamePhasePlayback
	GamePhaseDone
)

func (m Phase) Bytes() []byte {
	return []byte{byte(m)}
}

func PhaseFromBytes(bs []byte) Phase {
	return Phase(bs[0])
}
