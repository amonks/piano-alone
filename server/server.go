package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	state    serverState
	commands chan ([]byte)

	playersMu sync.RWMutex
	players   map[string]*PlayerSession
}

func New() *Server {
	return &Server{
		commands: make(chan []byte),
		players:  map[string]*PlayerSession{},
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.HandleWebsocket)
	mux.Handle("/", http.FileServer(http.Dir("./website")))
	mux.ServeHTTP(w, req)
}

type serverState int

const (
	serverStateUninitialized serverState = iota
	serverStateWaitingForPlayers
	serverStateRecording
	serverStatePlaying
	serverStateDone
)

var upgrader = websocket.Upgrader{}

func (s *Server) FindOrCreatePlayerSession(fingerprint string) *PlayerSession {
	if got := s.getPlayerSession(fingerprint); got != nil {
		return got
	}
	return s.createPlayerSession(fingerprint)
}

func (s *Server) getPlayerSession(fingerprint string) *PlayerSession {
	s.playersMu.RLock()
	defer s.playersMu.RUnlock()
	return s.players[fingerprint]
}

func (s *Server) createPlayerSession(fingerprint string) *PlayerSession {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()
	ps := NewPlayerSession(fingerprint)
	s.players[fingerprint] = ps
	return ps
}

func (s *Server) HandleWebsocket(w http.ResponseWriter, req *http.Request) {
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ws upgrade error: %s", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %s", err)
			break
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}

	}
}
