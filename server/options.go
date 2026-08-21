package server

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"monks.co/piano-alone/data"
)

// Options configure a Server. Every field has a working default, so a
// caller who wants the piece as it has always run passes Options{} and
// nothing else. The fields exist for a host that already owns the
// things a server needs — a database it replicates, a logger it ships,
// markup its other pages carry — and would rather lend them than have
// this package invent its own.
type Options struct {
	// DB is the store. Nil means open DBPath, which is what the
	// standalone command does. A handle passed here is the caller's:
	// the server neither closes it nor changes its settings, and the
	// caller is expected to have migrated it (db.Migrations,
	// db.Baseline).
	DB *sql.DB

	// DBPath is where to open a database when DB is nil. Empty means
	// "piano-alone.db" in the working directory.
	DBPath string

	// Logger receives the server's structured record of what it does:
	// phase changes, joins, submissions, the rendition write. Nil
	// means slog.Default().
	Logger *slog.Logger

	// BodyEnd renders host-owned markup at the end of every page's
	// body. Nil renders nothing, which is what the piece's own pages
	// have always done — it exists for a host whose pages all carry
	// something, not because these pages need chrome.
	BodyEnd func(*http.Request) template.HTML

	// ClientPath is the disklavier client binary served at
	// /latest-client. Empty, or a path that does not exist, answers
	// 404: the binary is a signed universal Mac build made out of
	// band, so a server that has not been given one should say it has
	// none rather than pretend.
	ClientPath string
}

// Operations. A host gates this server from outside — the package
// authorizes nothing itself — and Operation is how it tells the two
// kinds of request apart without enumerating routes it does not own.
const (
	// OpRead is everything the audience touches: the page, the player,
	// the finished MIDI files, and the disklavier's socket, which
	// receives a rendition and cannot drive anything (§ Server.fromClient).
	OpRead = "read"

	// OpConduct drives the performance: scheduling, beginning,
	// advancing, restarting, deleting, and the conductor's socket.
	OpConduct = "conduct"
)

// Operation classifies a request. Anything that is not plainly a read
// is a conduct, so a route added here later arrives in the stricter
// bucket rather than arriving ungated — the classification is a
// property of the request, and a host that wrote its policy against
// today's routes should not silently start admitting tomorrow's.
func Operation(r *http.Request) string {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return OpConduct
	}
	if strings.TrimSuffix(r.URL.Path, "/") == strings.TrimPrefix(data.PathControllerWS, "GET ") &&
		r.URL.Query().Get("role") == roleConductor {
		return OpConduct
	}
	return OpRead
}
