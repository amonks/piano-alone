package merger

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
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

	absDst, absSrc := abstrack.FromSMF(dst.Tracks[0]), abstrack.FromSMF(src.Tracks[0])
	absDst.Merge(absSrc)
	dst.Tracks[0] = absDst.ToSMF()

	return nil
}
