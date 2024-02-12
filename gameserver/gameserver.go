package gameserver

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/pianists"
)

const (
	lobbyDur    = time.Second * 5
	recordDur   = time.Minute*4 + time.Second*15
	playbackDur = time.Second * 5
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
	case game.MessageTypeJoin:
		shouldStart := gs.state.Players == nil || len(gs.state.Players) == 0
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
		if shouldStart {
			end := gs.setPhase(game.GamePhaseLobby, lobbyDur)
			gs.after(end, game.MessageTypeExpireLobby)
		}
		if notes := gs.state.Players[msg.Player].Notes; len(notes) > 0 {
			gs.sendTo(msg.Player, game.MessageTypeAssignment, notes)
		}
	case game.MessageTypeLeave:
		gs.state.Players[msg.Player].ConnectionState = game.ConnectionStateDisconnected
		gs.broadcast(game.MessageTypeBroadcastDisconnectedPlayer, []byte(msg.Player))
	case game.MessageTypeExpireLobby:
		var fingerprints []string
		for f := range gs.state.Players {
			fingerprints = append(fingerprints, f)
		}
		track := abstrack.FromSMF(gs.state.Score, 0)
		notes := track.CountNotes()
		for i, note := range notes {
			player := fingerprints[i%len(fingerprints)]
			if len(gs.state.Players[player].Notes) >= 100 {
				log.Printf("completed assignment with %d unassigned notes", len(notes)-i)
				break
			}
			gs.state.Players[player].Notes = append(gs.state.Players[player].Notes, note.Key)
		}
		for _, player := range gs.state.Players {
			gs.sendTo(player.Fingerprint, game.MessageTypeAssignment, player.Notes)
		}
		end := gs.setPhase(game.GamePhaseHero, recordDur)
		gs.after(end, game.MessageTypeExpireHero)
	case game.MessageTypeSubmitPartialTrack:
		smf, err := smf.ReadFrom(bytes.NewReader(msg.Data))
		if err != nil {
			return fmt.Errorf("smf parsing error: '%s'", err)
		}
		gs.partials[msg.Player] = smf
	case game.MessageTypeExpireHero:
		gs.setPhase(game.GamePhaseProcessing, 0)
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
		end := gs.setPhase(game.GamePhasePlayback, playbackDur)
		gs.after(end, game.MessageTypeExpirePlayback)
	case game.MessageTypeExpirePlayback:
		gs.setPhase(game.GamePhaseDone, 0)
	default:
		log.Printf("not handling message (type: '%s')", msg.Type.String())
	}
	return nil
}

func (gs *GameServer) after(delay <-chan time.Time, msgType game.MessageType) {
	go func() {
		<-delay
		gs.bus <- &game.Message{Type: msgType}
	}()
}

func (gs *GameServer) setPhase(phase game.GamePhase, dur time.Duration) <-chan time.Time {
	log.Println("phase:", phase)
	if dur == 0 {
		gs.state.Phase = game.NewPhase(phase)
		gs.broadcast(game.MessageTypeBroadcastPhase, gs.state.Phase.Bytes())
		return nil
	}
	gs.state.Phase = game.NewPhase(phase).WithExp(time.Now().Add(dur))
	gs.broadcast(game.MessageTypeBroadcastPhase, gs.state.Phase.Bytes())
	return time.After(time.Until(gs.state.Phase.Exp))
}

func (gs *GameServer) broadcast(msgType game.MessageType, data []byte) {
	gs.send <- &game.Message{
		Type:   msgType,
		Player: "*",
		Data:   data,
	}
}

func (gs *GameServer) sendTo(fingerprint string, msgType game.MessageType, data []byte) {
	gs.send <- &game.Message{
		Type:   msgType,
		Player: fingerprint,
		Data:   data,
	}
}
