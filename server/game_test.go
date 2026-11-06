package server

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	"monks.co/piano-alone/game"
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
	tr.Add(0, midi.NoteOn(0, key, 100))
	tr.Add(120, midi.NoteOff(0, key))
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
