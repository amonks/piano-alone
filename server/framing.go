package server

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// checkFraming refuses a part in which anything is declared that the
// part does not hold: a track, a system exclusive, or a meta event
// longer than the bytes left for it, or more tracks than there are.
//
// The MIDI parser allocates what an event declares before it reads the
// event, so a system exclusive claiming a quarter of a gigabyte is a
// quarter of a gigabyte allocated, from a part of twenty-eight bytes
// that anyone with a player's socket may send. Running out of memory
// is not a panic, and nothing recovers from it. After this check, a
// length the parser allocates is a length the part holds, which is at
// most maxPartBytes.
//
// That holds only while this walk and the parser agree on where each
// event starts. So where the parser reads a file some other way than
// it is laid out — it reads a track to its end-of-track event whatever
// length the track gave, six bytes of header whatever length the
// header gave, the bytes after a chunk it does not know as an event,
// and a length of any number of bytes — the part is refused rather
// than read two ways. The page's recorder writes none of those.
func checkFraming(bs []byte) error {
	// The header chunk: its name, its length, and the six bytes that
	// length has to say — the format, the number of tracks, and the
	// time division.
	const headerBytes = 4 + 4 + 6
	if len(bs) < headerBytes || string(bs[:4]) != "MThd" || binary.BigEndian.Uint32(bs[4:8]) != 6 {
		return errors.New("not a MIDI file: it has no header")
	}
	declared := int(binary.BigEndian.Uint16(bs[10:12]))
	tracks := 0
	for rest := bs[headerBytes:]; len(rest) > 0; tracks++ {
		if len(rest) < 8 {
			return errors.New("a chunk's header is cut short")
		}
		kind, size := string(rest[:4]), binary.BigEndian.Uint32(rest[4:8])
		rest = rest[8:]
		if kind != "MTrk" {
			return fmt.Errorf("it has a %q chunk, which is not a track", kind)
		}
		if uint64(size) > uint64(len(rest)) {
			return fmt.Errorf("a track declares %d bytes with %d left in the part", size, len(rest))
		}
		if err := checkTrackFraming(rest[:size]); err != nil {
			return err
		}
		rest = rest[size:]
	}
	if tracks != declared {
		return fmt.Errorf("its header declares %d tracks and it holds %d", declared, tracks)
	}
	return nil
}

// checkTrackFraming walks one track's events to its end-of-track,
// which has to be the last thing in it. An event is a delta time and
// then a channel message — a status byte, or the last one still
// running, and its data bytes — or a system exclusive or a meta event,
// each carrying its own length.
func checkTrackFraming(track []byte) error {
	// The running status: the last channel status, which a message
	// without one of its own repeats. A system exclusive or a meta
	// event clears it.
	var status byte
	for len(track) > 0 {
		_, rest, err := readVLQ(track)
		if err != nil {
			return fmt.Errorf("a delta time %w", err)
		}
		if len(rest) == 0 {
			return errors.New("an event is cut short")
		}
		first := rest[0]
		rest = rest[1:]

		var data int
		switch {
		case first == 0xF0 || first == 0xF7 || first == 0xFF:
			status = 0
			kind := "system exclusive"
			var metaType byte
			if first == 0xFF {
				if len(rest) == 0 {
					return errors.New("an event is cut short")
				}
				kind, metaType, rest = "meta event", rest[0], rest[1:]
			}
			length, after, err := readVLQ(rest)
			if err != nil {
				return fmt.Errorf("a %s's length %w", kind, err)
			}
			if uint64(length) > uint64(len(after)) {
				return fmt.Errorf("a %s declares %d bytes with %d left in its track", kind, length, len(after))
			}
			rest = after[length:]
			if first == 0xFF && metaType == 0x2F {
				if len(rest) > 0 {
					return errors.New("a track goes on past its end")
				}
				return nil
			}
			track = rest
			continue
		case first >= 0x80 && first < 0xF0:
			status = first
			data = channelDataBytes(status)
		case first < 0x80 && status != 0:
			// The first data byte of a message under the running status.
			data = channelDataBytes(status) - 1
		default:
			return fmt.Errorf("an event starts with %#x, which starts none", first)
		}
		if len(rest) < data {
			return errors.New("an event is cut short")
		}
		track = rest[data:]
	}
	return errors.New("a track has no end")
}

// channelDataBytes is how many data bytes follow a channel status: one
// for a program change or a channel pressure, two for the rest.
func channelDataBytes(status byte) int {
	if kind := status >> 4; kind == 0xC || kind == 0xD {
		return 1
	}
	return 2
}

// readVLQ reads a variable-length quantity, the form a MIDI file
// writes its delta times and lengths in: seven bits a byte, the most
// significant first, every byte but the last with its top bit set. The
// format gives one at most four bytes. The error completes a sentence
// naming what was being read.
func readVLQ(bs []byte) (value uint32, rest []byte, err error) {
	for i, b := range bs[:min(len(bs), 4)] {
		value = value<<7 | uint32(b&0x7F)
		if b&0x80 == 0 {
			return value, bs[i+1:], nil
		}
	}
	return 0, nil, errors.New("is cut short, or longer than four bytes")
}
