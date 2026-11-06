package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/recorder"
)

// What one performance can hold, whoever opens sockets. The player's
// socket is public and anonymous, so these bound the memory a stranger
// can make the server keep — at most maxPlayers parts of maxPartBytes,
// 16 MiB — and each player's part is bounded on its own, so no socket
// can crowd out another's. Both are far above anything the piece has
// played: its two audiences were twelve and six, and its largest
// rendition, six players' parts merged, is 7 KB.
const (
	// maxPlayers is how many players one performance holds.
	maxPlayers = 256
	// maxPartBytes is the largest part a player may submit.
	maxPartBytes = 64 << 10
)

// renditionBPM is the tempo a rendition is written at.
const renditionBPM = 120

// renditionTicks is the ticks a quarter note a rendition is written
// in, its header and its deltas both. It is 480 because the piece's
// first renditions counted their deltas at 480 under a header that
// said 960: db's 002 migration corrected them by rewriting the two
// header bytes rather than every delta, and at 480 the archive stays
// one resolution, in which a header saying 960 can only be one of
// those. At renditionBPM a tick is about a millisecond, finer than
// any phone times a touch.
const renditionTicks smf.MetricTicks = 480

// What a part may hold, in time and in density. The page records a
// part from the player's first key for the length of the score and
// five seconds more, ends it half a second after the last key, strikes
// a key only once it is released, and is played by a hand. A rewritten
// page is none of those: it could put its keys hours past the score
// (and the rendition runs as long as its longest part), strike one key
// thousands of times in an instant for the disklavier to play, or
// strike a key it never releases.
const (
	// partSlack is how far past the score's length a part may run. The
	// page's own parts end within six seconds of it.
	partSlack = time.Minute
	// minStrikeGap is how soon a key may be struck again. In the
	// renditions of the piece's two performances no player struck one
	// key twice within 95 ms, and no hand repeats a key twice as fast
	// as this.
	minStrikeGap = 20 * time.Millisecond
)

// play is the game loop: one goroutine, one state, messages handled
// strictly in the order they arrive, which is why nothing here takes a
// lock.
//
// A message that fails is logged and dropped. It used to return out of
// the loop, which stopped the performance — and since every kind of
// message reaches here from the public internet, that made a
// malformed frame a way to end the piece from off-site. There is no
// error a single message can produce that is worth more than the
// performance in flight.
func (s *Server) play(ctx context.Context) error {
	s.state = game.NewState()
	s.partials = map[string][]byte{}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-s.inbox:
			if err := s.handle(ctx, msg); err != nil {
				s.log.ErrorContext(ctx, "message handling failed",
					"piano.message_type", msg.Type.String(),
					"piano.origin", msg.Player,
					"error.message", err.Error(),
					"error.type", fmt.Sprintf("%T", err))
			}
		}
	}
}

// handle runs one message and turns a panic into an error.
//
// The loop is one goroutine handling bytes from the public internet,
// and the merge path walks MIDI a player submitted — abstrack's
// invariants are asserted with panics, and smf parsing is a third
// party's. A panic there would unwind the loop and take the process
// with it, mid-performance, from a room. One message failing is the
// bounded version of that, and it is logged with the type that did it.
func (s *Server) handle(ctx context.Context, msg *game.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic handling %s: %v", msg.Type, r)
		}
	}()
	return s.handleMessage(ctx, msg)
}

// player returns the state's record of a player's id, or nil. The
// tutorial messages used to index the map directly and assign through
// the result, so a player who sent one before joining took the process
// down with a nil dereference.
func (s *Server) player(id string) *game.Player {
	return s.state.Players[id]
}

func (s *Server) handleMessage(ctx context.Context, msg *game.Message) error {
	switch msg.Type {

	case game.MessageTypeRestart:
		conf := s.state.Configuration
		s.state = game.NewState()
		s.state.Configuration = conf
		s.partials = map[string][]byte{}
		s.sentSwitchToVideoModal = false
		s.broadcast(ctx, game.MessageTypeState, s.state.Bytes())
		s.log.InfoContext(ctx, "performance restarted")
		return nil

	case game.MessageTypeBeginPerformance:
		configuration, err := game.ConfigurationFromBytes(msg.Data)
		if err != nil {
			return err
		}
		// The conductor's socket can begin a configuration the HTTP
		// route never saw, so the loop asks too: a completed
		// performance's rendition is written once (db.ErrComplete), and
		// beginning it again would only lose this one's at the end. An
		// id the store cannot answer for goes ahead as it always has;
		// SaveRendition holds the line either way.
		if complete, err := s.db.IsComplete(ctx, configuration.PerformanceID); err == nil && complete {
			s.log.WarnContext(ctx, "refusing to begin a completed performance",
				"piano.performance_id", configuration.PerformanceID)
			return nil
		}
		// A performance begins with the audience in the room and nothing
		// else of the last one: its parts are not this one's to merge,
		// whoever has left is forgotten (a reload joins again), and those
		// still here start with no notes, tutorial, or part.
		s.state.Phase = game.GamePhaseUninitialized
		s.state.Configuration = configuration
		s.state.Rendition = nil
		for id, player := range s.state.Players {
			if player.ConnectionState != game.ConnectionStateConnected {
				delete(s.state.Players, id)
				continue
			}
			s.state.Players[id] = &game.Player{
				Fingerprint:     id,
				NoteCapacity:    player.NoteCapacity,
				ConnectionState: game.ConnectionStateConnected,
			}
		}
		s.partials = map[string][]byte{}
		s.sentSwitchToVideoModal = false
		s.broadcast(ctx, game.MessageTypeState, s.state.Bytes())
		s.setPhase(ctx, game.GamePhaseLobby)
		s.log.InfoContext(ctx, "performance begun",
			"piano.performance_id", configuration.PerformanceID,
			"piano.title", configuration.Title)
		return nil

	case game.MessageTypeAdvancePhase:
		return s.advance(ctx)

	case game.MessageTypeConductorConnected:
		s.state.ConductorIsConnected = true
		s.sendTo(ctx, msg.Player, game.MessageTypeState, s.state.Bytes())
		s.sendTo(ctx, toControllers, game.MessageTypeConductorConnected, nil)
		return nil

	case game.MessageTypeDisklavierConnected:
		s.state.DisklavierIsConnected = true
		s.sendTo(ctx, msg.Player, game.MessageTypeState, s.state.Bytes())
		s.sendTo(ctx, toControllers, game.MessageTypeDisklavierConnected, nil)
		return nil

	case game.MessageTypeConductorDisconnected:
		s.state.ConductorIsConnected = false
		s.sendTo(ctx, toControllers, game.MessageTypeConductorDisconnected, nil)
		return nil

	case game.MessageTypeDisklavierDisconnected:
		s.state.DisklavierIsConnected = false
		s.sendTo(ctx, toControllers, game.MessageTypeDisklavierDisconnected, nil)
		return nil

	case game.MessageTypeJoin:
		joinMessage, err := game.JoinMessageFromBytes(msg.Data)
		if err != nil {
			return err
		}
		// Who joins is the connection's to say, not the join
		// message's. A player used to be able to name any fingerprint
		// here and be registered as it.
		id := msg.Player
		if s.player(id) == nil {
			if len(s.state.Players) >= maxPlayers && !s.evictLeaver() {
				// Sent the state all the same: the page waits on it to
				// know where the performance is.
				s.log.WarnContext(ctx, "refusing a join: the performance is full",
					"piano.origin", id, "piano.players", len(s.state.Players))
				s.sendTo(ctx, id, game.MessageTypeState, s.state.Bytes())
				return nil
			}
			s.state.Players[id] = &game.Player{
				Fingerprint:  id,
				NoteCapacity: joinMessage.NoteCapacity,
			}
		}
		player := s.player(id)
		player.ConnectionState = game.ConnectionStateConnected
		s.sendTo(ctx, id, game.MessageTypeState, s.state.Bytes())
		s.sendTo(ctx, toControllers, game.MessageTypeBroadcastConnectedPlayer, player.Bytes())
		if notes := player.AssignedNotes; len(notes) > 0 {
			s.sendTo(ctx, id, game.MessageTypeAssignment, notes)
		}
		return nil

	case game.MessageTypeLeave:
		player := s.player(msg.Player)
		if player == nil {
			return nil
		}
		player.ConnectionState = game.ConnectionStateDisconnected
		s.sendTo(ctx, toControllers, game.MessageTypeBroadcastDisconnectedPlayer, []byte(msg.Player))
		return nil

	case game.MessageTypeStartTutorial:
		player := s.player(msg.Player)
		if player == nil {
			return nil
		}
		player.HasStartedTutorial = true
		s.sendTo(ctx, toControllers, game.MessageTypeStartTutorial, []byte(msg.Player))
		return nil

	case game.MessageTypeCompleteTutorial:
		player := s.player(msg.Player)
		if player == nil {
			return nil
		}
		player.HasStartedTutorial, player.HasCompletedTutorial = true, true
		s.sendTo(ctx, toControllers, game.MessageTypeCompleteTutorial, []byte(msg.Player))
		return nil

	case game.MessageTypeSubmitPartialTrack:
		// A part is a player's recording of their notes, so it is taken
		// while the performance is played and from someone in it — and
		// within the budget, which a resubmission's own last part frees.
		// Anyone may open a player's socket in any phase, and a part
		// taken whenever from whoever sat in memory until a Restart and
		// was merged into the next rendition.
		player := s.player(msg.Player)
		var refusal string
		switch {
		case s.state.Phase != game.GamePhaseHero:
			refusal = "outside Hero"
		case player == nil:
			refusal = "from a player who never joined"
		case len(player.AssignedNotes) == 0:
			// The rendition takes only a player's dealt keys from
			// their part, so there is nothing to take from this one.
			refusal = "from a player dealt no notes"
		case len(msg.Data) > maxPartBytes:
			refusal = fmt.Sprintf("larger than %d bytes", maxPartBytes)
		default:
			if err := checkPart(msg.Data); err != nil {
				refusal = err.Error()
			}
		}
		if refusal != "" {
			s.log.WarnContext(ctx, "refusing a part: "+refusal,
				"piano.origin", msg.Player,
				"piano.phase", s.state.Phase.String(),
				"piano.part_bytes", len(msg.Data))
			return nil
		}
		s.partials[msg.Player] = msg.Data
		player.HasSubmitted = true
		s.sendTo(ctx, toControllers, game.MessageTypeBroadcastSubmittedTrack, []byte(msg.Player))
		s.log.InfoContext(ctx, "part submitted",
			"piano.player", msg.Player,
			"piano.submitted", s.state.CountSubmittedTracks())
		if !s.sentSwitchToVideoModal {
			s.sentSwitchToVideoModal = true
			s.sendTo(ctx, roleDisklavier, game.MessageTypeBroadcastControllerModal, []byte("switch output to video"))
		}
		return nil

	default:
		s.log.DebugContext(ctx, "not handling message", "piano.message_type", msg.Type.String())
		return nil
	}
}

// evictLeaver makes room in a full roster by forgetting one player who
// has left holding nothing to come back to — no notes and no part —
// and reports whether it found one. Every page load joins, so a
// roster of everyone who ever opened the page fills by itself between
// performances; a leaver with notes may reload to play them, and keeps
// their place.
func (s *Server) evictLeaver() bool {
	for id, p := range s.state.Players {
		if p.ConnectionState != game.ConnectionStateConnected && len(p.AssignedNotes) == 0 && !p.HasSubmitted {
			delete(s.state.Players, id)
			return true
		}
	}
	return false
}

// readPart reads a part's bytes as MIDI. The parser is a third party's
// and panics on some inputs (an SMPTE time division among them), so a
// panic is a refusal like any other parse error.
func readPart(bs []byte) (f *smf.SMF, err error) {
	defer func() {
		if r := recover(); r != nil {
			f, err = nil, fmt.Errorf("not a MIDI file: %v", r)
		}
	}()
	f, err = smf.ReadFrom(bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("not a MIDI file: %w", err)
	}
	return f, nil
}

// checkPart refuses a part the merge could not take: abstrack reads
// the first track in metric ticks, and counts time at the tempo last
// set, so time passing before a tempo is an infinite duration, and a
// timeline long enough to overflow runs backwards — either one panics
// when the rendition is written, after every part was taken. The
// player's recorder sets its tempo first and records minutes.
func checkPart(bs []byte) error {
	f, err := readPart(bs)
	if err != nil {
		return err
	}
	if _, ok := f.TimeFormat.(smf.MetricTicks); !ok {
		return errors.New("its time format is not metric ticks")
	}
	if len(f.Tracks) == 0 {
		return errors.New("it has no tracks")
	}
	tempo := false
	for _, ev := range f.Tracks[0] {
		if ev.Delta > 0 && !tempo {
			return errors.New("time passes before its tempo is set")
		}
		var bpm float64
		if ev.Message.GetMetaTempo(&bpm) {
			if !(bpm > 0) || math.IsInf(bpm, 1) {
				return fmt.Errorf("its tempo is %v", bpm)
			}
			tempo = true
		}
	}
	var last time.Duration
	for _, ev := range abstrack.FromSMF(f, 0).Events {
		if ev.Timestamp < last {
			return errors.New("its timeline runs backwards")
		}
		last = ev.Timestamp
	}
	return nil
}

// mergePart merges what one part's bytes may contribute — the keys
// dealt, by onlyDealt, within what a page can record by end, by
// bounded — into track, and reports how many strikes that kept and how
// many of the part's events bounded left out. A panic in the merge's
// asserted invariants is an error.
func mergePart(track *abstrack.AbsTrack, bs []byte, dealt []uint8, end time.Duration) (merged *abstrack.AbsTrack, strikes, dropped int, err error) {
	defer func() {
		if r := recover(); r != nil {
			merged, strikes, dropped, err = nil, 0, 0, fmt.Errorf("merging: %v", r)
		}
	}()
	if err := checkPart(bs); err != nil {
		return nil, 0, 0, err
	}
	f, err := readPart(bs)
	if err != nil {
		return nil, 0, 0, err
	}
	part, dropped := bounded(onlyDealt(abstrack.FromSMF(f, 0), dealt), end)
	for _, ev := range part.Events {
		var channel, key, velocity uint8
		if ev.Message.GetNoteStart(&channel, &key, &velocity) {
			strikes++
		}
	}
	return track.Merge(part), strikes, dropped, nil
}

// bounded holds a part, as onlyDealt leaves it, to what the page can
// record: nothing after end; a key struck only once it is released,
// and not again within minStrikeGap of its last strike, and at
// recorder.Velocity, the one velocity the page strikes at; a key
// released only while it is down. It reports how many events it left out, which
// for a part the page recorded is none. The part ends where it said it
// did, or at end if that is sooner, and a key still down there is
// released there: the page stops recording on a timer, so a player
// holding a key when it does sends a strike with no release, and the
// disklavier would hold that key down.
func bounded(part *abstrack.AbsTrack, end time.Duration) (*abstrack.AbsTrack, int) {
	var (
		kept     = abstrack.New()
		dropped  int
		down     = map[uint8]bool{}
		struck   = map[uint8]time.Duration{}
		last     time.Duration
		partEnd  time.Duration
		hasEnded bool
	)
	for _, ev := range part.Events {
		var channel, key, velocity uint8
		switch {
		case ev.Message.Is(smf.MetaEndOfTrackMsg):
			partEnd, hasEnded = ev.Timestamp, true
			continue
		case ev.Timestamp > end:
			dropped++
			continue
		case ev.Message.GetNoteStart(&channel, &key, &velocity):
			if at, ok := struck[key]; down[key] || ok && ev.Timestamp-at < minStrikeGap {
				dropped++
				continue
			}
			down[key], struck[key] = true, ev.Timestamp
			ev.Message = smf.Message(midi.NoteOn(recorder.Channel, key, recorder.Velocity))
		case ev.Message.GetNoteEnd(&channel, &key):
			if !down[key] {
				dropped++
				continue
			}
			down[key] = false
		default:
			// onlyDealt has left only notes and the track's end.
			continue
		}
		kept.Append(ev)
		last = ev.Timestamp
	}
	if !hasEnded || partEnd < last {
		partEnd = last
	}
	partEnd = min(partEnd, end)
	for key := range 128 {
		if down[uint8(key)] {
			kept.Append(abstrack.AbsEvent{Timestamp: partEnd, Message: smf.Message(midi.NoteOff(recorder.Channel, uint8(key)))})
		}
	}
	kept.Append(abstrack.AbsEvent{Timestamp: partEnd, Message: smf.EOT})
	return kept, dropped
}

// scoreLength is how long the score's notes run, as the page reads
// them: its first track, at its own tempo.
func scoreLength(score []byte) (length time.Duration, err error) {
	defer func() {
		if r := recover(); r != nil {
			length, err = 0, fmt.Errorf("reading the score: %v", r)
		}
	}()
	f, err := readPart(score)
	if err != nil {
		return 0, err
	}
	if len(f.Tracks) == 0 {
		return 0, errors.New("the score has no tracks")
	}
	track := abstrack.FromSMF(f, 0)
	if len(track.Events) == 0 {
		return 0, nil
	}
	return track.Dur(), nil
}

// onlyDealt is what a part contributes to the rendition: the strikes
// and releases of the keys its player was dealt, on the recorder's
// channel, and the track's end — nothing else. A part is sent from a
// page anyone can rewrite, and merged whole it could play any key, on
// any channel, or send the disklavier anything else MIDI carries: a
// held pedal, a program change, a system exclusive. Its tempo goes
// too. FromSMF has already used it to put the part's notes in absolute
// time, and kept in the merge it would set the ticks every later note
// in the rendition is written in, whoever played it; assemble writes
// the rendition at a tempo of its own. The page records a tempo and
// the player's notes on recorder.Channel, so nothing a player played
// is lost.
//
// abstrack.Select is not this filter: it keeps the keys asked for and
// passes every other kind of event through, which is right for showing
// a player their score and wrong for taking what a page sends.
func onlyDealt(part *abstrack.AbsTrack, dealt []uint8) *abstrack.AbsTrack {
	keys := map[uint8]bool{}
	for _, key := range dealt {
		keys[key] = true
	}
	kept := abstrack.New()
	for _, ev := range part.Events {
		var channel, key, velocity uint8
		switch {
		case ev.Message.GetNoteStart(&channel, &key, &velocity), ev.Message.GetNoteEnd(&channel, &key):
			if channel == recorder.Channel && keys[key] {
				kept.Append(ev)
			}
		case ev.Message.Is(smf.MetaEndOfTrackMsg):
			kept.Append(ev)
		}
	}
	return kept
}

// advance moves the performance on one phase. Which phase it lands in
// depends only on where it was: the conductor presses a button, and
// nothing here is on a timer or waiting for a quorum.
func (s *Server) advance(ctx context.Context) error {
	switch s.state.Phase {
	case game.GamePhaseLobby:
		return s.assignNotes(ctx)

	case game.GamePhaseHero:
		return s.assemble(ctx)

	case game.GamePhasePlayback:
		s.setPhase(ctx, game.GamePhaseDone)
		return nil

	default:
		return nil
	}
}

// assignNotes deals the score's keys out to the connected players,
// most-played key first, round-robin, skipping anyone already at their
// capacity — three notes on a touch device, one otherwise. Everyone
// gets a different handful of the same song, which is the piece.
func (s *Server) assignNotes(ctx context.Context) error {
	if s.state.Configuration == nil {
		return errors.New("advancing out of the lobby with no performance configured")
	}
	var ids []string
	for id := range s.state.Players {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		s.log.InfoContext(ctx, "not advancing: nobody has joined")
		return nil
	}

	score, err := smf.ReadFrom(bytes.NewReader(s.state.Configuration.Score))
	if err != nil {
		return fmt.Errorf("parsing score: %w", err)
	}
	notes := abstrack.FromSMF(score, 0).CountNotes()
	for i, note := range notes {
		player := s.state.Players[ids[i%len(ids)]]
		if len(player.AssignedNotes) < player.NoteCapacity {
			player.AssignedNotes = append(player.AssignedNotes, note.Key)
		}
	}
	for _, player := range s.state.Players {
		s.sendTo(ctx, player.Fingerprint, game.MessageTypeAssignment, player.AssignedNotes)
	}
	s.setPhase(ctx, game.GamePhaseHero)
	s.log.InfoContext(ctx, "notes assigned",
		"piano.players", len(ids),
		"piano.distinct_notes", len(notes))
	return nil
}

// assemble merges every submitted part into one MIDI file, saves it,
// and sends it to the disklavier. This is the piece's one irreplaceable
// write: the rendition is what an audience played, and it exists
// nowhere else.
func (s *Server) assemble(ctx context.Context) error {
	s.setPhase(ctx, game.GamePhaseProcessing)

	// The rendition's one tempo, which ToSMF writes every part's notes
	// back into ticks at. The parts are merged in absolute time, so any
	// tempo would do; this is MIDI's default, and the one the page
	// records at.
	track := abstrack.New()
	track.Append(abstrack.AbsEvent{Message: smf.MetaTempo(renditionBPM)})
	// How long a part may run. The score was read when the notes were
	// dealt, so it reads now; if it somehow does not, the parts are held
	// to the slack alone rather than to nothing.
	end := partSlack
	if s.state.Configuration != nil {
		length, err := scoreLength(s.state.Configuration.Score)
		if err != nil {
			s.log.ErrorContext(ctx, "holding the parts to the slack alone",
				"error.message", err.Error(),
				"error.type", fmt.Sprintf("%T", err))
		}
		end += length
	}
	strikes := 0
	for id, bs := range s.partials {
		// Read again here rather than kept parsed since it was taken: a
		// parsed SMF is several times its bytes. checkPart passed it
		// when it was taken, so this cannot fail; if it somehow does,
		// one part is not worth the rendition. A part is taken only
		// from a player in the state, and Begin drops the parts with
		// the players, so every part has its player here.
		var dealt []uint8
		if player := s.player(id); player != nil {
			dealt = player.AssignedNotes
		}
		withPart, struck, dropped, err := mergePart(track, bs, dealt, end)
		if err != nil {
			s.log.ErrorContext(ctx, "leaving a part out of the rendition",
				"piano.player", id,
				"error.message", err.Error(),
				"error.type", fmt.Sprintf("%T", err))
			continue
		}
		if dropped > 0 {
			s.log.WarnContext(ctx, "leaving out what no page records",
				"piano.player", id,
				"piano.dropped_events", dropped)
		}
		track = withPart
		strikes += struck
	}
	if strikes == 0 {
		// Nothing was played — no part was submitted, none could be
		// merged, or every one merged holds no strike of a dealt key,
		// which is what the page sends when its player pressed nothing,
		// since it submits on a timer — so there is nothing to keep, and
		// a rendition is written once: an empty one would complete the
		// performance, taking it off the conductor's picker and refusing
		// it a second beginning, with nothing in it. It stays scheduled,
		// and the disklavier is sent nothing.
		s.state.Rendition = nil
		s.log.WarnContext(ctx, "nothing played: writing no rendition",
			"piano.parts", len(s.partials))
		s.setPhase(ctx, game.GamePhasePlayback)
		return nil
	}
	file := track.ToSMF(renditionTicks)
	var buf bytes.Buffer
	if _, err := file.WriteTo(&buf); err != nil {
		return fmt.Errorf("writing rendition: %w", err)
	}
	bs := buf.Bytes()
	s.state.Rendition = bs

	submitted := s.state.CountSubmittedTracks()
	if s.state.Configuration != nil {
		if err := s.db.SaveRendition(ctx, s.state.Configuration.PerformanceID, submitted, bs); err != nil {
			// The performance goes on: the disklavier still gets the
			// file below, and the rendition is in memory. Losing the
			// recording is worth an error a human will see, not a
			// stopped playback in a room full of people.
			s.log.ErrorContext(ctx, "saving rendition failed",
				"piano.performance_id", s.state.Configuration.PerformanceID,
				"error.message", err.Error(),
				"error.type", fmt.Sprintf("%T", err))
		}
	}

	s.sendTo(ctx, roleDisklavier, game.MessageTypeSendRenditionToDisklavier, bs)
	s.setPhase(ctx, game.GamePhasePlayback)
	s.log.InfoContext(ctx, "rendition assembled",
		"piano.parts", submitted,
		"piano.rendition_bytes", len(bs))
	return nil
}

func (s *Server) setPhase(ctx context.Context, phase game.Phase) {
	s.state.Phase = phase
	s.broadcast(ctx, game.MessageTypeBroadcastPhase, phase.Bytes())
	s.log.InfoContext(ctx, "phase changed", "piano.phase", phase.String())
}

func (s *Server) broadcast(ctx context.Context, msgType game.MessageType, payload []byte) {
	s.out(ctx, game.NewMessage(msgType, toEveryone, payload))
}

func (s *Server) sendTo(ctx context.Context, to string, msgType game.MessageType, payload []byte) {
	s.out(ctx, game.NewMessage(msgType, to, payload))
}

// out hands a message to the delivery goroutine. It is only ever
// called from the game loop, so a blocked send would deadlock the
// performance; the ctx case is what keeps a shutdown from wedging.
func (s *Server) out(ctx context.Context, m *game.Message) {
	select {
	case s.outbox <- m:
	case <-ctx.Done():
	}
}
