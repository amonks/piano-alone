package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/gameserver"
	"monks.co/piano-alone/proto"
)

type Server struct {
	players map[string]*Player

	inbox  chan *proto.Message
	outbox chan *proto.Message
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	gs := gameserver.New()
	s.players = map[string]*Player{}
	s.outbox = make(chan *proto.Message)
	s.inbox = make(chan *proto.Message)
	go func() {
		for m := range s.outbox {
			if m.Player == "*" {
				for _, player := range s.players {
					if err := player.conn.WriteMessage(websocket.BinaryMessage, m.Bytes()); err != nil {
						panic(err)
					}
				}
				continue
			}

			player, got := s.players[m.Player]
			if !got {
				panic("no such player " + m.Player)
			}
			if err := player.conn.WriteMessage(websocket.BinaryMessage, m.Bytes()); err != nil {
				panic(err)
			}
		}
	}()
	f, err := smf.ReadFile("example.mid")
	if err != nil {
		return err
	}
	return gs.Start(s.outbox, s.inbox, f)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./website")))
	mux.HandleFunc("/ws", s.HandleWebsocket)
	mux.ServeHTTP(w, req)
}

var upgrader = websocket.Upgrader{}

type Player struct {
	conn *websocket.Conn
}

func (s *Server) HandleWebsocket(w http.ResponseWriter, req *http.Request) {
	fingerprint := req.URL.Query().Get("fingerprint")
	if fingerprint == "" {
		http.Error(w, "no fingerprint specified", 400)
		return
	}
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, "ws upgrade error: "+err.Error(), 500)
		return
	}
	player := &Player{c}
	s.players[fingerprint] = player

	for {
		_, bs, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %s", err)
			c.Close()
			break
		}

		m := proto.MessageFromBytes(bs)
		s.inbox <- m
	}
}
