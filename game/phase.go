package game

import "errors"

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

func PhaseFromBytes(bs []byte) (Phase, error) {
	if len(bs) == 0 {
		return GamePhaseUninitialized, errors.New("decoding phase: empty")
	}
	return Phase(bs[0]), nil
}
