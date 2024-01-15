package abstrack

import (
	"bytes"

	"gitlab.com/gomidi/midi/v2/smf"
)

// An AbsTrack is semantically equivalent to an smf.Track,
// but it represents time using absolute timestamps rather
// than a delta series.
type AbsTrack struct {
	Events []AbsEvent
}

func New() *AbsTrack {
	return &AbsTrack{}
}

// An AbsEvent is semantically equivalent to an smf.Event,
// but it represents time using an absolute timestamp
// rather than a delta.
type AbsEvent struct {
	Timestamp uint32
	Message   smf.Message
}

func FromSMF(src smf.Track) *AbsTrack {
	dst := &AbsTrack{}
	var lastStamp uint32
	for _, event := range src {
		stamp := lastStamp + event.Delta
		dst.Events = append(dst.Events, AbsEvent{
			Timestamp: stamp,
			Message:   event.Message,
		})
		lastStamp = stamp
	}
	return dst
}

func (at *AbsTrack) ToSMF() smf.Track {
	dst := make(smf.Track, len(at.Events))
	var lastTimestamp uint32
	for i, event := range at.Events {
		dst[i] = smf.Event{
			Delta:   event.Timestamp - lastTimestamp,
			Message: event.Message,
		}
		lastTimestamp = event.Timestamp
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
	return &AbsTrack{dst}
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
