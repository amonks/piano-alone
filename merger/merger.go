package merger

import (
	"fmt"
	"sort"

	"gitlab.com/gomidi/midi/v2/smf"
)

// Merge combines two single-track SMF files into a single file, overwriting
// its first argument.
func Merge(dst, src *smf.SMF) error {
	if l := len(dst.Tracks); l != 1 {
		return fmt.Errorf("dst has %d tracks, should have 1", l)
	}
	if l := len(src.Tracks); l != 1 {
		return fmt.Errorf("src has %d tracks, should have 1", l)
	}

	// In the SMF file, event timestamps are expressed as deltas from the
	// previous event. Therefore, to interleave the events of multiple
	// tracks, we have to recalculate the deltas.
	//
	// 1. First, we'll go through both tracks, adding each event to a list
	//   (`stamped`) of events with absolute timestamps rather than deltas.
	// 2. We'll sort the list by those absolute timestamps
	// 3. We'll convert the list into a new track, converting the absolute
	//    timestamps with the deltas mandated by SMF.
	var stamped []stampedMsg
	stamped = append(stamped, convertToStamped(src.Tracks[0])...)
	stamped = append(stamped, convertToStamped(dst.Tracks[0])...)
	sort.Slice(stamped, func(a, b int) bool {
		return stamped[a].stamp < stamped[b].stamp
	})

	// TODO: dedupe tempo/meter messages
	var track smf.Track
	var lastStamp uint32
	for _, stamped := range stamped {
		track = append(track, smf.Event{
			Delta:   stamped.stamp - lastStamp,
			Message: stamped.msg,
		})
		lastStamp = stamped.stamp
	}

	dst.Tracks[0] = track
	return nil
}

type stampedMsg struct {
	msg   smf.Message
	stamp uint32
}

func convertToStamped(track smf.Track) []stampedMsg {
	var stamped []stampedMsg
	var lastStamp uint32
	for _, msg := range track {
		stamp := lastStamp + msg.Delta
		stamped = append(stamped, stampedMsg{
			stamp: stamp,
			msg:   msg.Message,
		})
		lastStamp = stamp
	}
	return stamped
}
