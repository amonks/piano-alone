package abstrack_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
)

func TestRoundTrip(t *testing.T) {
	example, err := smf.ReadFile("../example.mid")
	require.NoError(t, err)

	converted := abstrack.FromSMF(example, 0).ToSMF()
	assert.Equal(t, example.Tracks[0], converted)
}

func TestMerger(t *testing.T) {
	var trackA, trackB smf.Track

	trackA.Add(0, smf.MetaTempo(120))
	trackA.Add(0, midi.NoteOn(0, 100, 100))
	trackA.Add(1000, midi.NoteOn(0, 102, 100))
	trackA.Add(1000, midi.NoteOn(0, 104, 100))
	trackA.Add(500, smf.EOT)

	trackB.Add(0, smf.MetaTempo(120))
	trackB.Add(500, midi.NoteOn(0, 101, 100))
	trackB.Add(1000, midi.NoteOn(0, 103, 100))
	trackB.Add(1000, midi.NoteOn(0, 105, 100))
	trackB.Add(500, smf.EOT)

	smfA, smfB := smf.New(), smf.New()
	smfA.Add(trackA)
	smfB.Add(trackB)
	absA, absB := abstrack.FromSMF(smfA, 0), abstrack.FromSMF(smfB, 0)

	absMerged := absA.Merge(absB)
	merged := absMerged.ToSMF()

	var msgs []string
	for _, ev := range merged {
		msgs = append(msgs, fmt.Sprintf("%d: %s", ev.Delta, ev.Message.String()))
	}
	assert.Equal(t, []string{
		"0: MetaTempo bpm: 120.00",
		"0: NoteOn channel: 0 key: 100 velocity: 100",
		"250: NoteOn channel: 0 key: 101 velocity: 100",
		"250: NoteOn channel: 0 key: 102 velocity: 100",
		"250: NoteOn channel: 0 key: 103 velocity: 100",
		"250: NoteOn channel: 0 key: 104 velocity: 100",
		"250: NoteOn channel: 0 key: 105 velocity: 100",
		"250: MetaEndOfTrack",
	}, msgs)
}
