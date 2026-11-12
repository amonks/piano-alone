// Package server is piano-alone's web server and game loop: the pages
// the audience sees, the sockets the players and the two controllers
// hold open, and the state machine a conductor advances through a
// performance.
//
// It authorizes nothing. Handler returns everything on one mux and
// Operation says which of two kinds a request is, so a host can put
// its own decision in front of the whole thing and a route added here
// later arrives gated. Run standalone, with no host in front, every
// conductor verb is open — which is what the piece did for its first
// two performances, and is still the right default for someone running
// it on a laptop in a room.
package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"monks.co/piano-alone/assets"
	"monks.co/piano-alone/db"
	"monks.co/piano-alone/game"
)

// The two controller roles, which are also the names the game loop
// addresses them by beside the players' ids, and the two names that
// address more than one socket.
const (
	roleConductor  = "conductor"
	roleDisklavier = "disklavier"

	// toControllers addresses both controllers at once.
	toControllers = "controllers"

	// toEveryone broadcasts.
	toEveryone = "*"
)

const defaultDBPath = "piano-alone.db"

type Server struct {
	db     *db.DB
	ownsDB bool
	log    *slog.Logger

	bodyEnd    func(*http.Request) template.HTML
	clientPath string

	mux http.Handler

	// The open sockets, by whose they are: one for each controller's
	// role and one for each player's id.
	disklavierConn *peer
	conductorConn  *peer
	conns          map[string]*peer
	connMu         sync.RWMutex
	// seats holds one token for each player's socket that is open.
	seats chan struct{}

	inbox  chan *game.Message
	outbox chan *game.Message

	// The performance. These belong to the game loop and are touched
	// from nowhere else, which is why they carry no lock: play handles
	// one message at a time. All of it is in memory — only the
	// finished rendition is written down — so a restart mid-performance
	// loses everything in flight.
	state *game.State
	// partials are the parts players have submitted, by id, as
	// the bytes each sent.
	partials               map[string][]byte
	sentSwitchToVideoModal bool
}

// New builds a server. It starts nothing: Run drives the game loop and
// Handler serves the pages, and a host runs them under whatever
// supervision it already has.
func New(ctx context.Context, opts Options) (*Server, error) {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}

	var (
		store  *db.DB
		ownsDB bool
	)
	if opts.DB != nil {
		store = db.New(opts.DB)
	} else {
		path := opts.DBPath
		if path == "" {
			path = defaultDBPath
		}
		opened, err := db.Open(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("opening database at %s: %w", path, err)
		}
		if err := opened.Seed(ctx); err != nil {
			opened.Close()
			return nil, fmt.Errorf("seeding database: %w", err)
		}
		store, ownsDB = opened, true
	}

	s := &Server{
		db:         store,
		ownsDB:     ownsDB,
		log:        log,
		bodyEnd:    opts.BodyEnd,
		clientPath: opts.ClientPath,
		conns:      map[string]*peer{},
		seats:      make(chan struct{}, maxPlayerSockets),
		outbox:     make(chan *game.Message),
		inbox:      make(chan *game.Message),
	}
	s.mux = s.newMux()

	// The player is a build artifact, and the embed is of a directory
	// so that a checkout which has not run the generators still
	// compiles. That means a server can start without the thing it
	// exists to serve, which is worth saying once at boot rather than
	// leaving to whoever opens the page.
	if _, err := assets.Read("main.wasm"); err != nil {
		log.WarnContext(ctx, "the wasm player is not built: the page will load without it",
			"piano.asset", "main.wasm")
	}

	return s, nil
}

// Close releases what New opened. A server given a database handle
// closes nothing: the handle is its caller's.
func (s *Server) Close() error {
	if s.ownsDB {
		return s.db.Close()
	}
	return nil
}

// Handler serves the pages, the assets, the sockets, and the
// performance API, all rooted at "/". The server expects to own the
// whole namespace of whatever it is mounted on — this is a site, not a
// section of one — so there is no sub-path support.
func (s *Server) Handler() http.Handler { return s.mux }

// Run drives the game loop until ctx is cancelled. Nothing else in
// this package advances a performance, so a server whose Run is not
// running serves pages and accepts sockets that go nowhere.
func (s *Server) Run(ctx context.Context) error {
	go s.deliver(ctx)
	return s.play(ctx)
}

// Serve runs the HTTP server on addr until ctx is cancelled. It does
// not run the game loop; see Run.
func (s *Server) Serve(ctx context.Context, addr string) error {
	srv := &http.Server{Addr: addr, Handler: gzipMiddleware(s.Handler())}
	errs := make(chan error, 1)
	go func() { errs <- srv.ListenAndServe() }()
	select {
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
		defer cancel()
		return srv.Shutdown(shutdown)
	case err := <-errs:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// Go runs the game loop and the HTTP server together until ctx is
// cancelled or either fails. It is the standalone server; a host that
// owns its own listener runs Run and mounts Handler instead.
func (s *Server) Go(ctx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, 2)
	// Both of these return nil on a clean shutdown, and wrapping a nil
	// error yields an error that formats as "%!w(<nil>)" — which is
	// what the standalone command would have printed on its way to
	// exiting 1 every time someone pressed ^C.
	labelled := func(what string, err error) error {
		if err == nil {
			return nil
		}
		return fmt.Errorf("%s: %w", what, err)
	}
	go func() { errs <- labelled("game loop", s.Run(ctx)) }()
	go func() { errs <- labelled("http", s.Serve(ctx, addr)) }()

	s.log.Info("piano-alone listening", "server.address", addr)
	err := <-errs
	cancel()
	if second := <-errs; err == nil {
		err = second
	}
	return err
}

// What the sockets may cost, whoever opens them. A player's socket is
// public and anonymous, and each open one is two goroutines, a
// connection's buffers, and the frames waiting to be written to it.
const (
	// maxPlayerSockets is how many players' sockets are open at once:
	// one for every player a performance holds, and a quarter again
	// for pages caught reloading. The piece's audiences were twelve
	// and six. The controllers' sockets are not counted, so a full
	// house cannot keep the conductor or the disklavier out.
	maxPlayerSockets = maxPlayers + maxPlayers/4

	// playerQueue is how many frames may wait for a player's socket.
	// A whole performance sends one player about eight — the state,
	// its notes, each phase — so a socket this far behind has missed
	// the performance.
	playerQueue = 8

	// controllerQueue is how many frames may wait for a controller's
	// socket. A controller is told of every player, so it holds twice
	// what every player's socket at once could have it told
	// (inboundBurst messages each). These frames are small: a
	// player's id, or its record.
	controllerQueue = 2 * maxPlayerSockets * inboundBurst

	// writeTimeout bounds one socket write. A phone on bad venue wifi
	// is given this long to take a frame; a socket that takes longer
	// is closed.
	writeTimeout = 10 * time.Second

	// pingInterval is how often each socket is pinged, and readTimeout
	// how long one may go without answering — or sending anything —
	// before it is closed. A peer that has gone says nothing, and a
	// socket nothing is written to never fails, so without these a
	// phone out of range held its socket until the server restarted.
	// Three pings may go missing, since venue wifi drops a phone for
	// longer than it drops a packet, and a page does not reconnect by
	// itself: closing a socket that was only slow costs its player the
	// performance until they reload.
	pingInterval = 30 * time.Second
	readTimeout  = 3 * pingInterval

	// inboundBurst is how many of a socket's messages are read at
	// once, and inboundGap how far apart the rest are. A page sends
	// four messages in a performance. The loop tells the controllers
	// of most of them, so an unpaced socket could have a controller
	// told of it thousands of times a second, faster than the
	// controller's socket drains — and a socket that falls behind is
	// closed.
	inboundBurst = 8
	inboundGap   = time.Second

	// closeReplaced is the close code a socket is sent when a newer
	// one takes its place: a private-use code, since no standard one
	// says it.
	closeReplaced = 4000
)

// peer is one open socket: the connection, and the frames waiting for
// its writer.
type peer struct {
	conn  *websocket.Conn
	queue chan frame
	// gone is closed when the socket's read loop ends.
	gone chan struct{}
}

// frame is a message as the sockets carry it, encoded once however
// many sockets it is for.
type frame struct {
	kind  game.MessageType
	bytes []byte
}

// deliver fans the game loop's outbound messages out to the sockets.
// It waits on none of them. Each socket has a writer of its own, and
// delivery only queues a frame for it: this used to write to every
// socket in turn, so a socket that stopped reading held each message
// for ten seconds, and with it every other socket and the game loop,
// which hands delivery one message at a time. Anyone may open a
// player's socket and not read from it.
func (s *Server) deliver(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-s.outbox:
			f := frame{kind: m.Type, bytes: m.Bytes()}
			if m.Player == toEveryone {
				s.eachConn(func(id string, p *peer) { s.enqueue(ctx, id, p, f) })
				continue
			}
			delivered := false
			s.withConn(m.Player, func(id string, p *peer) {
				delivered = true
				s.enqueue(ctx, id, p, f)
			})
			if !delivered {
				s.log.DebugContext(ctx, "dropping message for absent recipient",
					"piano.message_type", m.Type.String(), "piano.recipient", m.Player)
			}
		}
	}
}

// enqueue leaves a frame for a socket's writer, and closes a socket
// with no room for it. A socket that far behind is not waited for,
// and what it has missed is not kept for it: a page that reloads
// joins again and is sent the state whole.
func (s *Server) enqueue(ctx context.Context, id string, p *peer, f frame) {
	select {
	case p.queue <- f:
	default:
		s.log.WarnContext(ctx, "closing a socket that fell behind",
			"piano.message_type", f.kind.String(),
			"piano.recipient", id,
			"piano.queued_frames", len(p.queue))
		p.conn.Close()
	}
}

// writeLoop is a socket's one writer: the frames queued for it, and a
// ping every pingInterval. A write that fails is the connection's
// problem, not the server's: this used to panic, which meant one
// player closing a laptop lid at the wrong moment took the performance
// down. The socket is closed, which ends its read loop too.
func (s *Server) writeLoop(ctx context.Context, id string, p *peer) {
	ping := time.NewTicker(pingInterval)
	defer ping.Stop()
	for {
		var (
			kind string
			err  error
		)
		select {
		case <-p.gone:
			return
		case f := <-p.queue:
			kind = f.kind.String()
			p.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err = p.conn.WriteMessage(websocket.BinaryMessage, f.bytes)
		case <-ping.C:
			kind = "ping"
			err = p.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout))
		}
		if err != nil {
			s.log.WarnContext(ctx, "websocket write failed",
				"piano.message_type", kind,
				"piano.recipient", id,
				"error.message", err.Error(),
				"error.type", fmt.Sprintf("%T", err))
			p.conn.Close()
			return
		}
	}
}

// send puts a message on the game loop's inbox without blocking past
// ctx. Every path into the loop goes through here, so a wedged or
// stopped loop makes readers return instead of parking a goroutine on
// a channel nobody is reading.
func (s *Server) send(ctx context.Context, m *game.Message) bool {
	select {
	case s.inbox <- m:
		return true
	case <-ctx.Done():
		return false
	}
}
