package gameserver

import (
	"bytes"
	"log"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/pianists"
)

type Message struct {
	Type   string
	Player string
	Data   []byte
}

type player struct {
	state       string
	fingerprint string
	pianist     string
	smf         *smf.SMF
}

type gamePhase byte

const (
	gamePhaseUninitialized gamePhase = iota
	gamePhaseLobby
	gamePhaseSplitting
	gamePhaseHero
	gamePhaseJoining
	gamePhasePlayback
	gamePhaseDone
)

type GameServer struct {
	phase   gamePhase
	exp     time.Time
	send    chan<- *Message
	recv    <-chan *Message
	players map[string]*player
	song    *smf.SMF
}

func New() *GameServer {
	return &GameServer{
		players: map[string]*player{},
		song:    nil,
	}
}

func (gs *GameServer) Start(send chan<- *Message, recv <-chan *Message, song *smf.SMF) {
	// lobby phase
	done := gs.setPhase(gamePhaseLobby, time.Minute)
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
	}

	// split tracks
	gs.setPhase(gamePhaseSplitting, 0)

	// hero phase
	done = gs.setPhase(gamePhaseHero, time.Minute)
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
	gs.setPhase(gamePhaseJoining, 0)
	file := gs.combinePlayerTracks()
	var buf bytes.Buffer
	file.WriteTo(&buf)
	bs := buf.Bytes()

	// playback phase
	gs.setPhase(gamePhasePlayback, 0)
	gs.broadcast("combined", bs)

	//done
	gs.setPhase(gamePhaseDone, 0)
}

func (gs *GameServer) setPhase(phase gamePhase, dur time.Duration) <-chan time.Time {
	gs.phase = phase
	if dur == 0 {
		gs.exp = time.Time{}
		return nil
	}
	gs.exp = time.Now().Add(dur)
	gs.broadcast("phase", []byte{byte(phase)})
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

func (gs *GameServer) splitTracksForPlayers() {
	playerCount := 0
	for _, p := range gs.players {
		if p.state == "connected" {
			playerCount += 1
		}
	}
	track := abstrack.FromSMF(gs.song.Tracks[0])
	notes := track.CountNotes()
	notesPerPlayer := len(notes) / playerCount
	log.Println("TODO", notesPerPlayer)
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
	gs.send <- &Message{
		Type:   msgType,
		Player: "*",
		Data:   data,
	}
}

func (gs *GameServer) sendTo(fingerprint string, msgType string, data []byte) {
	gs.send <- &Message{
		Type:   msgType,
		Player: fingerprint,
		Data:   data,
	}
}
