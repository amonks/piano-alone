package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/gorilla/websocket"

	"monks.co/piano-alone/game"
)

// These tests open the server's sockets over in-memory pipes, inside a
// synctest bubble. A pipe has no buffer, so a client that stops
// reading stalls the server's write at its first byte — which a real
// socket does only once the kernel's buffers have filled — and every
// timeout is the bubble's clock, so "held up for ten seconds" and "not
// held up at all" are exact rather than a race against a sleep.

// socketServer is a server with its delivery running and no game
// loop: a test says what the loop would by sending on the outbox, and
// reads what the sockets told the loop from the channel returned.
func socketServer(t *testing.T) (*Server, <-chan *game.Message) {
	t.Helper()
	s, err := New(t.Context(), Options{
		DBPath: filepath.Join(t.TempDir(), "performances.db"),
		Logger: slog.New(slog.DiscardHandler),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	go s.deliver(ctx)
	heard := make(chan *game.Message, 1<<16)
	go func() {
		for {
			select {
			case m := <-s.inbox:
				heard <- m
			case <-ctx.Done():
				return
			}
		}
	}()
	return s, heard
}

// pipeResponse is the ResponseWriter of a request that arrived on a
// pipe: it records a refusal, or hands the pipe to an upgrade.
type pipeResponse struct {
	*httptest.ResponseRecorder
	conn     net.Conn
	rw       *bufio.ReadWriter
	hijacked bool
}

func (w *pipeResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.hijacked = true
	return w.conn, w.rw, nil
}

// open dials one of the server's sockets. The request is read off the
// server's end of a pipe and handed to the handler, as net/http would.
func open(t *testing.T, s *Server, target string) (*websocket.Conn, *http.Response, error) {
	t.Helper()
	ctx := t.Context()
	dialer := websocket.Dialer{NetDialContext: func(context.Context, string, string) (net.Conn, error) {
		client, server := net.Pipe()
		go func() {
			defer server.Close()
			br := bufio.NewReader(server)
			req, err := http.ReadRequest(br)
			if err != nil {
				return
			}
			w := &pipeResponse{
				ResponseRecorder: httptest.NewRecorder(),
				conn:             server,
				rw:               bufio.NewReadWriter(br, bufio.NewWriter(server)),
			}
			s.Handler().ServeHTTP(w, req.WithContext(ctx))
			if !w.hijacked {
				w.Result().Write(server)
			}
		}()
		return client, nil
	}}
	conn, resp, err := dialer.Dial("ws://piano.test"+target, nil)
	if err == nil {
		t.Cleanup(func() { conn.Close() })
	}
	return conn, resp, err
}

func mustOpen(t *testing.T, s *Server, target string) *websocket.Conn {
	t.Helper()
	conn, _, err := open(t, s, target)
	if err != nil {
		t.Fatalf("opening %s: %v", target, err)
	}
	return conn
}

// listener is a client that reads whatever it is sent, as a page does,
// which is also what answers the server's pings.
type listener struct {
	got    chan *game.Message
	closed chan error
}

func listen(conn *websocket.Conn) *listener {
	l := &listener{got: make(chan *game.Message, 1<<16), closed: make(chan error, 1)}
	go func() {
		for {
			_, bs, err := conn.ReadMessage()
			if err != nil {
				l.closed <- err
				return
			}
			if m, err := game.MessageFromBytes(bs); err == nil {
				l.got <- m
			}
		}
	}()
	return l
}

// isClosed reports whether the server has closed the listener's
// socket, and how.
func (l *listener) isClosed() (error, bool) {
	select {
	case err := <-l.closed:
		l.closed <- err
		return err, true
	default:
		return nil, false
	}
}

func phase() *game.Message {
	return game.NewMessage(game.MessageTypeBroadcastPhase, toEveryone, game.GamePhaseLobby.Bytes())
}

// has reports whether the server still holds a socket for a player.
func (s *Server) has(fingerprint string) bool {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	_, ok := s.conns[game.PlayerID(fingerprint)]
	return ok
}

// A socket that stops reading holds up nobody else. Delivery used to
// write to each socket in turn, under the registry's lock, waiting up
// to ten seconds on each: so every stalled socket cost every broadcast
// ten seconds, and the game loop, which hands delivery one message at
// a time, stopped with it. Anyone may open a player's socket and not
// read from it.
func TestASocketThatStopsReadingHoldsUpNoOne(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, _ := socketServer(t)
		mustOpen(t, s, "/ws?fingerprint=stalled") // and never read
		page := listen(mustOpen(t, s, "/ws?fingerprint=page"))
		disklavier := listen(mustOpen(t, s, "/controller-ws?role=disklavier"))
		synctest.Wait()

		start := time.Now()
		for range 3 {
			s.outbox <- phase()
		}
		synctest.Wait()
		if waited := time.Since(start); waited != 0 {
			t.Errorf("three broadcasts took %v beside a socket that is not reading", waited)
		}
		for name, l := range map[string]*listener{"the page": page, "the disklavier": disklavier} {
			if len(l.got) != 3 {
				t.Errorf("%s was sent %d of 3 broadcasts", name, len(l.got))
			}
		}

		// The stalled socket's own write times out, and it is closed.
		if !s.has("stalled") {
			t.Fatal("the stalled socket was closed before its write could time out")
		}
		time.Sleep(writeTimeout)
		synctest.Wait()
		if s.has("stalled") {
			t.Errorf("a socket that read nothing for %v is still open", writeTimeout)
		}
		if _, closed := page.isClosed(); closed {
			t.Error("the page's socket was closed too")
		}
	})
}

// A socket has a few frames' room to fall behind in, and one that
// falls further is closed rather than waited for or kept in memory: a
// page that far behind has missed what it needed, and a reload joins
// again and is sent the state whole.
func TestASocketThatFallsBehindIsClosed(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, _ := socketServer(t)
		mustOpen(t, s, "/ws?fingerprint=stalled")
		page := listen(mustOpen(t, s, "/ws?fingerprint=page"))
		synctest.Wait()

		// One frame is with the stalled socket's writer, and the queue
		// holds the rest.
		const within, past = 1 + playerQueue, 4
		start := time.Now()
		for range within {
			s.outbox <- phase()
		}
		synctest.Wait()
		if !s.has("stalled") {
			t.Fatalf("a socket %d frames behind was closed: its queue holds %d", within, playerQueue)
		}
		for range past {
			s.outbox <- phase()
		}
		synctest.Wait()
		if s.has("stalled") {
			t.Errorf("a socket %d frames behind is still open", within+past)
		}
		if waited := time.Since(start); waited != 0 {
			t.Errorf("closing it took %v", waited)
		}
		if len(page.got) != within+past {
			t.Errorf("the page was sent %d of %d broadcasts", len(page.got), within+past)
		}
	})
}

// A peer that is gone — a phone out of range, a laptop lid closed —
// says nothing, and its socket used to stay open until something was
// written to it and failed, which for a quiet socket is never. Each
// socket is pinged, and one that answers nothing for readTimeout is
// closed and its player marked gone. A page answers pings without
// knowing it, and stays.
func TestASocketThatAnswersNothingIsClosed(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, heard := socketServer(t)
		gone := mustOpen(t, s, "/ws?fingerprint=gone")
		gone.SetPingHandler(func(string) error { return nil }) // reads its pings and answers none
		silent := listen(gone)
		page := listen(mustOpen(t, s, "/ws?fingerprint=page"))

		time.Sleep(readTimeout + pingInterval)
		synctest.Wait()
		if _, closed := silent.isClosed(); !closed || s.has("gone") {
			t.Errorf("a socket that answered no ping for %v is still open", readTimeout+pingInterval)
		}
		left := false
		for len(heard) > 0 {
			if m := <-heard; m.Type == game.MessageTypeLeave && m.Player == game.PlayerID("gone") {
				left = true
			}
		}
		if !left {
			t.Error("the game loop was not told the player left")
		}
		if err, closed := page.isClosed(); closed || !s.has("page") {
			t.Errorf("a page that answers its pings was closed: %v", err)
		}
	})
}

// The players' sockets open at once are capped: each is a goroutine,
// buffers, and a queue, and the socket is anonymous. The one past the
// cap is refused before it is upgraded, and a socket closing makes
// room. The controllers' sockets are not counted, so a full house
// cannot keep the conductor or the disklavier out.
func TestPlayersSocketsAreCapped(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, _ := socketServer(t)
		var first *websocket.Conn
		for i := range maxPlayerSockets {
			conn, resp, err := open(t, s, fmt.Sprintf("/ws?fingerprint=player%d", i))
			if err != nil {
				status := "no response"
				if resp != nil {
					status = resp.Status
				}
				t.Fatalf("socket %d of %d was refused (%s): %v", i+1, maxPlayerSockets, status, err)
			}
			if i == 0 {
				first = conn
			}
		}

		_, resp, err := open(t, s, "/ws?fingerprint=late")
		if !errors.Is(err, websocket.ErrBadHandshake) || resp == nil || resp.StatusCode != http.StatusServiceUnavailable {
			status := "no response"
			if resp != nil {
				status = resp.Status
			}
			t.Fatalf("socket %d: %s, %v, want it refused with a 503", maxPlayerSockets+1, status, err)
		}
		for _, role := range []string{roleConductor, roleDisklavier} {
			if _, _, err := open(t, s, "/controller-ws?role="+role); err != nil {
				t.Errorf("the %s was refused at a full house: %v", role, err)
			}
		}

		first.Close()
		synctest.Wait()
		if _, _, err := open(t, s, "/ws?fingerprint=late"); err != nil {
			t.Errorf("a socket was refused after another closed: %v", err)
		}
	})
}

// A socket's place — a player's id, or a controller's role — holds one
// socket, and a newer one closes the one it replaces. The older socket
// was already sent nothing, since messages go to whoever holds the
// place; left open it was a socket nobody counted, and the disklavier's
// place is public, so opening it over and over left any number of
// them. Its player has not left: the newer socket is theirs.
func TestANewerSocketClosesTheOneItReplaces(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, heard := socketServer(t)
		for _, target := range []string{"/controller-ws?role=disklavier", "/controller-ws?role=conductor", "/ws?fingerprint=page"} {
			older := listen(mustOpen(t, s, target))
			synctest.Wait()
			newer := listen(mustOpen(t, s, target))
			synctest.Wait()

			err, closed := older.isClosed()
			if !closed {
				t.Errorf("%s: the older socket is still open", target)
			} else if !websocket.IsCloseError(err, closeReplaced) {
				t.Errorf("%s: the older socket was closed with %v, want it told it was replaced", target, err)
			}
			if err, closed := newer.isClosed(); closed {
				t.Errorf("%s: the newer socket was closed: %v", target, err)
			}
		}
		s.outbox <- phase()
		synctest.Wait()
		for len(heard) > 0 {
			switch m := <-heard; m.Type {
			case game.MessageTypeLeave, game.MessageTypeDisklavierDisconnected, game.MessageTypeConductorDisconnected:
				t.Errorf("the game loop was told of a %s for a place that is still held", m.Type)
			}
		}
	})
}

// One socket's messages are read at a pace: a few at once, and then
// one every inboundGap. A page sends four in a performance. Unpaced, a
// socket could have the loop tell the controllers of it thousands of
// times a second, faster than a controller's own socket drains — and a
// socket that falls behind is closed.
func TestASocketIsReadAtAPace(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, heard := socketServer(t)
		conn := mustOpen(t, s, "/ws?fingerprint=eager")
		const sent = inboundBurst + 5
		go func() {
			for range sent {
				conn.WriteMessage(websocket.BinaryMessage, game.NewMessage(game.MessageTypeStartTutorial, "", nil).Bytes())
			}
		}()
		synctest.Wait()
		if len(heard) != inboundBurst {
			t.Errorf("%d messages were read at once, want %d", len(heard), inboundBurst)
		}
		time.Sleep(5 * inboundGap)
		synctest.Wait()
		if len(heard) != sent {
			t.Errorf("%d messages were read after %v, want all %d", len(heard), 5*inboundGap, sent)
		}
	})
}
