package game

import (
	"crypto/sha256"
	"encoding/hex"
)

type Player struct {
	ConnectionState ConnectionState
	// Fingerprint is the player's id (PlayerID), never the fingerprint
	// its socket presents. The field keeps its name because it is gob
	// on the wire, and the disklavier client, built and notarized by
	// hand, decodes it by that name.
	Fingerprint          string
	NoteCapacity         int
	AssignedNotes        []uint8
	HasStartedTutorial   bool
	HasCompletedTutorial bool
	HasSubmitted         bool
}

// PlayerID is the name a player goes by everywhere past its socket:
// in the state, in what the controllers hear, in the logs. A
// fingerprint is the player's whole credential — a socket that
// presents one acts as that player — so the server takes it at the
// socket and never repeats it, and every client hears of every player
// only by this one-way digest of it. A player's page knows its own
// fingerprint, and finds itself by its digest too.
func PlayerID(fingerprint string) string {
	sum := sha256.Sum256([]byte(fingerprint))
	return hex.EncodeToString(sum[:])
}

func PlayerFromBytes(bs []byte) (*Player, error) {
	p, err := decode[Player]("player", bs)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Player) Bytes() []byte { return encode(p) }
