package server

import (
	"bytes"
	"encoding/binary"
	"runtime"
	"strings"
	"testing"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	"monks.co/piano-alone/recorder"
)

// smfHeader is a MIDI file's header: format 0, one track, 480 ticks a
// quarter note.
const smfHeader = "MThd\x00\x00\x00\x06\x00\x00\x00\x01\x01\xe0"

// trackChunk is a track holding events, under an honest length.
func trackChunk(events string) string {
	var size [4]byte
	binary.BigEndian.PutUint32(size[:], uint32(len(events)))
	return "MTrk" + string(size[:]) + events
}

// allocated is how many bytes f allocates, whether or not they are
// still held when it returns.
func allocated(f func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// The tempo every part opens with, as the page's recorder writes it,
// and the end every track closes with.
const (
	tempoEvent = "\x00\xff\x51\x03\x07\xa1\x20"
	endOfTrack = "\x00\xff\x2f\x00"
)

// A part says how long each of its events is, and the MIDI parser
// allocates what an event says before reading it. A system exclusive
// claiming a quarter of a gigabyte was a quarter of a gigabyte
// allocated, from a part of twenty-eight bytes, by anyone with a
// player's socket — and running out of memory is not a panic, so
// nothing recovered from it. A part is refused when any length in it
// runs past the bytes it holds, before the parser reads it.
func TestAPartDeclaringMoreThanItHoldsIsRefused(t *testing.T) {
	const huge = "\xff\xff\xff\x7f" // 0x0FFFFFFF, the most a length can say
	for name, part := range map[string]string{
		"a system exclusive":            smfHeader + trackChunk("\x00\xf0"+huge),
		"a system exclusive, continued": smfHeader + trackChunk(tempoEvent+"\x00\xf7"+huge),
		"a meta event":                  smfHeader + trackChunk(tempoEvent+"\x00\xff\x01"+huge),
		"a track":                       smfHeader + "MTrk\x0f\xff\xff\xff" + tempoEvent + endOfTrack,
		// The bytes it claims are in the part, but not in its track.
		"a system exclusive running into the next track": smfHeader +
			trackChunk(tempoEvent+"\x00\xf0\x05\x01\x02\x03") + trackChunk(tempoEvent+endOfTrack),
	} {
		t.Run(name, func(t *testing.T) {
			var err error
			spent := allocated(func() { err = checkPart([]byte(part)) })
			if err == nil {
				t.Errorf("the part was taken")
			} else if !strings.Contains(err.Error(), "declares") {
				t.Errorf("refused as %q, want it refused for what it declares", err)
			}
			if spent > 1<<20 {
				t.Errorf("reading a %d-byte part allocated %d bytes", len(part), spent)
			}
		})
	}
}

// The check reads a part the way the parser will, so that a length it
// passes is the length the parser allocates. Where the parser reads a
// file some other way than it is laid out — more length bytes than the
// format has, a track that runs past the length it gave or stops short
// of it, a chunk that is not a track — the part is refused rather than
// read two ways.
func TestAPartTheParserWouldReadAnotherWayIsRefused(t *testing.T) {
	for name, part := range map[string]string{
		"a length of five bytes":       smfHeader + trackChunk(tempoEvent+"\x00\xff\x01\x80\x80\x80\x80\x00"+endOfTrack),
		"a delta time of five bytes":   smfHeader + trackChunk(tempoEvent+"\x80\x80\x80\x80\x00\x90\x3c\x64"+endOfTrack),
		"a length cut short":           smfHeader + trackChunk(tempoEvent+"\x00\xf0\x81"),
		"a track with no end":          smfHeader + trackChunk(tempoEvent+"\x00\x91\x3c\x64"),
		"a track past its end":         smfHeader + trackChunk(tempoEvent+endOfTrack+"\x00\xf0"+"\x00"),
		"an empty track":               smfHeader + trackChunk("") + trackChunk(tempoEvent+endOfTrack),
		"a chunk that is not a track":  smfHeader + "XFIH\x00\x00\x00\x02\x00\x00" + trackChunk(tempoEvent+endOfTrack),
		"a chunk header cut short":     smfHeader + trackChunk(tempoEvent+endOfTrack) + "MTrk\x00",
		"a header of another length":   "MThd\x00\x00\x00\x08\x00\x00\x00\x01\x01\xe0\x00\x00" + trackChunk(tempoEvent+endOfTrack),
		"a data byte with no status":   smfHeader + trackChunk(tempoEvent+"\x00\x3c\x64"+endOfTrack),
		"a status no file event has":   smfHeader + trackChunk(tempoEvent+"\x00\xf8"+endOfTrack),
		"a channel message cut short":  smfHeader + trackChunk(tempoEvent+"\x00\x91\x3c"),
		"a running status cut short":   smfHeader + trackChunk(tempoEvent+"\x00\x91\x3c\x64\x00\x3e"),
		"a meta event with no type":    smfHeader + trackChunk(tempoEvent+"\x00\xff"),
		"nothing after the delta time": smfHeader + trackChunk(tempoEvent+"\x00"),
	} {
		if err := checkFraming([]byte(part)); err == nil {
			t.Errorf("%s: the part's framing passed", name)
		}
	}
}

// What the page records passes, and so does everything else MIDI lets
// a track carry at its honest length: a meta event, a system
// exclusive, a channel message of one data byte, and notes under a
// running status, which the page's recorder writes whenever one
// message follows another of its kind.
func TestAPartWithHonestLengthsPasses(t *testing.T) {
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(120))
	tr.Add(0, smf.MetaText("a part"))
	tr.Add(0, midi.SysEx([]byte{0x7e, 0x7f, 0x09, 0x01}))
	tr.Add(0, midi.NoteOn(recorder.Channel, 60, 100))
	tr.Add(240, midi.NoteOn(recorder.Channel, 64, 100)) // running status
	tr.Add(240, midi.NoteOff(recorder.Channel, 60))
	tr.Add(0, midi.NoteOff(recorder.Channel, 64)) // running status
	tr.Add(10, midi.ProgramChange(recorder.Channel, 5))
	tr.Add(10, midi.ProgramChange(recorder.Channel, 6)) // running status, one data byte
	tr.Add(0, smf.MetaText(strings.Repeat("x", 300)))   // a length of two bytes
	tr.Close(0)
	f := smf.New()
	if err := f.Add(tr); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := f.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	everything := b.Bytes()
	if !bytes.Contains(everything, []byte{0x91, 60, 100, 0x81, 0x70, 64, 100}) {
		t.Fatalf("the part holds no running status: % x", everything)
	}

	r := recorder.New(120)
	r.Close()
	silent, err := r.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	for name, part := range map[string][]byte{
		"every kind of event":      everything,
		"the page's part":          partOf(t, 60, 0),
		"the page's part, silent":  silent,
		"a part at its size limit": part(t, maxPartBytes-200),
		"two tracks":               []byte("MThd\x00\x00\x00\x06\x00\x01\x00\x02\x01\xe0" + trackChunk(tempoEvent+endOfTrack) + trackChunk(endOfTrack)),
	} {
		if err := checkPart(part); err != nil {
			t.Errorf("%s: refused: %v", name, err)
		}
	}

	// Read by the parser, the part is what was written: the check and
	// the parser agree on where its events are.
	read, err := readPart(everything)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(read.Tracks[0]), len(tr); got != want {
		t.Errorf("the parser read %d events, want %d", got, want)
	}
}

// Whatever the check passes, the parser reads within the part's own
// size. Every single-byte change to a part that holds each kind of
// event is tried: a change the check and the parser read differently
// would show as a part that passes and then allocates far more than
// it is.
func TestNoPartThatPassesMakesTheParserAllocate(t *testing.T) {
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(120))
	tr.Add(0, midi.SysEx([]byte{0x7e, 0x7f}))
	tr.Add(0, midi.NoteOn(recorder.Channel, 60, 100))
	tr.Add(240, midi.NoteOn(recorder.Channel, 64, 100))
	tr.Add(0, smf.MetaText("hi"))
	tr.Add(240, midi.NoteOff(recorder.Channel, 60))
	tr.Close(0)
	// A second track, so the changes reach the bytes between two.
	var second smf.Track
	second.Add(0, smf.MetaText("x"))
	second.Close(0)
	f := smf.New()
	for _, track := range []smf.Track{tr, second} {
		if err := f.Add(track); err != nil {
			t.Fatal(err)
		}
	}
	var b bytes.Buffer
	if _, err := f.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	base := b.Bytes()

	passed := 0
	for i := range base {
		for v := range 256 {
			mutant := bytes.Clone(base)
			mutant[i] = byte(v)
			if checkFraming(mutant) != nil {
				continue
			}
			passed++
			if spent := allocated(func() { readPart(mutant) }); spent > 1<<20 {
				t.Fatalf("with byte %d set to %#x the part passes, and reading its %d bytes allocated %d",
					i, v, len(mutant), spent)
			}
		}
	}
	if passed < 256 {
		t.Errorf("only %d of the changed parts passed: the check refuses what it should read", passed)
	}
}
