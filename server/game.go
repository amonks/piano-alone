package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"gitlab.com/gomidi/midi/v2/smf"

	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
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
	s.partials = map[string]*smf.SMF{}

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

// player returns the state's record of a fingerprint, or nil. The
// tutorial messages used to index the map directly and assign through
// the result, so a player who sent one before joining took the process
// down with a nil dereference.
func (s *Server) player(fingerprint string) *game.Player {
	return s.state.Players[fingerprint]
}

func (s *Server) handleMessage(ctx context.Context, msg *game.Message) error {
	switch msg.Type {

	case game.MessageTypeRestart:
		conf := s.state.Configuration
		s.state = game.NewState()
		s.state.Configuration = conf
		s.partials = map[string]*smf.SMF{}
		s.sentSwitchToVideoModal = false
		s.broadcast(ctx, game.MessageTypeState, s.state.Bytes())
		s.log.InfoContext(ctx, "performance restarted")
		return nil

	case game.MessageTypeBeginPerformance:
		configuration, err := game.ConfigurationFromBytes(msg.Data)
		if err != nil {
			return err
		}
		s.state.Phase = game.GamePhaseUninitialized
		s.state.Configuration = configuration
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
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeConductorConnected, nil)
		return nil

	case game.MessageTypeDisklavierConnected:
		s.state.DisklavierIsConnected = true
		s.sendTo(ctx, msg.Player, game.MessageTypeState, s.state.Bytes())
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeDisklavierConnected, nil)
		return nil

	case game.MessageTypeConductorDisconnected:
		s.state.ConductorIsConnected = false
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeConductorDisconnected, nil)
		return nil

	case game.MessageTypeDisklavierDisconnected:
		s.state.DisklavierIsConnected = false
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeDisklavierDisconnected, nil)
		return nil

	case game.MessageTypeJoin:
		joinMessage, err := game.JoinMessageFromBytes(msg.Data)
		if err != nil {
			return err
		}
		// The fingerprint is the connection's, not the join message's.
		// A player used to be able to name any fingerprint here and be
		// registered as it.
		fingerprint := msg.Player
		if s.player(fingerprint) == nil {
			s.state.Players[fingerprint] = &game.Player{
				Fingerprint:  fingerprint,
				NoteCapacity: joinMessage.NoteCapacity,
			}
		}
		player := s.player(fingerprint)
		player.ConnectionState = game.ConnectionStateConnected
		s.sendTo(ctx, fingerprint, game.MessageTypeState, s.state.Bytes())
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeBroadcastConnectedPlayer, player.Bytes())
		if notes := player.AssignedNotes; len(notes) > 0 {
			s.sendTo(ctx, fingerprint, game.MessageTypeAssignment, notes)
		}
		return nil

	case game.MessageTypeLeave:
		player := s.player(msg.Player)
		if player == nil {
			return nil
		}
		player.ConnectionState = game.ConnectionStateDisconnected
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeBroadcastDisconnectedPlayer, []byte(msg.Player))
		return nil

	case game.MessageTypeStartTutorial:
		player := s.player(msg.Player)
		if player == nil {
			return nil
		}
		player.HasStartedTutorial = true
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeStartTutorial, []byte(msg.Player))
		return nil

	case game.MessageTypeCompleteTutorial:
		player := s.player(msg.Player)
		if player == nil {
			return nil
		}
		player.HasStartedTutorial, player.HasCompletedTutorial = true, true
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeCompleteTutorial, []byte(msg.Player))
		return nil

	case game.MessageTypeSubmitPartialTrack:
		partial, err := smf.ReadFrom(bytes.NewReader(msg.Data))
		if err != nil {
			return fmt.Errorf("parsing submitted track: %w", err)
		}
		s.partials[msg.Player] = partial
		if s.player(msg.Player) == nil {
			s.state.Players[msg.Player] = &game.Player{Fingerprint: msg.Player}
		}
		s.player(msg.Player).HasSubmitted = true
		s.sendTo(ctx, fingerprintControllers, game.MessageTypeBroadcastSubmittedTrack, []byte(msg.Player))
		s.log.InfoContext(ctx, "part submitted",
			"piano.fingerprint", msg.Player,
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
	var fingerprints []string
	for f := range s.state.Players {
		fingerprints = append(fingerprints, f)
	}
	if len(fingerprints) == 0 {
		s.log.InfoContext(ctx, "not advancing: nobody has joined")
		return nil
	}

	score, err := smf.ReadFrom(bytes.NewReader(s.state.Configuration.Score))
	if err != nil {
		return fmt.Errorf("parsing score: %w", err)
	}
	notes := abstrack.FromSMF(score, 0).CountNotes()
	for i, note := range notes {
		player := s.state.Players[fingerprints[i%len(fingerprints)]]
		if len(player.AssignedNotes) < player.NoteCapacity {
			player.AssignedNotes = append(player.AssignedNotes, note.Key)
		}
	}
	for _, player := range s.state.Players {
		s.sendTo(ctx, player.Fingerprint, game.MessageTypeAssignment, player.AssignedNotes)
	}
	s.setPhase(ctx, game.GamePhaseHero)
	s.log.InfoContext(ctx, "notes assigned",
		"piano.players", len(fingerprints),
		"piano.distinct_notes", len(notes))
	return nil
}

// assemble merges every submitted part into one MIDI file, saves it,
// and sends it to the disklavier. This is the piece's one irreplaceable
// write: the rendition is what an audience played, and it exists
// nowhere else.
func (s *Server) assemble(ctx context.Context) error {
	s.setPhase(ctx, game.GamePhaseProcessing)

	track := abstrack.New()
	for _, partial := range s.partials {
		if partial == nil {
			continue
		}
		track = track.Merge(abstrack.FromSMF(partial, 0))
	}
	file := smf.New()
	file.Add(track.ToSMF())
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
	s.out(ctx, game.NewMessage(msgType, fingerprintEveryone, payload))
}

func (s *Server) sendTo(ctx context.Context, fingerprint string, msgType game.MessageType, payload []byte) {
	s.out(ctx, game.NewMessage(msgType, fingerprint, payload))
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
