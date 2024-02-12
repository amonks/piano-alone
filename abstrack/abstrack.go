package abstrack

import (
	"bytes"
	"sort"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
)

// An AbsTrack is semantically equivalent to an smf.Track,
// but it represents time using absolute timestamps rather
// than a delta series.
type AbsTrack struct {
	Events []AbsEvent
	bpm    float64
}

func New() *AbsTrack {
	return &AbsTrack{}
}

func (t *AbsTrack) BPM() float64 {
	if t.bpm != 0 {
		return t.bpm
	}
	var bpm float64
	for _, ev := range t.Events {
		if !ev.Message.Is(smf.MetaTempoMsg) {
			continue
		}
		ev.Message.GetMetaTempo(&bpm)
		break
	}
	t.bpm = bpm
	return bpm
}

// An AbsEvent is semantically equivalent to an smf.Event, but it represents
// time using an absolute nanosecond timestamp rather than a delta.
type AbsEvent struct {
	Timestamp time.Duration
	Message   smf.Message
}

func FromSMF(file *smf.SMF, trackIndex int) *AbsTrack {
	var (
		src = file.Tracks[trackIndex]
		dst = make([]AbsEvent, len(src))

		ticks               = file.TimeFormat.(smf.MetricTicks)
		bpm   float64       = 0
		now   time.Duration = 0
	)
	for i, event := range src {
		var dur time.Duration
		if event.Delta > 0 {
			dur = ticks.Duration(bpm, event.Delta)
		}
		now = now + dur
		dst[i] = AbsEvent{
			Timestamp: now,
			Message:   event.Message,
		}

		if event.Message.Is(smf.MetaTempoMsg) {
			event.Message.GetMetaTempo(&bpm)
		}
	}
	return &AbsTrack{Events: dst}
}

func (at *AbsTrack) ToSMF() smf.Track {
	var (
		src = at.Events
		dst = make(smf.Track, len(src))

		ticks               = smf.MetricTicks(480)
		bpm   float64       = 0
		now   time.Duration = 0
	)
	for i, event := range src {
		if now > event.Timestamp {
			panic("something is amiss")
		}
		delta := event.Timestamp - now
		dst[i] = smf.Event{
			Delta:   ticks.Ticks(bpm, delta),
			Message: event.Message,
		}
		now = event.Timestamp
		if event.Message.Is(smf.MetaTempoMsg) {
			event.Message.GetMetaTempo(&bpm)
		}
	}
	return dst
}

type CountedKey struct {
	Key   uint8
	Count int
}

func (at *AbsTrack) CountNotes() []CountedKey {
	keyset := map[uint8]int{}
	var channel, key, velocity uint8
	for _, ev := range at.Events {
		if !ev.Message.GetNoteOn(&channel, &key, &velocity) {
			continue
		}
		keyset[key] += 1
	}
	var keys []CountedKey
	for key, count := range keyset {
		keys = append(keys, CountedKey{key, count})
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Count > keys[j].Count })
	return keys
}

func (at *AbsTrack) Select(notes []uint8) *AbsTrack {
	if at.Events[0].Timestamp < 0 {
		panic("invalid track")
	}
	noteSet := map[uint8]struct{}{}
	for _, note := range notes {
		noteSet[note] = struct{}{}
	}
	dst := New()
	var c, key, vel uint8
	for _, ev := range at.Events {
		if ev.Message.GetNoteStart(&c, &key, &vel) {
			if _, has := noteSet[key]; has {
				dst.Events = append(dst.Events, ev)
			}
		} else if ev.Message.GetNoteEnd(&c, &key) {
			if _, has := noteSet[key]; has {
				dst.Events = append(dst.Events, ev)
			}
		} else {
			dst.Events = append(dst.Events, ev)
		}
	}
	return dst
}

// Merge combines two AbsTracks, retaining sort order and deduplicating. The
// resulting track has only one MetaEndOfTrackMsg. Merge DOES NOT handle tempo
// changes or tracks with different BPM.
func (at *AbsTrack) Merge(from *AbsTrack) *AbsTrack {
	var dst []AbsEvent
	intoIdx, fromIdx := 0, 0
	for intoIdx < len(at.Events) && fromIdx < len(from.Events) {
		intoEvent, fromEvent := at.Events[intoIdx], from.Events[fromIdx]
		if intoEvent.Message.Type() == smf.MetaEndOfTrackMsg {
			intoIdx++
		} else if fromEvent.Message.Type() == smf.MetaEndOfTrackMsg {
			fromIdx++
		} else if intoEvent.Timestamp < fromEvent.Timestamp {
			dst = append(dst, intoEvent)
			intoIdx++
		} else if fromEvent.Timestamp < intoEvent.Timestamp {
			dst = append(dst, fromEvent)
			fromIdx++
		} else if intoEvent.Is(fromEvent) {
			intoIdx++
		} else if intoEvent.Less(fromEvent) {
			dst = append(dst, intoEvent)
			intoIdx++
		} else {
			dst = append(dst, fromEvent)
			fromIdx++
		}
	}
	if intoIdx < len(at.Events) {
		dst = append(dst, at.Events[intoIdx:]...)
	} else if fromIdx < len(from.Events) {
		dst = append(dst, from.Events[fromIdx:]...)
	}
	return &AbsTrack{Events: dst}
}

// Less returns true if `ev` should appear earlier in an
// SMF file than `other`.
func (ev AbsEvent) Less(other AbsEvent) bool {
	if ev.Timestamp != other.Timestamp {
		return ev.Timestamp < other.Timestamp
	}

	// This is intentionally reverse-lexical: we want
	// meta messages (first byte 0xFF) to sort before
	// non-meta messages.
	return bytes.Compare(ev.Message, other.Message) == 1
}

func (ev AbsEvent) Is(other AbsEvent) bool {
	return bytes.Equal(ev.Message, other.Message)
}
