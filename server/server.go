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
	"gitlab.com/gomidi/midi/v2/smf"

	"monks.co/piano-alone/assets"
	"monks.co/piano-alone/db"
	"monks.co/piano-alone/game"
)

// The two controller roles, which are also the fingerprints the game
// loop addresses them by.
const (
	roleConductor  = "conductor"
	roleDisklavier = "disklavier"

	// fingerprintControllers addresses both controllers at once.
	fingerprintControllers = "controllers"

	// fingerprintEveryone broadcasts.
	fingerprintEveryone = "*"
)

const defaultDBPath = "piano-alone.db"

type Server struct {
	db     *db.DB
	ownsDB bool
	log    *slog.Logger

	bodyEnd    func(*http.Request) template.HTML
	clientPath string

	mux http.Handler

	disklavierConn *websocket.Conn
	conductorConn  *websocket.Conn
	conns          map[string]*websocket.Conn
	connMu         sync.RWMutex

	inbox  chan *game.Message
	outbox chan *game.Message

	// The performance. These belong to the game loop and are touched
	// from nowhere else, which is why they carry no lock: play handles
	// one message at a time. All of it is in memory — only the
	// finished rendition is written down — so a restart mid-performance
	// loses everything in flight.
	state                  *game.State
	partials               map[string]*smf.SMF
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
		conns:      map[string]*websocket.Conn{},
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

// deliver fans the game loop's outbound messages out to the sockets.
// A write that fails is the connection's problem, not the server's:
// this used to panic, which meant one player closing a laptop lid at
// the wrong moment took the performance down.
func (s *Server) deliver(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-s.outbox:
			if m.Player == fingerprintEveryone {
				s.eachConn(func(fingerprint string, sock *websocket.Conn) {
					s.write(ctx, fingerprint, sock, m)
				})
				continue
			}
			delivered := false
			s.withConn(m.Player, func(fingerprint string, conn *websocket.Conn) {
				delivered = true
				s.write(ctx, fingerprint, conn, m)
			})
			if !delivered {
				s.log.DebugContext(ctx, "dropping message for absent recipient",
					"piano.message_type", m.Type.String(), "piano.recipient", m.Player)
			}
		}
	}
}

// writeTimeout bounds one socket write. Delivery is a single
// goroutine holding the registry's read lock, so a phone on bad venue
// wifi that stops reading would otherwise stall every other player,
// block a connect or disconnect behind the write lock, and then block
// the game loop on its own send. A write that fails is the
// connection's problem; without a deadline, a write that *hangs* is
// everyone's.
const writeTimeout = 10 * time.Second

func (s *Server) write(ctx context.Context, fingerprint string, conn *websocket.Conn, m *game.Message) {
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err := conn.WriteMessage(websocket.BinaryMessage, m.Bytes()); err != nil {
		s.log.WarnContext(ctx, "websocket write failed",
			"piano.message_type", m.Type.String(),
			"piano.recipient", fingerprint,
			"error.message", err.Error(),
			"error.type", fmt.Sprintf("%T", err))
		conn.Close()
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
