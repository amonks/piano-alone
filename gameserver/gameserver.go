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
	lobbyDur    = time.Second * 20
	recordDur   = time.Second * 5
	playbackDur = time.Second * 5
)

type GameServer struct {
	send chan<- *game.Message
	recv <-chan *game.Message

	state    *game.State
	partials map[string]*smf.SMF
}

func New() *GameServer {
	return &GameServer{
		state: game.NewState(),
	}
}

func (gs *GameServer) Start(send chan<- *game.Message, recv <-chan *game.Message, song *smf.SMF) error {
	gs.state.Score = song
	gs.send = send
	gs.recv = recv

	// lobby phase
	done := gs.setPhase(game.GamePhaseLobby, lobbyDur)
lobby:
	select {
	case msg := <-recv:
		switch msg.Type {
		case game.MessageTypeJoin:
			if _, got := gs.state.Players[msg.Player]; !got {
				gs.addPlayer(msg.Player)
				gs.sendTo(msg.Player, game.MessageTypeInitialState, gs.state.Bytes())
				gs.broadcast(game.MessageTypeBroadcastConnectedPlayer, gs.state.Players[msg.Player].Bytes())
			} else {
				gs.state.Players[msg.Player].ConnectionState = game.ConnectionStateConnected
			}
		case game.MessageTypeLeave:
			gs.state.Players[msg.Player].ConnectionState = game.ConnectionStateDisconnected
			gs.broadcast(game.MessageTypeBroadcastDisconnectedPlayer, []byte(msg.Player))
		default:
			log.Printf("unhandled message type '%s'", msg.Type.String())
		}
		goto lobby
	case <-done:
		fmt.Println("DONE")
	}

	// split tracks
	gs.setPhase(game.GamePhaseSplitting, 0)
	if err := gs.splitTracksForPlayers(); err != nil {
		return err
	}

	// hero phase
	done = gs.setPhase(game.GamePhaseHero, recordDur)
hero:
	select {
	case msg := <-recv:
		switch msg.Type {
		case game.MessageTypeSubmitPartialTrack:
			smf, err := smf.ReadFrom(bytes.NewReader(msg.Data))
			if err != nil {
				log.Printf("smf parsing error: '%s'", msg.Type)
			}
			gs.partials[msg.Player] = smf
		default:
			log.Printf("unhandled message type '%s'", msg.Type)
		}
		goto hero
	case <-done:
	}

	// combine tracks
	gs.setPhase(game.GamePhaseJoining, 0)
	gs.combinePlayerTracks()
	var buf bytes.Buffer
	if _, err := gs.state.Rendition.WriteTo(&buf); err != nil {
		return err
	}
	bs := buf.Bytes()

	// playback phase
	done = gs.setPhase(game.GamePhasePlayback, playbackDur)
	gs.broadcast(game.MessageTypeBroadcastCombinedTrack, bs)
	<-done

	// done
	gs.setPhase(game.GamePhaseDone, 0)
	return nil
}

func (gs *GameServer) setPhase(phase game.GamePhase, dur time.Duration) <-chan time.Time {
	log.Println("phase:", phase)
	if dur == 0 {
		gs.state.Phase = phase
		gs.state.PhaseExp = time.Time{}
		msg := &game.PhaseChangeMessage{Phase: phase}
		gs.broadcast(game.MessageTypeBroadcastPhase, msg.Bytes())
		return nil
	}
	gs.state.Phase = phase
	gs.state.PhaseExp = time.Now().Add(dur)
	msg := &game.PhaseChangeMessage{Phase: phase, Exp: gs.state.PhaseExp}
	gs.broadcast(game.MessageTypeBroadcastPhase, msg.Bytes())
	return time.After(time.Until(gs.state.PhaseExp))
}

func (gs *GameServer) addPlayer(fingerprint string) {
	gs.state.Players[fingerprint] = &game.Player{
		ConnectionState: game.ConnectionStateConnected,
		Fingerprint:     fingerprint,
		Pianist:         pianists.Hash(fingerprint),
	}
}

func (gs *GameServer) splitTracksForPlayers() error {
	var fingerprints []string
	for f := range gs.state.Players {
		fingerprints = append(fingerprints, f)
	}
	track := abstrack.FromSMF(gs.state.Score.Tracks[0])
	notes := track.CountNotes()
	for i, note := range notes {
		player := fingerprints[i%len(fingerprints)]
		gs.state.Players[player].Notes = append(gs.state.Players[player].Notes, note.Key)
	}
	for _, player := range gs.state.Players {
		split := smf.New()
		split.Add(track.Select(player.Notes).ToSMF())
		var buf bytes.Buffer
		if _, err := split.WriteTo(&buf); err != nil {
			return err
		}
		gs.sendTo(player.Fingerprint, game.MessageTypeAssignment, buf.Bytes())
	}
	return nil
}

func (gs *GameServer) combinePlayerTracks() {
	track := abstrack.New()
	for _, partial := range gs.partials {
		if partial == nil {
			continue
		}
		track = track.Merge(abstrack.FromSMF(partial.Tracks[0]))
	}
	file := smf.New()
	file.Add(track.ToSMF())
	gs.state.Rendition = file
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
