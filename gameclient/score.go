//go:bulid js && wasm

package gameclient

import (
	"sort"

	"monks.co/piano-alone/abstrack"
)

type Score struct {
	NoteTracks []NoteTrack
}

type NoteTrack struct {
	Note  uint8
	Track *abstrack.AbsTrack
}

func NewScore(track *abstrack.AbsTrack) *Score {
	if track.Events[0].Timestamp < 0 {
		panic("invalid track")
	}
	score := &Score{}
	notes := track.CountNotes()
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Key < notes[j].Key
	})
	score.NoteTracks = make([]NoteTrack, len(notes))
	for i, n := range notes {
		score.NoteTracks[i] = NoteTrack{
			Note:  n.Key,
			Track: track.Select([]uint8{n.Key}),
		}
	}
	return score
}
