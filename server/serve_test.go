package server

import (
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
