package server

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"monks.co/piano-alone/game"
)

// These run against a real listener on real time: a websocket
// handshake and a Range read are exactly the things a recorder cannot
// stand in for, and a real socket is the synctest carve-out.
//
// They exist because Serve is where the standalone server composes its
// own middleware, and nothing used to exercise it. Putting the
// library's gzip in front of the whole mux — rather than the two page
// routes it wrapped upstream — made `go run .`, the mode the README
// leads with, unable to accept a browser's websocket at all: the
// wrapper is not hijackable, so the handshake became a gzipped 500.

func runServer(t *testing.T) string {
	t.Helper()

	// A port from the OS, since these tests run concurrently with
	// whatever else the suite is doing.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()

	s, err := New(t.Context(), Options{
		DBPath: filepath.Join(t.TempDir(), "performances.db"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	go s.Go(ctx, addr)

	// Wait for the listener rather than sleeping at it. A liveness
	// backstop, not a timing assertion: the builder is slower than a
	// laptop and this is bounding a failure, not measuring one.
	deadline := time.Now().Add(30 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			conn.Close()
			return addr
		}
		if time.Now().After(deadline) {
			t.Fatalf("server never came up on %s: %v", addr, err)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestServeAcceptsABrowserWebsocket. Chrome and Firefox both send
// Accept-Encoding on the handshake, so a middleware that compresses
// whatever asked for gzip is a middleware that breaks every real
// browser connection while curl without the header still works.
func TestServeAcceptsABrowserWebsocket(t *testing.T) {
	addr := runServer(t)

	for _, tc := range []struct {
		name    string
		headers http.Header
	}{
		{"as a browser sends it", http.Header{"Accept-Encoding": {"gzip, deflate, br"}}},
		{"without Accept-Encoding", http.Header{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dialer := websocket.Dialer{HandshakeTimeout: 30 * time.Second}
			conn, resp, err := dialer.Dial("ws://"+addr+"/ws?fingerprint=a-player", tc.headers)
			if err != nil {
				status := "no response"
				if resp != nil {
					status = resp.Status
				}
				t.Fatalf("handshake failed (%s): %v", status, err)
			}
			defer conn.Close()

			// The socket is live both ways: a join is answered with the
			// state, which is what a player needs to render anything.
			join := game.JoinMessage{Fingerprint: "a-player", NoteCapacity: 1}
			if err := conn.WriteMessage(websocket.BinaryMessage,
				game.NewMessage(game.MessageTypeJoin, "", join.Bytes()).Bytes()); err != nil {
				t.Fatal(err)
			}
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			_, bs, err := conn.ReadMessage()
			if err != nil {
				t.Fatalf("no reply to a join: %v", err)
			}
			m, err := game.MessageFromBytes(bs)
			if err != nil {
				t.Fatal(err)
			}
			if m.Type != game.MessageTypeState {
				t.Errorf("reply to a join was %s, want the state", m.Type)
			}
		})
	}
}

// TestServeDoesNotCompressARangeRead. A range is computed against the
// identity encoding, so answering one with a gzip stream under a 206
// and the original Content-Range hands back bytes that are not the
// bytes asked for. The wasm player is the route this matters on: it is
// the one big binary the site serves.
func TestServeDoesNotCompressARangeRead(t *testing.T) {
	addr := runServer(t)

	req, err := http.NewRequestWithContext(t.Context(), "GET", "http://"+addr+"/main.wasm", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Range", "bytes=0-9")

	// The default transport would negotiate and transparently decode
	// its own gzip; this asks exactly what was sent.
	resp, err := (&http.Transport{DisableCompression: true}).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		t.Skip("the wasm player is not built in this checkout")
	}
	if enc := resp.Header.Get("Content-Encoding"); enc != "" {
		t.Errorf("Content-Encoding on a range read = %q, want none", enc)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != 10 {
		t.Errorf("bytes=0-9 returned %d bytes, want 10", len(body))
	}
}

// TestServeVariesOnAcceptEncoding: without it a cache can hand a
// compressed body to a client that did not ask for one.
func TestServeVariesOnAcceptEncoding(t *testing.T) {
	addr := runServer(t)

	req, err := http.NewRequestWithContext(t.Context(), "GET", "http://"+addr+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if vary := resp.Header.Get("Vary"); !strings.Contains(vary, "Accept-Encoding") {
		t.Errorf("Vary = %q, want it to name Accept-Encoding", vary)
	}
}

// socket is a client end of one of the server's sockets that keeps
// every frame it is sent, for a test to look through afterwards.
type socket struct {
	conn   *websocket.Conn
	frames [][]byte
}

func dial(t *testing.T, addr, target string) *socket {
	t.Helper()
	dialer := websocket.Dialer{HandshakeTimeout: 30 * time.Second}
	conn, resp, err := dialer.Dial("ws://"+addr+target, nil)
	if err != nil {
		status := "no response"
		if resp != nil {
			status = resp.Status
		}
		t.Fatalf("dialing %s (%s): %v", target, status, err)
	}
	t.Cleanup(func() { conn.Close() })
	return &socket{conn: conn}
}

func (s *socket) send(t *testing.T, msgType game.MessageType, data []byte) {
	t.Helper()
	s.write(t, game.NewMessage(msgType, "", data).Bytes())
}

func (s *socket) write(t *testing.T, frame []byte) {
	t.Helper()
	if err := s.conn.WriteMessage(websocket.BinaryMessage, frame); err != nil {
		t.Fatal(err)
	}
}

// await reads until a message of the given type arrives, and returns
// it. The deadline is a liveness backstop, not a timing assertion.
func (s *socket) await(t *testing.T, msgType game.MessageType) *game.Message {
	t.Helper()
	s.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	for {
		_, bs, err := s.conn.ReadMessage()
		if err != nil {
			t.Fatalf("waiting for %s: %v", msgType, err)
		}
		s.frames = append(s.frames, bs)
		m, err := game.MessageFromBytes(bs)
		if err != nil {
			t.Fatal(err)
		}
		if m.Type == msgType {
			return m
		}
	}
}

func post(t *testing.T, addr, target string, body []byte) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), "POST", "http://"+addr+target, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("POST %s: %s %s", target, resp.Status, b)
	}
}

// A fingerprint is a player's whole credential: a socket that presents
// one acts as that player. So no socket is ever told one — not a
// player's, which is sent the whole state, and not the disklavier's,
// which anyone may open and which hears of every player who joins,
// starts the tutorial, or submits. Every state message used to carry
// every player's fingerprint, so any player could take any other's
// place. Players are known everywhere by PlayerID instead, a digest
// the fingerprint cannot be recovered from, and a page finds itself in
// the state by its own.
func TestNoSocketIsToldAFingerprint(t *testing.T) {
	addr := runServer(t)
	conf := &game.Configuration{PerformanceID: "tonight", Title: "Tonight", Composer: "Somebody", Score: partOf(t, 60, 0)}
	post(t, addr, "/performances", (&game.Performance{Configuration: conf}).Bytes())

	disklavier := dial(t, addr, "/controller-ws?role=disklavier")
	disklavier.await(t, game.MessageTypeState)

	fingerprints := []string{strings.Repeat("a1", 128), strings.Repeat("b2", 128)}
	var players []*socket
	for _, f := range fingerprints {
		p := dial(t, addr, "/ws?fingerprint="+f)
		p.send(t, game.MessageTypeJoin, game.JoinMessage{Fingerprint: f, NoteCapacity: 3}.Bytes())
		p.await(t, game.MessageTypeState)
		p.send(t, game.MessageTypeStartTutorial, nil)
		players = append(players, p)
	}
	for range players {
		disklavier.await(t, game.MessageTypeStartTutorial)
	}

	post(t, addr, "/performances/tonight/begin", nil)
	var states []*game.State
	for _, p := range players {
		m := p.await(t, game.MessageTypeState)
		state, err := game.StateFromBytes(m.Data)
		if err != nil {
			t.Fatal(err)
		}
		states = append(states, state)
	}
	post(t, addr, "/performances/tonight/advance", nil)
	for _, p := range players {
		p.await(t, game.MessageTypeAssignment)
	}
	disklavier.await(t, game.MessageTypeBroadcastPhase) // Lobby
	disklavier.await(t, game.MessageTypeBroadcastPhase) // Hero

	for i, f := range fingerprints {
		if me := states[i].Players[game.PlayerID(f)]; me == nil || me.Fingerprint != game.PlayerID(f) {
			t.Errorf("player %d is not in its own state by its id: %+v", i, states[i].Players)
		}
	}
	for name, s := range map[string]*socket{"disklavier": disklavier, "player 0": players[0], "player 1": players[1]} {
		for _, frame := range s.frames {
			for i, f := range fingerprints {
				if bytes.Contains(frame, []byte(f)) {
					m, _ := game.MessageFromBytes(frame)
					t.Errorf("the %s socket was told player %d's fingerprint in a %s", name, i, m.Type)
				}
			}
		}
	}
}

// A socket message is bounded by the largest one the protocol has: a
// player's part at its cap, from the longest fingerprint, read whole;
// one past that closes the socket rather than being read into memory.
// The bound used to be 8 MiB on every socket — a hundred times any
// part, on sockets anyone may open.
func TestSocketMessagesAreBoundedByTheProtocol(t *testing.T) {
	addr := runServer(t)

	f := strings.Repeat("ab", maxFingerprintBytes/2)
	largest := game.NewMessage(game.MessageTypeSubmitPartialTrack, f, make([]byte, maxPartBytes)).Bytes()
	if len(largest) > maxPlayerMessageBytes {
		t.Fatalf("the largest part a player sends is %d bytes on the wire, over the %d-byte limit", len(largest), maxPlayerMessageBytes)
	}
	p := dial(t, addr, "/ws?fingerprint="+f)
	p.write(t, largest)
	p.send(t, game.MessageTypeJoin, game.JoinMessage{NoteCapacity: 1}.Bytes())
	p.await(t, game.MessageTypeState)

	p.write(t, game.NewMessage(game.MessageTypeSubmitPartialTrack, f, make([]byte, maxPlayerMessageBytes)).Bytes())
	p.send(t, game.MessageTypeJoin, game.JoinMessage{NoteCapacity: 1}.Bytes())
	p.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	_, bs, err := p.conn.ReadMessage()
	if err == nil {
		m, _ := game.MessageFromBytes(bs)
		t.Fatalf("a message past the limit was read, and the socket answered with a %s", m.Type)
	}
	if !websocket.IsCloseError(err, websocket.CloseMessageTooBig) {
		t.Errorf("a message past the limit: %v, want the socket closed as too big", err)
	}

	// The conductor's socket can begin a performance, whose
	// configuration carries a whole score, so it takes what scheduling
	// one does.
	c := dial(t, addr, "/controller-ws?role=conductor")
	c.await(t, game.MessageTypeState)
	conf := &game.Configuration{PerformanceID: "big", Title: "Big", Composer: "Somebody", Score: make([]byte, 1<<20)}
	c.send(t, game.MessageTypeBeginPerformance, conf.Bytes())
	m := c.await(t, game.MessageTypeState)
	if state, err := game.StateFromBytes(m.Data); err != nil || state.Configuration == nil || state.Configuration.PerformanceID != "big" {
		t.Errorf("the conductor's begin with a 1 MiB score: state %+v, %v", state, err)
	}
}
