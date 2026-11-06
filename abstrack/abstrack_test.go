package abstrack_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/songs"
)

func TestRoundTrip(t *testing.T) {
	example := songs.PreludeBergamasqueSMF
	converted := abstrack.FromSMF(example, 0).ToSMF(example.TimeFormat.(smf.MetricTicks))
	assert.Equal(t, example.TimeFormat, converted.TimeFormat)
	assert.Equal(t, example.Tracks[0], converted.Tracks[0])
}

// The file ToSMF writes declares the ticks a quarter note its deltas
// are counted in, so a MIDI player reading it sounds every event when
// the track says. ToSMF counted its deltas at 480 whatever the header
// said, and the file a caller made for them said 960, so it played at
// twice its speed.
func TestToSMFPlaysAtItsOwnTimes(t *testing.T) {
	track := abstrack.New()
	track.Append(
		abstrack.AbsEvent{Timestamp: 0, Message: smf.MetaTempo(120)},
		abstrack.AbsEvent{Timestamp: 0, Message: smf.Message(midi.NoteOn(0, 60, 100))},
		abstrack.AbsEvent{Timestamp: 500 * time.Millisecond, Message: smf.Message(midi.NoteOff(0, 60))},
		abstrack.AbsEvent{Timestamp: 3 * time.Second, Message: smf.Message(midi.NoteOn(0, 64, 100))},
		abstrack.AbsEvent{Timestamp: 3250 * time.Millisecond, Message: smf.Message(midi.NoteOff(0, 64))},
	)
	for _, ticks := range []smf.MetricTicks{96, 480, 960} {
		var b bytes.Buffer
		if _, err := track.ToSMF(ticks).WriteTo(&b); err != nil {
			t.Fatal(err)
		}
		var got []time.Duration
		rd := smf.ReadTracksFrom(&b).Do(func(ev smf.TrackEvent) {
			var ch, key, vel uint8
			if ev.Message.GetNoteStart(&ch, &key, &vel) || ev.Message.GetNoteEnd(&ch, &key) {
				got = append(got, time.Duration(ev.AbsMicroSeconds)*time.Microsecond)
			}
		})
		if err := rd.Error(); err != nil {
			t.Fatal(err)
		}
		want := []time.Duration{0, 500 * time.Millisecond, 3 * time.Second, 3250 * time.Millisecond}
		assert.Equal(t, want, got, "at %d ticks a quarter note", ticks)
	}
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
	merged := absMerged.ToSMF(smfA.TimeFormat.(smf.MetricTicks))

	var msgs []string
	for _, ev := range merged.Tracks[0] {
		msgs = append(msgs, fmt.Sprintf("%d: %s", ev.Delta, ev.Message.String()))
	}
	assert.Equal(t, []string{
		"0: MetaTempo bpm: 120.00",
		"0: NoteOn channel: 0 key: 100 velocity: 100",
		"500: NoteOn channel: 0 key: 101 velocity: 100",
		"500: NoteOn channel: 0 key: 102 velocity: 100",
		"500: NoteOn channel: 0 key: 103 velocity: 100",
		"500: NoteOn channel: 0 key: 104 velocity: 100",
		"500: NoteOn channel: 0 key: 105 velocity: 100",
		"500: MetaEndOfTrack",
	}, msgs)
}
