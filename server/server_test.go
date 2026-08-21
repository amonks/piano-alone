package server

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"testing/synctest"

	"monks.co/piano-alone/game"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	s, err := New(t.Context(), Options{
		DBPath: filepath.Join(t.TempDir(), "performances.db"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// TestOperationClassification pins the split a host gates on. The
// routes belong to this package, not to the host, so a route added
// here later must land in the stricter bucket rather than arriving
// open — which is what the method check below is for.
func TestOperationClassification(t *testing.T) {
	for _, tc := range []struct {
		method, target string
		want           string
	}{
		{"GET", "/", OpRead},
		{"GET", "/download", OpRead},
		{"GET", "/main.wasm", OpRead},
		{"GET", "/latest-client", OpRead},
		{"GET", "/performances/featured", OpRead},
		{"GET", "/performances/scheduled", OpRead},
		{"GET", "/performances/abc/midi/x.midi", OpRead},
		{"GET", "/ws?fingerprint=abc", OpRead},

		// The disklavier receives a rendition and can drive nothing
		// (TestSocketRefusesControlFromNonConductor), so its socket is
		// a read; the conductor's is the performance's controls.
		{"GET", "/controller-ws", OpRead},
		{"GET", "/controller-ws?role=disklavier", OpRead},
		{"GET", "/controller-ws?role=conductor", OpConduct},

		{"POST", "/performances", OpConduct},
		{"POST", "/performances/abc/begin", OpConduct},
		{"POST", "/performances/abc/advance", OpConduct},
		{"POST", "/performances/abc/restart", OpConduct},
		{"DELETE", "/performances/abc", OpConduct},

		// Anything unrecognised and unsafe is a conduct.
		{"PUT", "/whatever", OpConduct},
		{"POST", "/some/route/added/later", OpConduct},
	} {
		req := httptest.NewRequest(tc.method, tc.target, nil)
		if got := Operation(req); got != tc.want {
			t.Errorf("Operation(%s %s) = %q, want %q", tc.method, tc.target, got, tc.want)
		}
	}
}

// The socket tests drive fromClient against a bare Server: they need
// an inbox and a logger and nothing else, and building one directly
// keeps them inside a synctest bubble, where "nothing was sent" is a
// settled fact rather than a sleep long enough to look like one.
func bareServer() (*Server, chan *game.Message) {
	inbox := make(chan *game.Message)
	return &Server{log: slog.New(slog.DiscardHandler), inbox: inbox}, inbox
}

// drain forwards everything the loop would have received, so a test
// can assert on what arrived. Without a drainer a refused message and
// a blocked send look alike — which is how the first version of this
// test hung instead of failing when the check was removed.
func drain(ctx context.Context, inbox <-chan *game.Message) chan *game.Message {
	got := make(chan *game.Message, 8)
	go func() {
		for {
			select {
			case m := <-inbox:
				got <- m
			case <-ctx.Done():
				return
			}
		}
	}()
	return got
}

// TestSocketRefusesControlFromNonConductor is the fix for the hole
// that gating only the HTTP verbs would have left: every socket
// message used to reach the game loop, so anyone who could open one —
// a player, or anyone claiming the disklavier's slot — could advance a
// phase, restart, or begin a performance.
func TestSocketRefusesControlFromNonConductor(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, inbox := bareServer()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		got := drain(ctx, inbox)

		for _, origin := range []string{"a-player-fingerprint", roleDisklavier} {
			for _, msgType := range []game.MessageType{
				game.MessageTypeAdvancePhase,
				game.MessageTypeRestart,
				game.MessageTypeBeginPerformance,
			} {
				if s.fromClient(ctx, origin, game.NewMessage(msgType, "", nil)) {
					t.Errorf("%s was allowed to send %s", origin, msgType)
				}
			}
		}

		synctest.Wait()
		select {
		case m := <-got:
			t.Errorf("a control message reached the game loop: %s", m.Type)
		default:
		}
	})
}

// TestSocketStampsTheSender: the Player field arrives from the
// internet. Believing it let one player submit another player's part,
// or complete their tutorial.
func TestSocketStampsTheSender(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, inbox := bareServer()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		got := drain(ctx, inbox)

		claimed := game.NewMessage(game.MessageTypeCompleteTutorial, "someone-elses-fingerprint", nil)
		if !s.fromClient(ctx, "my-fingerprint", claimed) {
			t.Fatal("message was refused")
		}

		synctest.Wait()
		select {
		case m := <-got:
			if m.Player != "my-fingerprint" {
				t.Errorf("sender = %q, want the connection's own fingerprint", m.Player)
			}
		default:
			t.Fatal("message never reached the loop")
		}
	})
}

// TestConductorMessagesReachTheLoop: the refusal above must not be so
// broad that the conductor cannot conduct.
func TestConductorMessagesReachTheLoop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, inbox := bareServer()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		got := drain(ctx, inbox)

		if !s.fromClient(ctx, roleConductor, game.NewMessage(game.MessageTypeAdvancePhase, "", nil)) {
			t.Fatal("the conductor was refused its own control message")
		}

		synctest.Wait()
		select {
		case m := <-got:
			if m.Type != game.MessageTypeAdvancePhase {
				t.Errorf("loop received %s", m.Type)
			}
		default:
			t.Error("the conductor's message never reached the loop")
		}
	})
}

// TestReservedFingerprintsAreRefused: the routing vocabulary and the
// fingerprint space are one namespace, so a player answering to
// "conductor" would be handed the controllers' messages.
func TestReservedFingerprintsAreRefused(t *testing.T) {
	s := newTestServer(t)
	for _, fingerprint := range []string{roleConductor, roleDisklavier, fingerprintControllers, fingerprintEveryone, ""} {
		req := httptest.NewRequest("GET", "/ws?fingerprint="+url.QueryEscape(fingerprint), nil)
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("fingerprint %q got %d, want 400", fingerprint, w.Code)
		}
	}
}

// TestMalformedInputIsAnswered, not panicked through. Both sockets and
// the schedule endpoint take gob from the public internet, and every
// decoder used to panic on a bad frame.
func TestMalformedInputIsAnswered(t *testing.T) {
	s := newTestServer(t)

	req := httptest.NewRequest("POST", "/performances", strings.NewReader("this is not gob"))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("malformed performance got %d, want 400", w.Code)
	}

	// A performance that decodes but names nothing cannot be stored
	// under an id, and would otherwise land as a row keyed on "".
	empty := (&game.Performance{}).Bytes()
	req = httptest.NewRequest("POST", "/performances", strings.NewReader(string(empty)))
	w = httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("performance with no id got %d, want 400", w.Code)
	}
}

// TestMissingPerformanceIs404: the id is in the URL, so a bad one is
// the caller's mistake. It used to be a 500 carrying the SQL error.
func TestMissingPerformanceIs404(t *testing.T) {
	s := newTestServer(t)
	for _, target := range []string{
		"/performances/nope/midi/x.midi",
	} {
		req := httptest.NewRequest("GET", target, nil)
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("GET %s got %d, want 404", target, w.Code)
		}
	}

	req := httptest.NewRequest("POST", "/performances/nope/begin", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("beginning a missing performance got %d, want 404", w.Code)
	}
}

// TestPagesCarryTheHostsBodyEnd: an app embedding this server puts the
// markup all its pages carry through BodyEnd, and a page that drops it
// is a page missing whatever the host puts on every page.
func TestPagesCarryTheHostsBodyEnd(t *testing.T) {
	s, err := New(t.Context(), Options{
		DBPath:  filepath.Join(t.TempDir(), "performances.db"),
		BodyEnd: func(*http.Request) template.HTML { return "<!--host-markup-->" },
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	for _, target := range []string{"/", "/download"} {
		req := httptest.NewRequest("GET", target, nil)
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("GET %s got %d", target, w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "<!--host-markup-->") {
			t.Errorf("GET %s does not carry the host's body-end markup", target)
		}
		if !strings.Contains(body, "</body>") {
			t.Fatalf("GET %s is not a whole document", target)
		}
		if strings.Index(body, "<!--host-markup-->") > strings.Index(body, "</body>") {
			t.Errorf("GET %s puts the host's markup after </body>", target)
		}
	}
}

// TestStandalonePagesCarryNothingExtra: BodyEnd is a seam, not a
// requirement. The piece's own pages are the piece's.
func TestStandalonePagesCarryNothingExtra(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET / got %d", w.Code)
	}
	if body := w.Body.String(); !strings.Contains(body, "Life Online") {
		t.Error("the page is not the piece's page")
	}
}

// TestNoClientBinaryIs404: the disklavier client is built and notarized
// out of band on a Mac, so a deployment without one says it has none
// rather than serving an empty 200 that a curl | chmod +x would leave
// an operator holding.
func TestNoClientBinaryIs404(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/latest-client", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("GET /latest-client with no binary got %d, want 404", w.Code)
	}
}

// TestUpgradeAcceptsAForwardedOrigin. Behind the proxy the request's
// Host is the backend's name while Origin is the name the browser
// typed, so gorilla's default same-origin check refuses every real
// connection — a failure that would only show up in production.
func TestUpgradeAcceptsAForwardedOrigin(t *testing.T) {
	req := httptest.NewRequest("GET", "/ws?fingerprint=abc", nil)
	req.Host = "monks-piano-alone-fly-ord"
	req.Header.Set("Origin", "https://piano.computer")
	if !upgrader.CheckOrigin(req) {
		t.Error("the upgrade check refuses a browser arriving through the proxy")
	}
}
