package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/songs"
	"monks.co/piano-alone/templates"
)

type Handler struct {
	controllerConn *websocket.Conn
	conns          map[string]*websocket.Conn
	connMu         sync.RWMutex

	inbox  chan *game.Message
	outbox chan *game.Message
}

func NewHandler() *Handler {
	return &Handler{
		conns:  map[string]*websocket.Conn{},
		outbox: make(chan *game.Message),
		inbox:  make(chan *game.Message),
	}
}

func (s *Handler) Start(ctx context.Context) error {
	gs := NewGame()
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
	f := songs.ExcerptSMF
	gs.Start(ctx, s.outbox, s.inbox, f)
	return nil
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()

	// pages
	mux.Handle("GET /", gzipMiddleware(templ.Handler(templates.ComingSoon())))
	mux.Handle("GET /app", gzipMiddleware(templ.Handler(templates.App())))
	mux.Handle("GET /download", gzipMiddleware(templ.Handler(templates.Download())))

	// files
	mux.HandleFunc("GET /main.wasm", file("main.wasm"))

	// API
	mux.HandleFunc(data.PathLatestClientVersion, text(data.CurrentVersion))
	mux.HandleFunc(data.PathLatestClientDownload, file("macos-client-universal"))
	mux.HandleFunc(data.PathPlayerWS, s.HandlePlayerWebsocket)
	mux.HandleFunc(data.PathControllerWS, s.HandleControllerWebsocket)
	mux.HandleFunc(data.PathRestart, s.HandleRestart)
	mux.HandleFunc(data.PathAdvance, s.HandleAdvance)

	mux.ServeHTTP(w, req)
}

func (s *Handler) addControllerConn(conn *websocket.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	s.controllerConn = conn
}
func (s *Handler) removeControllerConn() {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	s.controllerConn = nil
}

func (s *Handler) addConn(fingerprint string, conn *websocket.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	s.conns[fingerprint] = conn
}
func (s *Handler) removeConn(fingerprint string) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	delete(s.conns, fingerprint)
}
func (s *Handler) withConn(fingerprint string, f func(*websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	if fingerprint == "controller" {
		f(s.controllerConn)
	} else {
		f(s.conns[fingerprint])
	}
}
func (s *Handler) eachConn(f func(string, *websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	if s.controllerConn != nil {
		f("controller", s.controllerConn)
	}
	for fingerprint, sock := range s.conns {
		f(fingerprint, sock)
	}
}

func file(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("<- %s", req.URL.Path)
		http.ServeFile(w, req, filename)
	}
}

func text(t string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("<- %s", req.URL.Path)
		w.Write([]byte(t))
	}
}

var upgrader = websocket.Upgrader{}

func (s *Handler) HandleRestart(w http.ResponseWriter, req *http.Request) {
	s.inbox <- game.NewMessage(game.MessageTypeRestart, "", nil)
	w.Write([]byte("ok"))
}

func (s *Handler) HandleAdvance(w http.ResponseWriter, req *http.Request) {
	s.inbox <- game.NewMessage(game.MessageTypeAdvancePhase, "", nil)
	w.Write([]byte("ok"))
}

func (s *Handler) HandleControllerWebsocket(w http.ResponseWriter, req *http.Request) {
	log.Printf("<- controller")
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ws upgrade error: %s\n", err.Error())
		return
	}
	s.addControllerConn(c)

	for {
		_, bs, err := c.ReadMessage()
		if err != nil {
			c.Close()
			break
		}

		m := game.MessageFromBytes(bs)
		s.inbox <- m
	}
	s.removeControllerConn()
}

func (s *Handler) HandlePlayerWebsocket(w http.ResponseWriter, req *http.Request) {
	log.Printf("<- player")
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
	s.inbox <- game.NewMessage(game.MessageTypeLeave, fingerprint, nil)
}
