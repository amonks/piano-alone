package server

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/recorder"
)

// loopServer is a server whose game loop a test drives one message at
// a time through handle, with an outbox deep enough that nothing the
// loop says blocks it.
func loopServer(t *testing.T) *Server {
	t.Helper()
	s := newTestServer(t)
	s.state = game.NewState()
	s.partials = map[string][]byte{}
	s.outbox = make(chan *game.Message, 1<<16)
	return s
}

func (s *Server) receive(t *testing.T, msgType game.MessageType, from string, data []byte) {
	t.Helper()
	if err := s.handle(t.Context(), game.NewMessage(msgType, from, data)); err != nil {
		t.Logf("%s from %.12s: %v", msgType, from, err)
	}
}

func (s *Server) join(t *testing.T, fingerprint string) {
	t.Helper()
	s.receive(t, game.MessageTypeJoin, fingerprint, game.JoinMessage{Fingerprint: fingerprint, NoteCapacity: 3}.Bytes())
}

// part is a recorded part: a valid SMF, padded with a text event to at
// least size bytes.
func part(t *testing.T, size int) []byte { return partOf(t, 60, size) }

// partOf is a part playing key, shaped as the player's recorder makes
// one: its tempo first, then the notes.
func partOf(t *testing.T, key uint8, size int) []byte {
	t.Helper()
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(120))
	tr.Add(0, midi.NoteOn(recorder.Channel, key, 100))
	tr.Add(120, midi.NoteOff(recorder.Channel, key))
	if size > 0 {
		tr.Add(0, smf.MetaText(strings.Repeat("x", size)))
	}
	tr.Close(0)
	f := smf.New()
	if err := f.Add(tr); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := f.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// A part is the recording of a player's notes, so it is taken only
// while the performance is being played and only from someone in it.
// The player's socket is public and open in every phase; a part taken
// from anyone at any time sat in memory until a Restart, and was
// merged into the next performance's rendition.
func TestPartsAreTakenOnlyInHeroFromPlayers(t *testing.T) {
	s := loopServer(t)
	s.join(t, "player")
	s.player("player").AssignedNotes = []uint8{60}
	for _, phase := range []game.Phase{game.GamePhaseUninitialized, game.GamePhaseLobby, game.GamePhaseProcessing, game.GamePhasePlayback, game.GamePhaseDone} {
		s.state.Phase = phase
		s.receive(t, game.MessageTypeSubmitPartialTrack, "player", part(t, 0))
		if len(s.partials) != 0 || s.player("player").HasSubmitted {
			t.Fatalf("a part was taken in %s", phase)
		}
	}

	s.state.Phase = game.GamePhaseHero
	s.receive(t, game.MessageTypeSubmitPartialTrack, "stranger", part(t, 0))
	if len(s.partials) != 0 || s.player("stranger") != nil {
		t.Errorf("a part from someone who never joined was taken: %d parts, player %+v", len(s.partials), s.player("stranger"))
	}
	s.receive(t, game.MessageTypeSubmitPartialTrack, "player", part(t, 0))
	if len(s.partials) != 1 || !s.player("player").HasSubmitted {
		t.Errorf("the player's part in Hero was not taken")
	}
}

// What one performance holds is bounded whoever opens sockets: so
// many players, and so many bytes in each one's part — far above
// anything the piece has played (its largest rendition, six players'
// parts merged, is 7 KB). The roster counts what is worth keeping: a
// visitor who has left holding no notes and no part makes room for a
// new one, so tabs opened and closed between performances never fill
// it. Only sockets held open can, and someone refused is still sent
// the state, so their page knows where the performance is.
func TestAPerformanceHoldsBoundedPlayersAndParts(t *testing.T) {
	s := loopServer(t)
	for i := range maxPlayers {
		s.join(t, fmt.Sprintf("ghost%d", i))
		s.receive(t, game.MessageTypeLeave, fmt.Sprintf("ghost%d", i), nil)
	}
	s.join(t, "early")
	if s.player("early") == nil || len(s.state.Players) != maxPlayers {
		t.Fatalf("a join at a roster of leavers: early=%+v, %d players", s.player("early"), len(s.state.Players))
	}

	s = loopServer(t)
	for i := range maxPlayers {
		s.join(t, fmt.Sprintf("p%d", i))
	}
	drainOutbox(s)
	s.join(t, "late")
	if s.player("late") != nil || len(s.state.Players) != maxPlayers {
		t.Errorf("a join past a full room of connected players was taken: %d players", len(s.state.Players))
	}
	if m := drainOutbox(s); len(m) != 1 || m[0].Type != game.MessageTypeState || m[0].Player != "late" {
		t.Errorf("the refused player was sent %+v, want the state", m)
	}
	// Someone already in it rejoins at the cap — a page reload.
	s.receive(t, game.MessageTypeLeave, "p0", nil)
	s.join(t, "p0")
	if p := s.player("p0"); p == nil || p.ConnectionState != game.ConnectionStateConnected {
		t.Errorf("a player rejoining at the cap: %+v", p)
	}
	// A player who left holding notes may come back to play them, so
	// their record is not room for a stranger.
	s.player("p1").AssignedNotes = []uint8{60}
	for i := 2; i < maxPlayers; i++ {
		s.player(fmt.Sprintf("p%d", i)).HasSubmitted = true
	}
	s.receive(t, game.MessageTypeLeave, "p1", nil)
	s.join(t, "stranger")
	if s.player("stranger") != nil || s.player("p1") == nil {
		t.Error("a stranger took the place of a player who left holding notes")
	}

	s.state.Phase = game.GamePhaseHero
	s.player("p0").AssignedNotes = []uint8{60}
	s.receive(t, game.MessageTypeSubmitPartialTrack, "p0", part(t, maxPartBytes))
	if len(s.partials) != 0 {
		t.Error("a part over the per-part cap was taken")
	}
	s.receive(t, game.MessageTypeSubmitPartialTrack, "p0", part(t, maxPartBytes-200))
	if len(s.partials) != 1 {
		t.Error("a part just under the per-part cap was refused")
	}
}

// drainOutbox is what the loop has said since it was last drained.
func drainOutbox(s *Server) []*game.Message {
	var out []*game.Message
	for {
		select {
		case m := <-s.outbox:
			out = append(out, m)
		default:
			return out
		}
	}
}

// A part the merge could not take — no tracks, a time format other
// than metric ticks, time passing before a tempo is set (which the
// merge would turn into an infinite duration), a timeline that runs
// backwards — is refused when it is submitted, so no one part can stop
// the rendition being assembled.
func TestPartsTheMergeCannotTakeAreRefused(t *testing.T) {
	smfOf := func(ticks smf.MetricTicks, tracks ...smf.Track) []byte {
		f := smf.New()
		f.TimeFormat = ticks
		for _, tr := range tracks {
			if err := f.Add(tr); err != nil {
				t.Fatal(err)
			}
		}
		var b bytes.Buffer
		if _, err := f.WriteTo(&b); err != nil {
			t.Fatal(err)
		}
		return b.Bytes()
	}
	var tempoless smf.Track
	tempoless.Add(0, midi.NoteOn(0, 61, 100))
	tempoless.Add(120, midi.NoteOff(0, 61))
	tempoless.Close(0)
	// At one tick a quarter note and four beats a minute, the longest
	// delta MIDI can carry is about 130 years; three of them overflow.
	var slowAndLong smf.Track
	slowAndLong.Add(0, smf.MetaTempo(4))
	for range 3 {
		slowAndLong.Add(1<<28-1, midi.NoteOn(0, 62, 100))
	}
	slowAndLong.Close(0)
	smpte := []byte("MThd\x00\x00\x00\x06\x00\x00\x00\x01\xe7\x28MTrk\x00\x00\x00\x04\x00\xff\x2f\x00")

	for name, bad := range map[string][]byte{
		"no tracks":            []byte("MThd\x00\x00\x00\x06\x00\x01\x00\x00\x01\xe0"),
		"no tempo":             smfOf(960, tempoless),
		"a timeline overflows": smfOf(1, slowAndLong),
		"SMPTE time":           smpte,
		"not MIDI":             []byte("hello"),
	} {
		s := loopServer(t)
		s.join(t, "a")
		s.player("a").AssignedNotes = []uint8{61, 62}
		s.state.Phase = game.GamePhaseHero
		s.receive(t, game.MessageTypeSubmitPartialTrack, "a", bad)
		if len(s.partials) != 0 || s.player("a").HasSubmitted {
			t.Errorf("%s: the part was taken", name)
		}
	}
}

// Beginning a performance begins it with the audience in the room and
// nothing else of the last one: its parts are not merged into this
// one's rendition, players who have left are forgotten, and those
// still connected start with no notes, no tutorial, and no part.
func TestBeginningAPerformanceClearsTheLastOne(t *testing.T) {
	s := loopServer(t)
	s.join(t, "stays")
	s.join(t, "left")
	s.receive(t, game.MessageTypeLeave, "left", nil)
	s.state.Phase = game.GamePhaseHero
	for _, f := range []string{"stays", "left"} {
		p := s.player(f)
		p.AssignedNotes, p.HasStartedTutorial, p.HasCompletedTutorial = []uint8{60}, true, true
		s.receive(t, game.MessageTypeSubmitPartialTrack, f, part(t, 0))
	}
	if len(s.partials) != 2 {
		t.Fatalf("setup: %d parts", len(s.partials))
	}
	s.state.Rendition = []byte("the last rendition")
	s.sentSwitchToVideoModal = true

	conf := &game.Configuration{PerformanceID: "next", Title: "Next", Composer: "Somebody"}
	s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())

	if s.state.Phase != game.GamePhaseLobby {
		t.Fatalf("phase after begin = %s", s.state.Phase)
	}
	if len(s.partials) != 0 {
		t.Errorf("the last performance's parts survived: %d", len(s.partials))
	}
	if s.player("left") != nil {
		t.Error("a player who had left survived into the next performance")
	}
	want := game.Player{Fingerprint: "stays", NoteCapacity: 3, ConnectionState: game.ConnectionStateConnected}
	if p := s.player("stays"); p == nil || p.Fingerprint != want.Fingerprint || p.NoteCapacity != want.NoteCapacity ||
		p.ConnectionState != want.ConnectionState || len(p.AssignedNotes) != 0 || p.HasStartedTutorial || p.HasCompletedTutorial || p.HasSubmitted {
		t.Errorf("the connected player after begin = %+v, want a fresh %+v", p, want)
	}
	if s.state.Rendition != nil || s.sentSwitchToVideoModal {
		t.Error("the last performance's rendition or video cue survived")
	}
}

// The parts are kept as the bytes each player sent and parsed when the
// performance is assembled: the rendition is every part merged, saved
// against the performance, and a part that cannot be merged costs only
// itself.
func TestTheRenditionMergesThePartsTaken(t *testing.T) {
	s := loopServer(t)
	conf := &game.Configuration{PerformanceID: "tonight", Title: "Tonight", Composer: "Somebody", Score: partOf(t, 60, 0)}
	if err := s.db.SchedulePerformance(t.Context(), &game.Performance{Configuration: conf}); err != nil {
		t.Fatal(err)
	}
	s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())
	s.join(t, "a")
	s.join(t, "b")
	s.player("a").AssignedNotes = []uint8{60}
	s.player("b").AssignedNotes = []uint8{64}
	s.state.Phase = game.GamePhaseHero
	s.receive(t, game.MessageTypeSubmitPartialTrack, "a", partOf(t, 60, 0))
	s.receive(t, game.MessageTypeSubmitPartialTrack, "b", partOf(t, 64, 0))
	// Parts that could not have been taken, planted past the intake:
	// assembly skips each rather than losing the rendition to it.
	var tempoless smf.Track
	tempoless.Add(0, midi.NoteOn(0, 61, 100))
	tempoless.Add(120, midi.NoteOff(0, 61))
	tempoless.Close(0)
	noTempo := smf.New()
	if err := noTempo.Add(tempoless); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := noTempo.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	s.partials["planted"] = b.Bytes()
	s.partials["garbage"] = []byte("not MIDI")
	s.receive(t, game.MessageTypeAdvancePhase, roleConductor, nil)

	if s.state.Phase != game.GamePhasePlayback {
		t.Fatalf("phase after assembly = %s", s.state.Phase)
	}
	saved, err := s.db.GetMIDIFile(t.Context(), "tonight")
	if err != nil {
		t.Fatal(err)
	}
	rendition, err := smf.ReadFrom(bytes.NewReader(saved))
	if err != nil {
		t.Fatalf("the saved rendition does not parse: %v", err)
	}
	played := map[uint8]bool{}
	for _, tr := range rendition.Tracks {
		for _, ev := range tr {
			var ch, key, vel uint8
			if ev.Message.GetNoteStart(&ch, &key, &vel) {
				played[key] = true
			}
		}
	}
	if !played[60] || !played[64] || len(played) != 2 {
		t.Errorf("the rendition plays %v, want both parts' keys", played)
	}
}

// A completed performance is not begun again — not from the conductor's
// CLI, which is answered 409, and not from the conductor's socket,
// which the loop refuses — since performing it would overwrite the
// rendition its audience played. Playing the piece again is scheduling
// a new performance; deleting one is the destructive verb.
func TestACompletedPerformanceIsNotBegunAgain(t *testing.T) {
	s := loopServer(t)
	conf := &game.Configuration{PerformanceID: "done", Title: "Done", Composer: "Somebody", Score: partOf(t, 60, 0)}
	if err := s.db.SchedulePerformance(t.Context(), &game.Performance{Configuration: conf}); err != nil {
		t.Fatal(err)
	}
	if err := s.db.SaveRendition(t.Context(), "done", 6, []byte("what they played")); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	req := httptest.NewRequestWithContext(ctx, "POST", "/performances/done/begin", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("POST begin on a completed performance: %d %q, want 409", w.Code, w.Body.String())
	}
	select {
	case m := <-s.inbox:
		t.Errorf("the refused begin reached the loop: %s", m.Type)
	default:
	}

	s.join(t, "a")
	s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())
	if s.state.Phase != game.GamePhaseUninitialized || s.state.Configuration != nil {
		t.Errorf("the loop began a completed performance: phase %s, configuration %+v", s.state.Phase, s.state.Configuration)
	}
	if s.player("a") == nil {
		t.Error("the refused begin cleared the room")
	}
}

// partPlaying is a part shaped as the player's recorder makes one, its
// tempo first, then each key struck and released — and then whatever
// else a rewritten page might put in it.
func partPlaying(t *testing.T, keys []uint8, extra ...midi.Message) []byte {
	t.Helper()
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(120))
	for _, key := range keys {
		tr.Add(0, midi.NoteOn(recorder.Channel, key, 100))
		tr.Add(120, midi.NoteOff(recorder.Channel, key))
	}
	for _, m := range extra {
		tr.Add(10, m)
	}
	tr.Close(0)
	f := smf.New()
	if err := f.Add(tr); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := f.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// A part is the player's recording of the keys they were dealt, sent
// from a page anyone can rewrite, so the rendition takes from it the
// dealt keys alone — their strikes and releases, on the channel the
// recorder plays them on. A part used to be merged whole: any player
// could add any key to the rendition, on any channel, or anything else
// MIDI carries (a held pedal, a program change) for the disklavier to
// play. A player dealt nothing has no part to give.
func TestARenditionPlaysOnlyTheNotesDealt(t *testing.T) {
	s := loopServer(t)
	conf := &game.Configuration{PerformanceID: "tonight", Title: "Tonight", Composer: "Somebody", Score: partOf(t, 60, 0)}
	if err := s.db.SchedulePerformance(t.Context(), &game.Performance{Configuration: conf}); err != nil {
		t.Fatal(err)
	}
	s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())
	for _, f := range []string{"a", "b", "late"} {
		s.join(t, f)
	}
	s.player("a").AssignedNotes = []uint8{60}
	s.player("b").AssignedNotes = []uint8{64, 67}
	s.state.Phase = game.GamePhaseHero

	s.receive(t, game.MessageTypeSubmitPartialTrack, "a", partPlaying(t, []uint8{60, 61, 64},
		midi.ControlChange(1, 64, 127), midi.ProgramChange(1, 40),
		midi.NoteOn(9, 60, 127), midi.NoteOff(9, 60)))
	s.receive(t, game.MessageTypeSubmitPartialTrack, "b", partPlaying(t, []uint8{67}))
	s.receive(t, game.MessageTypeSubmitPartialTrack, "late", partPlaying(t, []uint8{62}))
	if _, ok := s.partials["late"]; ok || s.player("late").HasSubmitted {
		t.Error("a part was taken from a player dealt no notes")
	}
	s.receive(t, game.MessageTypeAdvancePhase, roleConductor, nil)

	saved, err := s.db.GetMIDIFile(t.Context(), "tonight")
	if err != nil {
		t.Fatal(err)
	}
	rendition, err := smf.ReadFrom(bytes.NewReader(saved))
	if err != nil {
		t.Fatalf("the saved rendition does not parse: %v", err)
	}
	played := map[uint8]bool{}
	for _, tr := range rendition.Tracks {
		for _, ev := range tr {
			var ch, key, vel uint8
			var bpm float64
			switch {
			case ev.Message.GetNoteStart(&ch, &key, &vel):
				played[key] = true
				if ch != recorder.Channel {
					t.Errorf("the rendition strikes %d on channel %d", key, ch)
				}
			case ev.Message.GetNoteEnd(&ch, &key):
				if ch != recorder.Channel {
					t.Errorf("the rendition releases %d on channel %d", key, ch)
				}
			case ev.Message.GetMetaTempo(&bpm), ev.Message.Is(smf.MetaEndOfTrackMsg):
			default:
				t.Errorf("the rendition carries a %s from a part", ev.Message)
			}
		}
	}
	if len(played) != 2 || !played[60] || !played[67] {
		t.Errorf("the rendition plays %v, want the dealt keys that were played: 60 and 67", played)
	}
}

// A part's tempo times that part's notes and nobody else's. Each part
// is read into absolute time at its own tempo before the merge, and
// the rendition is written back to ticks at one tempo of its own — so
// a player whose part declares an extreme tempo, which a rewritten
// page can, cannot shift when anyone else's keys sound. The parts'
// tempos used to be merged into the rendition, and the last one seen
// set the ticks every later note was written in: rounded to a coarse
// tick, or overflowing one, which cut seconds out of a rest.
func TestAPartsTempoMovesNoOneElsesNotes(t *testing.T) {
	// a's part: 60 struck at half a second, and again after a twelve
	// second rest, at the recorder's 120 bpm and 960 ticks a beat.
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(120))
	tr.Add(960, midi.NoteOn(recorder.Channel, 60, 100))
	tr.Add(960, midi.NoteOff(recorder.Channel, 60))
	tr.Add(24*960, midi.NoteOn(recorder.Channel, 60, 100))
	tr.Add(960, midi.NoteOff(recorder.Channel, 60))
	tr.Close(0)
	f := smf.New()
	if err := f.Add(tr); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := f.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	aPart := b.Bytes()

	// When a's notes sound in the rendition, beside b's part.
	timesBeside := func(bPart []byte) []time.Duration {
		s := loopServer(t)
		conf := &game.Configuration{PerformanceID: "tonight", Title: "Tonight", Composer: "Somebody", Score: partOf(t, 60, 0)}
		if err := s.db.SchedulePerformance(t.Context(), &game.Performance{Configuration: conf}); err != nil {
			t.Fatal(err)
		}
		s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())
		s.join(t, "a")
		s.join(t, "b")
		s.player("a").AssignedNotes = []uint8{60}
		s.player("b").AssignedNotes = []uint8{64}
		s.state.Phase = game.GamePhaseHero
		s.receive(t, game.MessageTypeSubmitPartialTrack, "a", aPart)
		s.receive(t, game.MessageTypeSubmitPartialTrack, "b", bPart)
		if len(s.partials) != 2 {
			t.Fatalf("%d parts taken, want both", len(s.partials))
		}
		s.receive(t, game.MessageTypeAdvancePhase, roleConductor, nil)
		saved, err := s.db.GetMIDIFile(t.Context(), "tonight")
		if err != nil {
			t.Fatal(err)
		}
		rendition, err := smf.ReadFrom(bytes.NewReader(saved))
		if err != nil {
			t.Fatalf("the saved rendition does not parse: %v", err)
		}
		var times []time.Duration
		for _, ev := range abstrack.FromSMF(rendition, 0).Events {
			var ch, key, vel uint8
			if (ev.Message.GetNoteStart(&ch, &key, &vel) || ev.Message.GetNoteEnd(&ch, &key)) && key == 60 {
				times = append(times, ev.Timestamp)
			}
		}
		return times
	}

	want := timesBeside(partOf(t, 64, 0))
	if len(want) != 4 {
		t.Fatalf("a's notes beside an ordinary part: %v", want)
	}
	// b's part sets the recorder's tempo, as checkPart wants, and then
	// changes it: from then on, the last tempo the merge saw set the
	// ticks every later note in the rendition was written in.
	for _, bpm := range []float64{60_000_000, 3.58} {
		var tr smf.Track
		tr.Add(0, smf.MetaTempo(120))
		tr.Add(0, midi.NoteOn(recorder.Channel, 64, 100))
		tr.Add(1, smf.MetaTempo(bpm))
		tr.Add(0, midi.NoteOff(recorder.Channel, 64))
		tr.Close(0)
		f := smf.New()
		if err := f.Add(tr); err != nil {
			t.Fatal(err)
		}
		var b bytes.Buffer
		if _, err := f.WriteTo(&b); err != nil {
			t.Fatal(err)
		}
		if got := timesBeside(b.Bytes()); !slices.Equal(got, want) {
			t.Errorf("beside a part at %v bpm, a's notes sound at %v, want %v", bpm, got, want)
		}
	}
}

// timed is a MIDI message at a time from the start of a part.
type timed struct {
	at  time.Duration
	msg midi.Message
}

// partAt is a part with its events at the times given, which must be
// in order, ending at end: at the recorder's 120 bpm and 960 ticks a
// beat, 1920 ticks a second.
func partAt(t *testing.T, end time.Duration, evs ...timed) []byte {
	t.Helper()
	ticks := func(d time.Duration) uint32 { return uint32(d * 1920 / time.Second) }
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(120))
	var now time.Duration
	for _, ev := range evs {
		tr.Add(ticks(ev.at-now), ev.msg)
		now = ev.at
	}
	tr.Close(ticks(end - now))
	f := smf.New()
	if err := f.Add(tr); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := f.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// A part is bounded in time and in density, by what the page can
// record: it records for the length of the score and a few seconds
// more, it strikes a key only once it is released, and no hand
// repeats one key faster than minStrikeGap. A rewritten page could put
// a key hours past the score (or end its part days past it, which the
// rendition's length followed), pack thousands of strikes of one key
// into one instant, or strike a key twice without releasing it. What
// the page records is taken whole — including a key still held when
// it stops recording, which is released where the part ends rather
// than left down on the disklavier.
func TestAPartIsBoundedInTimeAndDensity(t *testing.T) {
	s := loopServer(t)
	score := partOf(t, 60, 0)
	conf := &game.Configuration{PerformanceID: "tonight", Title: "Tonight", Composer: "Somebody", Score: score}
	if err := s.db.SchedulePerformance(t.Context(), &game.Performance{Configuration: conf}); err != nil {
		t.Fatal(err)
	}
	s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())
	s.join(t, "rewritten")
	s.join(t, "real")
	s.player("rewritten").AssignedNotes = []uint8{60, 62}
	s.player("real").AssignedNotes = []uint8{64}
	s.state.Phase = game.GamePhaseHero

	on := func(key uint8) midi.Message { return midi.NoteOn(recorder.Channel, key, 100) }
	off := func(key uint8) midi.Message { return midi.NoteOff(recorder.Channel, key) }

	// The rewritten page: one ordinary strike, at the loudest velocity
	// MIDI carries; a thousand strikes of the same key in one instant;
	// a key struck twice and never released; a strike two hours past
	// the score; and an end three hours past it.
	evs := []timed{{0, midi.NoteOn(recorder.Channel, 60, 127)}, {100 * time.Millisecond, off(60)}}
	for range 1000 {
		evs = append(evs, timed{time.Second, on(60)}, timed{time.Second, off(60)})
	}
	evs = append(evs,
		timed{2 * time.Second, on(62)},
		timed{2500 * time.Millisecond, on(62)},
		timed{2 * time.Hour, on(60)},
		timed{2*time.Hour + time.Second, off(60)},
	)
	s.receive(t, game.MessageTypeSubmitPartialTrack, "rewritten", partAt(t, 3*time.Hour, evs...))

	// The real page: ten strikes a tenth of a second apart, the last
	// still held when the recording stops, half a second before the
	// part ends.
	var real []timed
	for i := range 10 {
		at := time.Duration(i) * 100 * time.Millisecond
		real = append(real, timed{at, on(64)})
		if i < 9 {
			real = append(real, timed{at + 50*time.Millisecond, off(64)})
		}
	}
	s.receive(t, game.MessageTypeSubmitPartialTrack, "real", partAt(t, 1400*time.Millisecond, real...))
	if len(s.partials) != 2 {
		t.Fatalf("%d parts taken, want both", len(s.partials))
	}
	s.receive(t, game.MessageTypeAdvancePhase, roleConductor, nil)

	saved, err := s.db.GetMIDIFile(t.Context(), "tonight")
	if err != nil {
		t.Fatal(err)
	}
	rendition, err := smf.ReadFrom(bytes.NewReader(saved))
	if err != nil {
		t.Fatalf("the saved rendition does not parse: %v", err)
	}
	track := abstrack.FromSMF(rendition, 0)
	scoreFile, err := smf.ReadFrom(bytes.NewReader(score))
	if err != nil {
		t.Fatal(err)
	}
	if bound := abstrack.FromSMF(scoreFile, 0).Dur() + partSlack; track.Dur() > bound {
		t.Errorf("the rendition runs %v, past the %v a part may", track.Dur(), bound)
	}
	strikes := map[uint8]int{}
	sounding := map[uint8]bool{}
	var lastKey uint8
	for _, ev := range track.Events {
		var ch, key, vel uint8
		switch {
		case ev.Message.GetNoteStart(&ch, &key, &vel):
			if sounding[key] {
				t.Errorf("%d is struck again at %v while it sounds", key, ev.Timestamp)
			}
			if vel != recorder.Velocity {
				t.Errorf("%d is struck at %v at velocity %d, want the page's %d", key, ev.Timestamp, vel, recorder.Velocity)
			}
			strikes[key]++
			sounding[key] = true
			lastKey = key
		case ev.Message.GetNoteEnd(&ch, &key):
			sounding[key] = false
			lastKey = key
		}
	}
	if want := map[uint8]int{60: 2, 62: 1, 64: 10}; !maps.Equal(strikes, want) {
		t.Errorf("the rendition strikes %v, want %v", strikes, want)
	}
	for key, down := range sounding {
		if down {
			t.Errorf("%d is left held at the end of the rendition", key)
		}
	}
	if lastKey != 62 {
		t.Errorf("the last key the rendition moves is %d, want 62, released where the longer part ends", lastKey)
	}
}

// A rendition is what an audience played, written once. When no part
// holds a strike of a dealt key — nobody submitted one, none could be
// taken, or every player pressed nothing, which still sends a part,
// since the page submits on a timer — nothing was played, and nothing
// is written: the performance stays scheduled, so the conductor can
// begin it again, and the disklavier is sent nothing. Written, an
// empty rendition would complete the performance: off the conductor's
// picker, and refused a second beginning, with nothing in it.
func TestAPerformanceWithNoPartsWritesNoRendition(t *testing.T) {
	for name, plant := range map[string]func(s *Server){
		"no part submitted":       func(*Server) {},
		"no part could be merged": func(s *Server) { s.partials["a"] = []byte("not MIDI") },
		// What the page sends when its player pressed nothing: it
		// submits on a timer, so a quiet room still sends every part.
		"every part is empty": func(s *Server) {
			r := recorder.New(120)
			r.Close()
			bs, err := r.Bytes()
			if err != nil {
				t.Fatal(err)
			}
			s.partials["a"] = bs
		},
		"no part plays a dealt key": func(s *Server) { s.partials["a"] = partPlaying(t, []uint8{61}) },
	} {
		t.Run(name, func(t *testing.T) {
			s := loopServer(t)
			conf := &game.Configuration{PerformanceID: "quiet", Title: "Quiet", Composer: "Somebody", Score: partOf(t, 60, 0)}
			if err := s.db.SchedulePerformance(t.Context(), &game.Performance{Configuration: conf}); err != nil {
				t.Fatal(err)
			}
			s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())
			s.join(t, "a")
			s.player("a").AssignedNotes = []uint8{60}
			s.state.Phase = game.GamePhaseHero
			plant(s)
			drainOutbox(s)
			s.receive(t, game.MessageTypeAdvancePhase, roleConductor, nil)

			if s.state.Phase != game.GamePhasePlayback {
				t.Errorf("phase after assembly = %s", s.state.Phase)
			}
			if s.state.Rendition != nil {
				t.Errorf("a rendition of %d bytes is held", len(s.state.Rendition))
			}
			if complete, err := s.db.IsComplete(t.Context(), "quiet"); err != nil || complete {
				t.Errorf("the performance is complete (%v, %v), want it still scheduled", complete, err)
			}
			if saved, err := s.db.GetMIDIFile(t.Context(), "quiet"); err != nil || len(saved) != 0 {
				t.Errorf("a rendition of %d bytes was saved (%v)", len(saved), err)
			}
			for _, m := range drainOutbox(s) {
				if m.Type == game.MessageTypeSendRenditionToDisklavier {
					t.Errorf("the disklavier was sent a rendition of %d bytes", len(m.Data))
				}
			}
			s.receive(t, game.MessageTypeBeginPerformance, roleConductor, conf.Bytes())
			if s.state.Phase != game.GamePhaseLobby {
				t.Errorf("the performance could not be begun again: phase %s", s.state.Phase)
			}
		})
	}
}
