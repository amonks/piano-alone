package gameserver

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/pianists"
	"monks.co/piano-alone/proto"
)

const (
	lobbyDur    = time.Second * 10
	recordDur   = time.Second * 5
	playbackDur = time.Second * 5
)

type player struct {
	state       string
	fingerprint string
	pianist     string
	notes       []uint8
	smf         *smf.SMF
}

type GameServer struct {
	phase   proto.GamePhase
	exp     time.Time
	send    chan<- *proto.Message
	recv    <-chan *proto.Message
	players map[string]*player
	song    *smf.SMF
}

func New() *GameServer {
	return &GameServer{
		players: map[string]*player{},
		song:    nil,
	}
}

func (gs *GameServer) Start(send chan<- *proto.Message, recv <-chan *proto.Message, song *smf.SMF) error {
	gs.song = song
	gs.send = send
	gs.recv = recv

	// lobby phase
	done := gs.setPhase(proto.GamePhaseLobby, lobbyDur)
lobby:
	select {
	case msg := <-recv:
		switch msg.Type {
		case "join":
			if _, got := gs.players[msg.Player]; !got {
				gs.addPlayer(msg.Player)
				gs.sendTo(msg.Player, "state", []byte("TODO"))
				gs.broadcast("connected", []byte(msg.Player))
			} else {
				gs.players[msg.Player].state = "connected"
			}
		case "leave":
			gs.players[msg.Player].state = "disconnected"
		default:
			log.Printf("unhandled message type '%s'", msg.Type)
		}
		goto lobby
	case <-done:
		fmt.Println("DONE")
	}

	// split tracks
	gs.setPhase(proto.GamePhaseSplitting, 0)
	if err := gs.splitTracksForPlayers(); err != nil {
		return err
	}

	// hero phase
	done = gs.setPhase(proto.GamePhaseHero, recordDur)
hero:
	select {
	case msg := <-recv:
		switch msg.Type {
		case "track":
			smf, err := smf.ReadFrom(bytes.NewReader(msg.Data))
			if err != nil {
				log.Printf("smf parsing error: '%s'", msg.Type)
			}
			gs.players[msg.Player].smf = smf
		default:
			log.Printf("unhandled message type '%s'", msg.Type)
		}
		goto hero
	case <-done:
	}

	// combine tracks
	gs.setPhase(proto.GamePhaseJoining, 0)
	file := gs.combinePlayerTracks()
	var buf bytes.Buffer
	if _, err := file.WriteTo(&buf); err != nil {
		return err
	}
	bs := buf.Bytes()

	// playback phase
	done = gs.setPhase(proto.GamePhasePlayback, playbackDur)
	gs.broadcast("combined", bs)
	<-done

	// done
	gs.setPhase(proto.GamePhaseDone, 0)
	return nil
}

func (gs *GameServer) setPhase(phase proto.GamePhase, dur time.Duration) <-chan time.Time {
	log.Println("phase:", phase)
	gs.broadcast("phase", []byte{byte(phase)})
	gs.phase = phase
	if dur == 0 {
		gs.exp = time.Time{}
		return nil
	}
	gs.exp = time.Now().Add(dur)
	return time.After(time.Until(gs.exp))
}

func (gs *GameServer) addPlayer(fingerprint string) {
	gs.players[fingerprint] = &player{
		state:       "connected",
		fingerprint: fingerprint,
		pianist:     pianists.Hash(fingerprint),
		smf:         nil,
	}
}

func (gs *GameServer) splitTracksForPlayers() error {
	var fingerprints []string
	for f := range gs.players {
		fingerprints = append(fingerprints, f)
	}
	track := abstrack.FromSMF(gs.song.Tracks[0])
	notes := track.CountNotes()
	for i, note := range notes {
		player := fingerprints[i%len(fingerprints)]
		gs.players[player].notes = append(gs.players[player].notes, note.Key)
	}
	for _, player := range gs.players {
		split := smf.New()
		split.Add(track.Select(player.notes).ToSMF())
		var buf bytes.Buffer
		if _, err := split.WriteTo(&buf); err != nil {
			return err
		}
		gs.sendTo(player.fingerprint, "split", buf.Bytes())
	}
	return nil
}

func (gs *GameServer) combinePlayerTracks() *smf.SMF {
	track := abstrack.New()
	for _, player := range gs.players {
		if player.smf == nil {
			continue
		}
		track = track.Merge(abstrack.FromSMF(player.smf.Tracks[0]))
	}
	file := smf.New()
	file.Add(track.ToSMF())
	return file
}

func (gs *GameServer) broadcast(msgType string, data []byte) {
	gs.send <- &proto.Message{
		Type:   msgType,
		Player: "*",
		Data:   data,
	}
}

func (gs *GameServer) sendTo(fingerprint string, msgType string, data []byte) {
	gs.send <- &proto.Message{
		Type:   msgType,
		Player: fingerprint,
		Data:   data,
	}
}
