package recorder

import (
	"bytes"
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

// Recorder records notes into a MIDI SMF file.
type Recorder struct {
	file   *smf.SMF
}

func New() *Recorder {
	return &Recorder{
		file: smf.New(),
	}
}

type Event struct {
	Timestamp time.Time
	Message   midi.Message
}

func At(msg midi.Message, when time.Time) Event {
	return Event{when, msg}
}

func Now(msg midi.Message) Event {
	return Event{time.Now(), msg}
}

// Record, given a channel full of midi messages, writes them into the SMF
// until the channel is closed.
func (r *Recorder) Record(bpm float64, c <-chan Event) {
	ticks := r.file.TimeFormat.(smf.MetricTicks)

	var tr smf.Track
	tr.Add(0, smf.MetaTempo(bpm))

	var lastNano int64
	for msg := range c {
		thisNano := msg.Timestamp.UnixNano()
		if lastNano == 0 {
			lastNano = thisNano
		}
		deltaNano := thisNano - lastNano
		deltaTicks := ticks.Ticks(bpm, time.Duration(deltaNano))
		lastNano = thisNano
		tr.Add(deltaTicks, msg.Message)
	}
	tr.Close(960)
	r.file.Add(tr)
}

func (r *Recorder) SMF() (*smf.SMF, error) {
	if len(r.file.Tracks) == 0 || !r.file.Tracks[0].IsClosed() {
		return nil, fmt.Errorf("recorder is not closed")
	}
	return r.file, nil
}

func (r *Recorder) Bytes() ([]byte, error) {
	if !r.file.Tracks[0].IsClosed() {
		return nil, fmt.Errorf("recorder is not closed")
	}
	var buf bytes.Buffer
	_, err := r.file.WriteTo(&buf)
	if err != nil {
		return nil, fmt.Errorf("error writing smf to bytebuffer: %w", err)
	}
	return buf.Bytes(), nil
}
