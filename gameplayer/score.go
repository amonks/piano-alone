//go:bulid js && wasm

package main

import (
	"sort"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
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

func BuildTutorialScore(notes []uint8) *Score {
	var (
		out = &Score{NoteTracks: make([]NoteTrack, len(notes))}
	)
	for i, note := range notes {
		var (
			track     = abstrack.New()
			noteStart = time.Duration(i) * time.Second
			noteEnd   = noteStart + time.Second
		)
		track.Append(
			abstrack.AbsEvent{
				Timestamp: noteStart,
				Message:   smf.Message(midi.NoteOn(1, note, 100)),
			},
			abstrack.AbsEvent{
				Timestamp: noteEnd,
				Message:   smf.Message(midi.NoteOff(1, note)),
			},
		)
		out.NoteTracks[i] = NoteTrack{
			Note:  note,
			Track: track,
		}
	}
	for i, note := range notes {
		var (
			noteStart = time.Duration(len(notes)+1) * time.Second
			noteEnd   = noteStart + time.Second
			track     = out.NoteTracks[i].Track
		)
		track.Append(
			abstrack.AbsEvent{
				Timestamp: noteStart,
				Message:   smf.Message(midi.NoteOn(1, note, 100)),
			},
			abstrack.AbsEvent{
				Timestamp: noteEnd,
				Message:   smf.Message(midi.NoteOff(1, note)),
			},
		)
	}
	return out
}
