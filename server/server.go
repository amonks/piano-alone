package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/gameserver"
)

type Server struct {
	conns  map[string]*websocket.Conn
	connMu sync.RWMutex

	inbox  chan *game.Message
	outbox chan *game.Message
}

func (s *Server) addConn(fingerprint string, conn *websocket.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	s.conns[fingerprint] = conn
}
func (s *Server) removeConn(fingerprint string) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	delete(s.conns, fingerprint)
}
func (s *Server) withConn(fingerprint string, f func(*websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	f(s.conns[fingerprint])
}
func (s *Server) eachConn(f func(string, *websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	for fingerprint, sock := range s.conns {
		f(fingerprint, sock)
	}
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	gs := gameserver.New()
	s.conns = map[string]*websocket.Conn{}
	s.outbox = make(chan *game.Message)
	s.inbox = make(chan *game.Message)
	go func() {
		for m := range s.outbox {
			if m.Player == "*" {
				log.Printf("broadcast: '%s'", m.Type)
				s.eachConn(func(fingerprint string, sock *websocket.Conn) {
					if err := sock.WriteMessage(websocket.BinaryMessage, m.Bytes()); err != nil {
						panic(err)
					}
				})
				continue
			}

			log.Printf("send: '%s'", m.Type)
			s.withConn(m.Player, func(conn *websocket.Conn) {
				if conn == nil {
					panic("no such player " + m.Player)
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, m.Bytes()); err != nil {
					panic(err)
				}
			})
		}
	}()
	f, err := smf.ReadFile("example-2.mid")
	if err != nil {
		return err
	}
	gs.Start(s.outbox, s.inbox, f)
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./website")))
	mux.HandleFunc("/ws", s.HandleWebsocket)
	mux.ServeHTTP(w, req)
}

var upgrader = websocket.Upgrader{}

func (s *Server) HandleWebsocket(w http.ResponseWriter, req *http.Request) {
	fingerprint := req.URL.Query().Get("fingerprint")
	if fingerprint == "" {
		http.Error(w, "no fingerprint specified", 400)
		return
	}
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ws upgrade error: %s\n", err.Error())
		return
	}
	s.addConn(fingerprint, c)

	for {
		_, bs, err := c.ReadMessage()
		if err != nil {
			c.Close()
			break
		}

		m := game.MessageFromBytes(bs)
		s.inbox <- m
	}
	s.removeConn(fingerprint)
	s.inbox <- &game.Message{
		Type:   game.MessageTypeLeave,
		Player: fingerprint,
	}
}
