package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"

	"monks.co/piano-alone/assets"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/db"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/templates"
)

// Request-body limits. A performance carries a whole MIDI score and a
// socket message carries a player's recorded part, so these are
// generous — they are here to bound what an anonymous caller can make
// the server allocate, not to police the piece.
const (
	maxPerformanceBytes   = 8 << 20
	maxSocketMessageBytes = 8 << 20
)

// newMux builds the route table once. It used to be rebuilt on every
// request, which worked but re-registered fourteen patterns per page
// view.
func (s *Server) newMux() http.Handler {
	mux := http.NewServeMux()

	// pages
	mux.HandleFunc("GET /", s.handleApp)
	mux.HandleFunc("GET /download", s.page(templates.Download()))

	// files
	mux.Handle("GET /main.wasm", s.asset("main.wasm", "application/wasm"))

	// the disklavier client, and the version it checks for updates
	mux.HandleFunc(data.PathLatestClientVersion, s.handleClientVersion)
	mux.HandleFunc(data.PathLatestClientDownload, s.handleClientDownload)

	// sockets
	mux.HandleFunc(data.PathPlayerWS, s.handlePlayerWebsocket)
	mux.HandleFunc(data.PathControllerWS, s.handleControllerWebsocket)

	// performances
	mux.HandleFunc(data.PathSchedulePerformance, s.handleSchedulePerformance)
	mux.HandleFunc(data.PathScheduledPerformances, s.handleScheduledPerformances)
	mux.HandleFunc(data.PathFeaturedPerformances, s.handleFeaturedPerformances)
	mux.HandleFunc(data.PathBeginPerformance, s.handleBeginPerformance)
	mux.HandleFunc(data.PathDeletePerformance, s.handleDeletePerformance)
	mux.HandleFunc(data.PathMIDIFile, s.handleMIDIFile)

	// the current performance
	mux.HandleFunc(data.PathRestart, s.handleRestart)
	mux.HandleFunc(data.PathAdvance, s.handleAdvance)

	return mux
}

func (s *Server) render(w http.ResponseWriter, req *http.Request, c templ.Component) {
	ctx := req.Context()
	if s.bodyEnd != nil {
		ctx = templates.WithBodyEnd(ctx, func() template.HTML { return s.bodyEnd(req) })
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.Render(ctx, w); err != nil {
		s.fail(w, req, err, "rendering page")
	}
}

func (s *Server) page(c templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) { s.render(w, req, c) }
}

// fail answers an error without telling the caller what it was. The
// detail goes to the log, where an operator can read it; the response
// says only that the request did not work, because these pages are
// public and a SQL error is not the audience's business.
func (s *Server) fail(w http.ResponseWriter, req *http.Request, err error, doing string) {
	s.log.ErrorContext(req.Context(), "request failed",
		"http.route", req.Pattern,
		"error.message", fmt.Sprintf("%s: %s", doing, err),
		"error.type", fmt.Sprintf("%T", err))
	http.Error(w, "internal error", http.StatusInternalServerError)
}

// notFound answers a performance id that is not in the table. It is
// separate from fail because a bad id is the caller's mistake and a
// broken database is ours.
func (s *Server) notFound(w http.ResponseWriter, req *http.Request, err error) bool {
	if !errors.Is(err, db.ErrNotFound) {
		return false
	}
	http.Error(w, "no such performance", http.StatusNotFound)
	return true
}

func (s *Server) asset(name, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bs, err := assets.Read(name)
		if errors.Is(err, fs.ErrNotExist) {
			// The wasm player is built, not committed. A server built
			// without running the generators serves the page and this
			// 404 rather than failing to compile, so say plainly which
			// artifact is missing.
			s.log.ErrorContext(req.Context(), "asset missing from build", "piano.asset", name)
			http.Error(w, name+" was not built", http.StatusNotFound)
			return
		} else if err != nil {
			s.fail(w, req, err, "reading asset "+name)
			return
		}
		w.Header().Set("Content-Type", contentType)
		// A zero modtime keeps Last-Modified off the response: the
		// bytes are baked into the binary, so the only honest answer
		// is the build's, which this package does not know.
		http.ServeContent(w, req, name, time.Time{}, bytes.NewReader(bs))
	})
}

func (s *Server) handleApp(w http.ResponseWriter, req *http.Request) {
	ps, err := s.db.GetFeaturedPerformances(req.Context())
	if err != nil {
		s.fail(w, req, err, "reading featured performances")
		return
	}
	s.render(w, req, templates.App(ps))
}

func (s *Server) handleClientVersion(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, data.CurrentVersion)
}

// handleClientDownload serves the signed universal Mac binary that
// drives the disklavier. It is built and notarized out of band on a
// Mac and handed to the server as a file, so a deployment without one
// says so rather than serving nothing under a 200.
func (s *Server) handleClientDownload(w http.ResponseWriter, req *http.Request) {
	if s.clientPath == "" {
		http.Error(w, "no client binary is configured", http.StatusNotFound)
		return
	}
	f, err := os.Open(s.clientPath)
	if errors.Is(err, fs.ErrNotExist) {
		s.log.WarnContext(req.Context(), "client binary is configured but absent", "piano.client_path", s.clientPath)
		http.Error(w, "no client binary is available", http.StatusNotFound)
		return
	} else if err != nil {
		s.fail(w, req, err, "opening client binary")
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		s.fail(w, req, err, "reading client binary")
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, req, "piano-telephone", info.ModTime(), f)
}

func (s *Server) handleSchedulePerformance(w http.ResponseWriter, req *http.Request) {
	bs, err := io.ReadAll(http.MaxBytesReader(w, req.Body, maxPerformanceBytes))
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}
	p, err := game.PerformanceFromBytes(bs)
	if err != nil {
		http.Error(w, "malformed performance", http.StatusBadRequest)
		return
	}
	if p.Configuration == nil || p.Configuration.PerformanceID == "" {
		http.Error(w, "performance has no id", http.StatusBadRequest)
		return
	}
	if err := s.db.SchedulePerformance(req.Context(), p); err != nil {
		s.fail(w, req, err, "scheduling performance")
		return
	}
	s.log.InfoContext(req.Context(), "performance scheduled",
		"piano.performance_id", p.Configuration.PerformanceID,
		"piano.title", p.Configuration.Title)
	io.WriteString(w, "ok")
}

func (s *Server) handleFeaturedPerformances(w http.ResponseWriter, req *http.Request) {
	ps, err := s.db.GetFeaturedPerformances(req.Context())
	if err != nil {
		s.fail(w, req, err, "reading featured performances")
		return
	}
	s.render(w, req, templates.Performances(ps))
}

func (s *Server) handleScheduledPerformances(w http.ResponseWriter, req *http.Request) {
	ps, err := s.db.GetScheduledPerformances(req.Context())
	if err != nil {
		s.fail(w, req, err, "reading scheduled performances")
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(game.PerformancesToBytes(ps))
}

func (s *Server) handleBeginPerformance(w http.ResponseWriter, req *http.Request) {
	perf, err := s.db.GetPerformance(req.Context(), req.PathValue("id"))
	if err != nil {
		if s.notFound(w, req, err) {
			return
		}
		s.fail(w, req, err, "reading performance")
		return
	}
	s.control(req, game.MessageTypeBeginPerformance, perf.Configuration.Bytes())
	io.WriteString(w, "ok")
}

func (s *Server) handleMIDIFile(w http.ResponseWriter, req *http.Request) {
	bs, err := s.db.GetMIDIFile(req.Context(), req.PathValue("id"))
	if err != nil {
		if s.notFound(w, req, err) {
			return
		}
		s.fail(w, req, err, "reading rendition")
		return
	}
	if len(bs) == 0 {
		http.Error(w, "that performance has no rendition yet", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "audio/midi")
	w.Write(bs)
}

func (s *Server) handleDeletePerformance(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if err := s.db.DeletePerformance(req.Context(), id); err != nil {
		s.fail(w, req, err, "deleting performance")
		return
	}
	s.log.InfoContext(req.Context(), "performance deleted", "piano.performance_id", id)
	io.WriteString(w, "ok")
}

func (s *Server) handleRestart(w http.ResponseWriter, req *http.Request) {
	s.control(req, game.MessageTypeRestart, nil)
	io.WriteString(w, "ok")
}

func (s *Server) handleAdvance(w http.ResponseWriter, req *http.Request) {
	s.control(req, game.MessageTypeAdvancePhase, nil)
	io.WriteString(w, "ok")
}

// control puts a conducting message on the game loop's inbox. These
// arrive over HTTP, where the host has already made its decision about
// the request, so they skip the origin check messages off a socket go
// through.
func (s *Server) control(req *http.Request, t game.MessageType, payload []byte) {
	s.send(req.Context(), game.NewMessage(t, roleConductor, payload))
}

// upgrader accepts a cross-origin upgrade, because the check gorilla
// does by default compares the Origin header's host against the
// request's Host — and behind a reverse proxy the request's Host is
// the backend's name, not the name the browser typed, so the default
// refuses every real connection. A host that cares about which origins
// may open a socket has to enforce it in front of this package, where
// the forwarded name is known.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

func (s *Server) handleControllerWebsocket(w http.ResponseWriter, req *http.Request) {
	role := roleDisklavier
	if req.URL.Query().Get("role") == roleConductor {
		role = roleConductor
	}
	s.socket(w, req, role)
}

func (s *Server) handlePlayerWebsocket(w http.ResponseWriter, req *http.Request) {
	fingerprint := req.URL.Query().Get("fingerprint")
	if fingerprint == "" {
		http.Error(w, "no fingerprint specified", http.StatusBadRequest)
		return
	}
	if fingerprint == roleConductor || fingerprint == roleDisklavier ||
		fingerprint == fingerprintControllers || fingerprint == fingerprintEveryone {
		// The routing vocabulary and the fingerprint space share one
		// namespace, so a player claiming one of these names would be
		// handed the controllers' messages.
		http.Error(w, "reserved fingerprint", http.StatusBadRequest)
		return
	}
	s.socket(w, req, fingerprint)
}

// socket runs one connection's read loop. origin is who the connection
// is — a player's fingerprint, or one of the two controller roles —
// and it is stamped onto every message the connection sends. Nothing a
// client writes decides who it is: the Player field arrives from the
// internet, and believing it let any socket act as any player.
func (s *Server) socket(w http.ResponseWriter, req *http.Request, origin string) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		// Upgrade has already written its own response.
		s.log.WarnContext(req.Context(), "websocket upgrade failed",
			"piano.origin", origin,
			"error.message", err.Error(),
			"error.type", fmt.Sprintf("%T", err))
		return
	}
	ctx := req.Context()
	s.addConn(ctx, origin, conn)
	s.log.InfoContext(ctx, "socket opened", "piano.origin", origin)

	conn.SetReadLimit(maxSocketMessageBytes)
	for {
		_, bs, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		m, err := game.MessageFromBytes(bs)
		if err != nil {
			s.log.WarnContext(ctx, "dropping malformed socket message",
				"piano.origin", origin,
				"error.message", err.Error(),
				"error.type", fmt.Sprintf("%T", err))
			continue
		}
		s.fromClient(ctx, origin, m)
	}

	s.removeConn(ctx, origin, conn)
	s.log.InfoContext(ctx, "socket closed", "piano.origin", origin)
}

// fromClient stamps a socket message with its origin and refuses the
// ones that origin has no business sending. Only the conductor's
// socket can drive the performance; the disklavier's socket receives a
// rendition and says when it connects, and a player's socket speaks
// only for that player. Every message used to reach the loop from
// every socket, so anyone who could open one could advance a phase or
// submit another player's part.
func (s *Server) fromClient(ctx context.Context, origin string, m *game.Message) bool {
	if m.Type.Control() && origin != roleConductor {
		s.log.WarnContext(ctx, "refusing control message from non-conductor",
			"piano.origin", origin, "piano.message_type", m.Type.String())
		return false
	}
	if m.Type.ServerOnly() {
		// The presence flags are the server's own account of which
		// controller sockets are open, synthesized where they open and
		// close. No client has business asserting one, and a player
		// that did could tell the conductor a disklavier was connected
		// that is not there.
		s.log.WarnContext(ctx, "refusing server-only message from a client",
			"piano.origin", origin, "piano.message_type", m.Type.String())
		return false
	}
	m.Player = origin
	return s.send(ctx, m)
}

func (s *Server) addConn(ctx context.Context, fingerprint string, conn *websocket.Conn) {
	s.connMu.Lock()
	switch fingerprint {
	case roleDisklavier:
		s.disklavierConn = conn
	case roleConductor:
		s.conductorConn = conn
	default:
		s.conns[fingerprint] = conn
	}
	s.connMu.Unlock()

	switch fingerprint {
	case roleDisklavier:
		s.send(ctx, game.NewMessage(game.MessageTypeDisklavierConnected, roleDisklavier, nil))
	case roleConductor:
		s.send(ctx, game.NewMessage(game.MessageTypeConductorConnected, roleConductor, nil))
	}
}

// removeConn clears a slot only if it still holds the connection that
// is leaving. A reconnect — a page refresh mid-performance — registers
// the new socket before the old one's read loop notices it is closed,
// and a delete keyed on the fingerprint alone would take the new
// socket out of the registry: the player would stay connected and
// receive nothing, no assignment and no phase changes, while the
// server thought it had already told them.
//
// The two controller slots need the same guard, and for the same
// reason plus one: an impostor who grabs the disklavier slot and drops
// it would otherwise nil out whatever conn is registered by then,
// including a real disklavier that reconnected in between.
func (s *Server) removeConn(ctx context.Context, fingerprint string, conn *websocket.Conn) {
	s.connMu.Lock()
	var current *websocket.Conn
	switch fingerprint {
	case roleDisklavier:
		if current = s.disklavierConn; current == conn {
			s.disklavierConn = nil
		}
	case roleConductor:
		if current = s.conductorConn; current == conn {
			s.conductorConn = nil
		}
	default:
		if current = s.conns[fingerprint]; current == conn {
			delete(s.conns, fingerprint)
		}
	}
	s.connMu.Unlock()

	// Someone else holds the slot now, so there is nothing to announce:
	// they are connected, and saying otherwise would mark them gone.
	if current != conn {
		return
	}

	// A non-blocking send, because this runs on the way out of a
	// hijacked request. Once Run has returned nothing drains the
	// inbox, and a hijacked request's context is cancelled only when
	// its handler returns — which it cannot do while parked here. The
	// game loop is already stopped in that case, so the message it
	// would have received has nowhere to be useful.
	msg := game.NewMessage(game.MessageTypeLeave, fingerprint, nil)
	switch fingerprint {
	case roleDisklavier:
		msg = game.NewMessage(game.MessageTypeDisklavierDisconnected, roleDisklavier, nil)
	case roleConductor:
		msg = game.NewMessage(game.MessageTypeConductorDisconnected, roleConductor, nil)
	}
	select {
	case s.inbox <- msg:
	case <-ctx.Done():
	}
}

func (s *Server) withConn(fingerprint string, f func(string, *websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	switch fingerprint {
	case fingerprintControllers:
		if s.disklavierConn != nil {
			f(roleDisklavier, s.disklavierConn)
		}
		if s.conductorConn != nil {
			f(roleConductor, s.conductorConn)
		}
	case roleDisklavier:
		if s.disklavierConn != nil {
			f(roleDisklavier, s.disklavierConn)
		}
	case roleConductor:
		if s.conductorConn != nil {
			f(roleConductor, s.conductorConn)
		}
	default:
		if c, ok := s.conns[fingerprint]; ok {
			f(fingerprint, c)
		}
	}
}

func (s *Server) eachConn(f func(string, *websocket.Conn)) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	if s.disklavierConn != nil {
		f(roleDisklavier, s.disklavierConn)
	}
	if s.conductorConn != nil {
		f(roleConductor, s.conductorConn)
	}
	for fingerprint, sock := range s.conns {
		f(fingerprint, sock)
	}
}
