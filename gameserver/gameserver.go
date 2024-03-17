package gameserver

import (
	"bytes"
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/pianists"
)

const (
	maxNotesPerPlayer = 3
)

type GameServer struct {
	send chan<- *game.Message
	bus  chan *game.Message

	state    *game.State
	partials map[string]*smf.SMF
}

func New() *GameServer {
	return &GameServer{
		state:    game.NewState(),
		partials: map[string]*smf.SMF{},
	}
}

func (gs *GameServer) Start(send chan<- *game.Message, recv <-chan *game.Message, song *smf.SMF) {
	gs.state.Score = song
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

	for {
		msg := <-msgs
		if err := gs.handleMessage(msg); err != nil {
			log.Printf("msg '%s' produced error: %s", msg.Type, err)
		}
	}
}

func (gs *GameServer) handleMessage(msg *game.Message) error {
	switch msg.Type {

	case game.MessageTypeRestart:
		score := gs.state.Score
		gs.state = game.NewState()
		gs.state.Score = score
		gs.partials = map[string]*smf.SMF{}
		gs.broadcast(game.MessageTypeInitialState, gs.state.Bytes())
		return nil

	case game.MessageTypeAdvancePhase:
		switch gs.state.Phase.Type {
		case game.GamePhaseLobby:
			var fingerprints []string
			for f := range gs.state.Players {
				fingerprints = append(fingerprints, f)
			}
			track := abstrack.FromSMF(gs.state.Score, 0)
			notes := track.CountNotes()
			for i, note := range notes {
				player := fingerprints[i%len(fingerprints)]
				if len(gs.state.Players[player].Notes) >= maxNotesPerPlayer {
					log.Printf("completed assignment with %d unassigned notes", len(notes)-i)
					break
				}
				gs.state.Players[player].Notes = append(gs.state.Players[player].Notes, note.Key)
			}
			for _, player := range gs.state.Players {
				gs.sendTo(player.Fingerprint, game.MessageTypeAssignment, player.Notes)
			}
			gs.setPhase(game.GamePhaseHero)
			gs.sendTo("controller", game.MessageTypeBroadcastControllerModal, []byte("switch output to video"))
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
			gs.state.Rendition = file
			var buf bytes.Buffer
			if _, err := gs.state.Rendition.WriteTo(&buf); err != nil {
				return err
			}
			gs.broadcast(game.MessageTypeBroadcastCombinedTrack, buf.Bytes())
			gs.setPhase(game.GamePhasePlayback)
			return nil

		case game.GamePhasePlayback:
			gs.setPhase(game.GamePhaseDone)
		}
		return nil

	case game.MessageTypeControllerJoin:
		gs.sendTo("controller", game.MessageTypeInitialState, gs.state.Bytes())
		return nil

	case game.MessageTypeJoin:
		if _, got := gs.state.Players[msg.Player]; !got {
			gs.state.Players[msg.Player] = &game.Player{
				ConnectionState: game.ConnectionStateConnected,
				Fingerprint:     msg.Player,
				Pianist:         pianists.Hash(msg.Player),
			}
		}
		gs.state.Players[msg.Player].ConnectionState = game.ConnectionStateConnected
		gs.sendTo(msg.Player, game.MessageTypeInitialState, gs.state.Bytes())
		gs.broadcast(game.MessageTypeBroadcastConnectedPlayer, gs.state.Players[msg.Player].Bytes())
		if notes := gs.state.Players[msg.Player].Notes; len(notes) > 0 {
			gs.sendTo(msg.Player, game.MessageTypeAssignment, notes)
		}
		return nil

	case game.MessageTypeLeave:
		if _, has := gs.state.Players[msg.Player]; !has {
			return nil
		}
		gs.state.Players[msg.Player].ConnectionState = game.ConnectionStateDisconnected
		gs.broadcast(game.MessageTypeBroadcastDisconnectedPlayer, []byte(msg.Player))
		if gs.state.CountConnectedPlayers() == 0 {
			score := gs.state.Score
			gs.state = game.NewState()
			gs.state.Score = score
			gs.partials = map[string]*smf.SMF{}
			gs.broadcast(game.MessageTypeInitialState, gs.state.Bytes())
			return nil
		}
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
		gs.broadcast(game.MessageTypeBroadcastSubmittedTrack, []byte(msg.Player))
		return nil

	default:
		log.Printf("not handling message (type: '%s')", msg.Type.String())
		return nil
	}
}

func (gs *GameServer) setPhase(phase game.GamePhase) {
	log.Println("phase:", phase)
	gs.state.Phase = game.NewPhase(phase)
	gs.broadcast(game.MessageTypeBroadcastPhase, gs.state.Phase.Bytes())
}

func (gs *GameServer) broadcast(msgType game.MessageType, data []byte) {
	gs.send <- game.NewMessage(msgType, "*", data)
}

func (gs *GameServer) sendTo(fingerprint string, msgType game.MessageType, data []byte) {
	gs.send <- game.NewMessage(msgType, fingerprint, data)
}
