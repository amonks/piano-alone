package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/db"
	"monks.co/piano-alone/game"
)

type GameServer struct {
	db *db.DB

	send chan<- *game.Message
	bus  chan *game.Message

	state    *game.State
	partials map[string]*smf.SMF

	sentSwitchToVideoModal bool
}

func NewGame(db *db.DB) *GameServer {
	return &GameServer{
		db:       db,
		state:    game.NewState(),
		partials: map[string]*smf.SMF{},
	}
}

func (gs *GameServer) Start(ctx context.Context, send chan<- *game.Message, recv <-chan *game.Message, songTitle, songComposer string, song *smf.SMF) error {
	gs.send = send
	gs.bus = make(chan *game.Message)

	msgs := make(chan *game.Message)
	go func() {
		var m *game.Message
		var ok bool
		for {
			select {
			case m, ok = <-gs.bus:
			case m, ok = <-recv:
			}
			if !ok {
				return
			}
			log.Printf("handle: '%s'", m.Type.String())
			msgs <- m
		}
	}()

	done := ctx.Done()
	for {
		select {
		case <-done:
			return nil

		case msg := <-msgs:
			if err := gs.handleMessage(msg); err != nil {
				log.Printf("msg '%s' produced error: %s", msg.Type, err)
				return err
			}
		}
	}
}

func (gs *GameServer) handleMessage(msg *game.Message) error {
	switch msg.Type {

	case game.MessageTypeRestart:
		conf := gs.state.Configuration
		gs.state = game.NewState()
		gs.state.Configuration = conf
		gs.partials = map[string]*smf.SMF{}
		gs.broadcast(game.MessageTypeState, gs.state.Bytes())
		return nil

	case game.MessageTypeBeginPerformance:
		configuration := game.ConfigurationFromBytes(msg.Data)
		gs.state.Phase = game.GamePhaseUninitialized
		gs.state.Configuration = configuration
		gs.broadcast(game.MessageTypeState, gs.state.Bytes())
		gs.setPhase(game.GamePhaseLobby)
		return nil

	case game.MessageTypeAdvancePhase:
		switch gs.state.Phase {
		case game.GamePhaseLobby:
			var fingerprints []string
			for f := range gs.state.Players {
				fingerprints = append(fingerprints, f)
			}
			if len(fingerprints) == 0 {
				return nil
			}
			r := bytes.NewReader(gs.state.Configuration.Score)
			f, err := smf.ReadFrom(r)
			if err != nil {
				return err
			}
			track := abstrack.FromSMF(f, 0)
			notes := track.CountNotes()
			for i, note := range notes {
				fingerprint := fingerprints[i%len(fingerprints)]
				player := gs.state.Players[fingerprint]
				if len(player.AssignedNotes) < player.NoteCapacity {
					player.AssignedNotes = append(player.AssignedNotes, note.Key)
				}
			}
			for _, player := range gs.state.Players {
				gs.sendTo(player.Fingerprint, game.MessageTypeAssignment, player.AssignedNotes)
			}
			gs.setPhase(game.GamePhaseHero)
			return nil

		case game.GamePhaseHero:
			gs.setPhase(game.GamePhaseProcessing)
			track := abstrack.New()
			for _, partial := range gs.partials {
				if partial == nil {
					continue
				}
				track = track.Merge(abstrack.FromSMF(partial, 0))
			}
			file := smf.New()
			file.Add(track.ToSMF())
			var buf bytes.Buffer
			if _, err := file.WriteTo(&buf); err != nil {
				return err
			}
			bs := buf.Bytes()
			gs.state.Rendition = bs
			if err := gs.db.SaveRendition(gs.state.Configuration.PerformanceID, gs.state.CountSubmittedTracks(), bs); err != nil {
				return err
			}
			gs.sendTo("disklavier", game.MessageTypeSendRenditionToDisklavier, bs)
			gs.setPhase(game.GamePhasePlayback)
			return nil

		case game.GamePhasePlayback:
			gs.setPhase(game.GamePhaseDone)
			return nil

		default:
			return nil
		}

	case game.MessageTypeConductorConnected:
		gs.state.ConductorIsConnected = true
		gs.sendTo(msg.Player, game.MessageTypeState, gs.state.Bytes())
		gs.sendTo("controllers", game.MessageTypeConductorConnected, nil)
		return nil

	case game.MessageTypeDisklavierConnected:
		gs.state.DisklavierIsConnected = true
		gs.sendTo(msg.Player, game.MessageTypeState, gs.state.Bytes())
		gs.sendTo("controllers", game.MessageTypeDisklavierConnected, nil)
		return nil

	case game.MessageTypeConductorDisconnected:
		gs.state.ConductorIsConnected = false
		gs.sendTo("controllers", game.MessageTypeConductorDisconnected, nil)
		return nil

	case game.MessageTypeDisklavierDisconnected:
		gs.state.DisklavierIsConnected = false
		gs.sendTo("controllers", game.MessageTypeDisklavierDisconnected, nil)
		return nil

	case game.MessageTypeJoin:
		joinMessage := game.JoinMessageFromBytes(msg.Data)
		fingerprint, noteCapacity := joinMessage.Fingerprint, joinMessage.NoteCapacity
		if _, got := gs.state.Players[fingerprint]; !got {
			gs.state.Players[fingerprint] = &game.Player{
				ConnectionState: game.ConnectionStateConnected,
				Fingerprint:     msg.Player,
				NoteCapacity:    noteCapacity,
			}
		}
		gs.state.Players[fingerprint].ConnectionState = game.ConnectionStateConnected
		gs.sendTo(msg.Player, game.MessageTypeState, gs.state.Bytes())
		gs.sendTo("controllers", game.MessageTypeBroadcastConnectedPlayer, gs.state.Players[fingerprint].Bytes())
		if notes := gs.state.Players[fingerprint].AssignedNotes; len(notes) > 0 {
			gs.sendTo(fingerprint, game.MessageTypeAssignment, notes)
		}
		return nil

	case game.MessageTypeLeave:
		if _, has := gs.state.Players[msg.Player]; !has {
			return nil
		}
		gs.state.Players[msg.Player].ConnectionState = game.ConnectionStateDisconnected
		gs.sendTo("controllers", game.MessageTypeBroadcastDisconnectedPlayer, []byte(msg.Player))
		return nil

	case game.MessageTypeStartTutorial:
		gs.state.Players[msg.Player].HasStartedTutorial = true
		gs.sendTo("controllers", game.MessageTypeStartTutorial, []byte(msg.Player))
		return nil
	case game.MessageTypeCompleteTutorial:
		gs.state.Players[msg.Player].HasStartedTutorial = true
		gs.state.Players[msg.Player].HasCompletedTutorial = true
		gs.sendTo("controllers", game.MessageTypeCompleteTutorial, []byte(msg.Player))
		return nil

	case game.MessageTypeSubmitPartialTrack:
		smf, err := smf.ReadFrom(bytes.NewReader(msg.Data))
		if err != nil {
			return fmt.Errorf("smf parsing error: '%s'", err)
		}
		gs.partials[msg.Player] = smf
		if _, ok := gs.state.Players[msg.Player]; !ok {
			gs.state.Players[msg.Player] = &game.Player{}
		}
		gs.state.Players[msg.Player].HasSubmitted = true
		gs.sendTo("controllers", game.MessageTypeBroadcastSubmittedTrack, []byte(msg.Player))
		if !gs.sentSwitchToVideoModal {
			gs.sentSwitchToVideoModal = true
			gs.sendTo("disklavier", game.MessageTypeBroadcastControllerModal, []byte("switch output to video"))
		}
		return nil

	default:
		log.Printf("not handling message (type: '%s')", msg.Type.String())
		return nil
	}
}

func (gs *GameServer) setPhase(phase game.Phase) {
	gs.state.Phase = phase
	gs.broadcast(game.MessageTypeBroadcastPhase, gs.state.Phase.Bytes())
}

func (gs *GameServer) broadcast(msgType game.MessageType, data []byte) {
	gs.send <- game.NewMessage(msgType, "*", data)
}

func (gs *GameServer) sendTo(fingerprint string, msgType game.MessageType, data []byte) {
	gs.send <- game.NewMessage(msgType, fingerprint, data)
}
