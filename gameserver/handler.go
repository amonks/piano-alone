package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/db"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/songs"
	"monks.co/piano-alone/templates"
)

type Handler struct {
	db *db.DB

	disklavierConn *websocket.Conn
	conductorConn  *websocket.Conn
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
	db, err := db.OpenDB(os.Getenv("SQLITE_DATABASE_PATH"))
	if err != nil {
		return err
	}

	s.db = db
	gs := NewGame(db)

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

			log.Printf("send: '%s' to '%s'", m.Type, m.Player)
			s.withConn(m.Player, func(conn *websocket.Conn) {
				if conn == nil {
					log.Printf("no such player %s", m.Player)
					return
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, m.Bytes()); err != nil {
					panic(err)
				}
			})
		}
	}()
	return gs.Start(ctx, s.outbox, s.inbox,
		"Prelude in C♯ Minor", "Sergei Rachmaninoff", songs.ExcerptSMF)
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()

	// pages
	mux.Handle("GET /", gzipMiddleware(http.HandlerFunc(s.HandleApp)))
	mux.Handle("GET /download", gzipMiddleware(templ.Handler(templates.Download())))

	// files
	mux.HandleFunc("GET /main.wasm", file("main.wasm"))

	// API: client update
	mux.HandleFunc(data.PathLatestClientVersion, text(data.CurrentVersion))
	mux.HandleFunc(data.PathLatestClientDownload, file("macos-client-universal"))
	// API: ws
	mux.HandleFunc(data.PathPlayerWS, s.HandlePlayerWebsocket)
	mux.HandleFunc(data.PathControllerWS, s.HandleControllerWebsocket)
	// API: performances
	mux.HandleFunc(data.PathSchedulePerformance, s.HandleSchedulePerformance)
	mux.HandleFunc(data.PathScheduledPerformances, s.HandleScheduledPerformances)
	mux.HandleFunc(data.PathFeaturedPerformances, s.HandleFeaturedPerformances)
	mux.HandleFunc(data.PathBeginPerformance, s.HandleBeginPerformance)
	mux.HandleFunc(data.PathDeletePerformance, s.HandleDeletePerformance)
	mux.HandleFunc(data.PathMIDIFile, s.HandleMIDIFile)
	// API: current performance
	mux.HandleFunc(data.PathRestart, s.HandleRestart)
	mux.HandleFunc(data.PathAdvance, s.HandleAdvance)

	mux.ServeHTTP(w, req)
}

func (s *Handler) addConn(fingerprint string, conn *websocket.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	switch fingerprint {
	case "disklavier":
		s.disklavierConn = conn
	case "conductor":
		s.conductorConn = conn
	default:
		s.conns[fingerprint] = conn
	}
}
func (s *Handler) removeConn(fingerprint string) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	switch fingerprint {
	case "disklavier":
		s.disklavierConn = nil
		go func() { s.inbox <- game.NewMessage(game.MessageTypeDisklavierDisconnected, "", nil) }()
	case "conductor":
		s.conductorConn = nil
		go func() { s.inbox <- game.NewMessage(game.MessageTypeConductorDisconnected, "", nil) }()
	default:
		delete(s.conns, fingerprint)
	}
}
func (s *Handler) withConn(fingerprint string, f func(*websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	switch fingerprint {
	case "controllers":
		if s.disklavierConn != nil {
			f(s.disklavierConn)
		}
		if s.conductorConn != nil {
			f(s.conductorConn)
		}
	case "disklavier":
		if s.disklavierConn != nil {
			f(s.disklavierConn)
		}
	case "conductor":
		if s.conductorConn != nil {
			f(s.conductorConn)
		}
	default:
		if c, ok := s.conns[fingerprint]; ok {
			f(c)
		}
	}
}
func (s *Handler) eachConn(f func(string, *websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	if s.disklavierConn != nil {
		f("disklavier", s.disklavierConn)
	}
	if s.conductorConn != nil {
		f("conductor", s.conductorConn)
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

func (s *Handler) HandleApp(w http.ResponseWriter, req *http.Request) {
	ps, err := s.db.GetFeaturedPerformances()
	if err != nil {
		http.Error(w, "internal error", 500)
	}
	h := templ.Handler(templates.App(ps))
	h.ServeHTTP(w, req)
}

var upgrader = websocket.Upgrader{}

func (s *Handler) HandleRestart(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathRestart)
	s.inbox <- game.NewMessage(game.MessageTypeRestart, "", nil)
	w.Write([]byte("ok"))
}

func (s *Handler) HandleAdvance(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathAdvance)
	s.inbox <- game.NewMessage(game.MessageTypeAdvancePhase, "", nil)
	w.Write([]byte("ok"))
}

func (s *Handler) HandleSchedulePerformance(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathSchedulePerformance)
	bs, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	p := game.PerformanceFromBytes(bs)
	if err := s.db.SchedulePerformance(p); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("ok"))
}

func (s *Handler) HandleFeaturedPerformances(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathFeaturedPerformances)
	ps, err := s.db.GetFeaturedPerformances()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	h := templ.Handler(templates.Performances(ps))
	h.ServeHTTP(w, req)
}

func (s *Handler) HandleScheduledPerformances(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathScheduledPerformances)
	ps, err := s.db.GetScheduledPerformances()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	bs := game.PerformancesToBytes(ps)
	w.Write(bs)
}

func (s *Handler) HandleBeginPerformance(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathBeginPerformance)
	id := req.PathValue("id")
	perf, err := s.db.GetPerformance(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.inbox <- game.NewMessage(game.MessageTypeBeginPerformance, "", perf.Configuration.Bytes())
	w.Write([]byte("ok"))
}

func (s *Handler) HandleMIDIFile(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathMIDIFile)
	id := req.PathValue("id")
	fmt.Println("id", id)
	bs, err := s.db.GetMIDIFile(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-type", "audio/midi")
	w.Write(bs)
}

func (s *Handler) HandleDeletePerformance(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathDeletePerformance)
	id := req.PathValue("id")
	if err := s.db.DeletePerformance(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("ok"))
}

func (s *Handler) HandleControllerWebsocket(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathControllerWS)
	role := req.URL.Query().Get("role")
	if role != "conductor" {
		role = "disklavier"
	}
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ws upgrade error: %s\n", err.Error())
		return
	}
	s.addConn(role, c)

	for {
		_, bs, err := c.ReadMessage()
		if err != nil {
			c.Close()
			break
		}

		m := game.MessageFromBytes(bs)
		s.inbox <- m
	}
	s.removeConn(role)
}

func (s *Handler) HandlePlayerWebsocket(w http.ResponseWriter, req *http.Request) {
	log.Println("<-", data.PathPlayerWS)
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
