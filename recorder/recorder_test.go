package recorder_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gomidi/midi/v2"
	"monks.co/piano-alone/recorder"
)

func TestRecorder(t *testing.T) {
	rec := recorder.New()
	c := make(chan recorder.Event)
	done := make(chan struct{})
	go func() { rec.Record(120, c); done <- struct{}{} }()
	for _, ev := range []struct {
		i   int
		msg midi.Message
	}{
		{1, midi.NoteOn(0, 81, 100)},
		{2, midi.NoteOn(0, 82, 100)},
		{3, midi.NoteOn(0, 83, 100)},
	} {
		s := fmt.Sprintf("2006-01-02T15:04:0%dZ", ev.i)
		at, err := time.Parse(time.RFC3339, s)
		require.NoError(t, err)
		c <- recorder.At(ev.msg, at)
	}
	close(c)
	<-done
	s, err := rec.SMF()
	require.NoError(t, err)

	track := s.Tracks[0]
	var msgs []string
	for _, ev := range track {
		msgs = append(msgs, ev.Message.String())
	}
	assert.Equal(t, []string{
		"MetaTempo bpm: 120.00",
		"NoteOn channel: 0 key: 81 velocity: 100",
		"NoteOn channel: 0 key: 82 velocity: 100",
		"NoteOn channel: 0 key: 83 velocity: 100",
		"MetaEndOfTrack",
	}, msgs)
}
