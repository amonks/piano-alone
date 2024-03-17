package recorder

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

// Recorder records notes into a MIDI SMF file.
type Recorder struct {
	bpm      float64
	lastNano int64
	ticks    smf.MetricTicks
	track    smf.Track
	file     *smf.SMF
	mu       sync.Mutex
}

func New(bpm float64) *Recorder {
	file := smf.New()
	r := &Recorder{
		file:  file,
		bpm:   bpm,
		ticks: file.TimeFormat.(smf.MetricTicks),
	}
	r.track.Add(0, smf.MetaTempo(bpm))
	return r
}

type Event struct {
	Timestamp time.Time
	Message   midi.Message
}

func (m Event) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func EventFromBytes(bs []byte) Event {
	buf := bytes.NewReader(bs)
	dec := gob.NewDecoder(buf)
	var m Event
	if err := dec.Decode(&m); err != nil {
		panic(err)
	}
	return m
}

func At(msg midi.Message, when time.Time) Event {
	return Event{when, msg}
}

func Now(msg midi.Message) Event {
	return Event{time.Now(), msg}
}

func (r *Recorder) Record(msg Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.track.IsClosed() {
		return
	}

	thisNano := msg.Timestamp.UnixNano()
	if r.lastNano == 0 {
		r.lastNano = thisNano
	}
	deltaNano := thisNano - r.lastNano
	deltaTicks := r.ticks.Ticks(r.bpm, time.Duration(deltaNano))
	r.lastNano = thisNano
	r.track.Add(deltaTicks, msg.Message)
}

func (r *Recorder) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.track.Close(960)
	r.file.Add(r.track)
}

func (r *Recorder) SMF() (*smf.SMF, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.track.IsClosed() {
		return nil, fmt.Errorf("recorder is not closed")
	}
	return r.file, nil
}

func (r *Recorder) Bytes() ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.track.IsClosed() {
		return nil, fmt.Errorf("recorder is not closed")
	}
	var buf bytes.Buffer
	_, err := r.file.WriteTo(&buf)
	if err != nil {
		return nil, fmt.Errorf("error writing smf to bytebuffer: %w", err)
	}
	return buf.Bytes(), nil
}
