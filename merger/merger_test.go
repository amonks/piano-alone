package merger_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/merger"
)

func TestMerger(t *testing.T) {
	var trackA, trackB smf.Track

	trackA.Add(0, smf.MetaTempo(120))
	trackA.Add(0, midi.NoteOn(0, 100, 100))
	trackA.Add(10, midi.NoteOn(0, 102, 100))
	trackA.Add(10, midi.NoteOn(0, 104, 100))

	trackB.Add(0, smf.MetaTempo(120))
	trackB.Add(5, midi.NoteOn(0, 101, 100))
	trackB.Add(10, midi.NoteOn(0, 103, 100))
	trackB.Add(10, midi.NoteOn(0, 105, 100))

	smfA, smfB := smf.New(), smf.New()
	smfA.Add(trackA)
	smfB.Add(trackB)

	merger.Merge(smfA, smfB)

	var got []string
	for _, msg := range smfA.Tracks[0] {
		if !msg.Message.Is(midi.NoteOnMsg) {
			continue
		}
		var channel, note, velocity uint8
		msg.Message.GetNoteOn(&channel, &note, &velocity)
		got = append(got, fmt.Sprintf("%d, %d, %d, %d", msg.Delta, channel, note, velocity))
	}
	assert.Equal(t, []string{
		"0, 0, 100, 100",
		"5, 0, 101, 100",
		"5, 0, 102, 100",
		"5, 0, 103, 100",
		"5, 0, 104, 100",
		"5, 0, 105, 100",
	}, got)
}
