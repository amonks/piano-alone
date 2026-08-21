# Life Online (Piano Telephone)

An audience, each given a handful of a score's keys, plays a piano
piece together on their phones. The parts are collected, merged into
one MIDI file, and performed on a disklavier.

Written for [_Piano, Alone in a Room: Translating/Transcoding_][call],
and performed twice in 2024. The proposal — what the piece is about, and
why it is a circle — is in [proposal/piano-alone.md](proposal/piano-alone.md).

[call]: https://docs.google.com/document/d/13bX58CqNuS43leFi2C0TOePTz89e5GglPmOBu-RSdYY/edit

Live at [piano.computer](https://piano.computer).

## What runs

Four programs out of one module:

| Program | What it is |
| --- | --- |
| `monks.co/piano-alone` | the server: the website, the sockets, and the game loop |
| `./gameplayer` | the player, compiled to wasm and embedded in the server |
| `./gamecontroller` | a terminal program with two roles — `conductor` drives the performance, `disklavier` plays the finished MIDI out a port |
| `./cmd/*` | MIDI utilities used while making the piece; not part of a performance |

## Running it

```sh
go run . -addr 0.0.0.0:8080 -db piano-alone.db
```

That is the whole thing: it creates the database, seeds one unplayed
performance, and serves the site. Then, from a machine with a MIDI
cable in it:

```sh
go run ./gamecontroller -role disklavier -baseURL http://localhost:8080
go run ./gamecontroller -role conductor  -baseURL http://localhost:8080
```

The conductor picks a performance and presses Advance to move through
the phases; nothing here is on a timer.

| Phase | What is happening |
| --- | --- |
| Lobby | players open the page and are counted |
| Hero | everyone plays their own handful of notes against a scrolling score |
| Processing | the parts are merged into one MIDI file |
| Playback | the disklavier plays it |
| Done | the page comes back, with the MIDI file to download |

All of that is in memory except the finished rendition, so restarting
the server mid-performance loses the performance.

**Run standalone, every conductor verb is open to whoever can reach
it** — scheduling, advancing, deleting. That is the right default for a
laptop in a room and the wrong one for the public internet; see
[Embedding](#embedding).

## Building

Needs Go and npm. [`run`](https://github.com/amonks/run) drives the
task graph in `tasks.toml`:

```sh
run generate   # templ, the stylesheet, the wasm player
run dev        # the above, then the server, rebuilding as you edit
```

Or by hand — note that `templ` and `stringer` are not dependencies of
this module, so outside the monorepo that owns them you need
`go run github.com/a-h/templ/cmd/templ@latest` and
`go run golang.org/x/tools/cmd/stringer@latest` in their place:

```sh
go tool templ generate
npm install && ./node_modules/.bin/postcss -o assets/files/style.css css/style.css
GOOS=js GOARCH=wasm go build -buildvcs=false -ldflags="-s -w" -o assets/files/main.wasm ./gameplayer
cp -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" assets/files/wasm_exec.js
```

The built assets are embedded, so the server binary carries the client
it serves. They are also gitignored, and the embed is of the directory
rather than of each file, so a fresh checkout compiles before it has
generated anything — a missing asset surfaces as an error on the page
that wanted it rather than as a build failure.

Two generated things are committed instead, because reproducing them
costs more than keeping them:

- `c2d/c2d_wasm.go`, a typed wrapper over `CanvasRenderingContext2D`
  generated from TypeScript's DOM library by
  `generate_canvas_api_wrapper/`. Regenerate it by hand (it needs yarn
  and tsc) when the DOM changes, which is roughly never.
- the `*_string.go` files, from `stringer`. Regenerate with
  `go generate ./...` after changing one of the enums.

### The disklavier client binary

`GET /latest-client` serves a signed, notarized universal Mac binary of
`./gamecontroller`, so an operator can `curl | chmod +x` it at a venue.
Building one needs a Mac, a Developer ID certificate, and
`$APP_SPECIFIC_PASSWORD`: `run build-client`. Hand the result to a
running server with `-client`, and a server without one answers 404
rather than serving an empty 200.

## Embedding

The piece is importable. `server.New` takes what a host already owns
rather than inventing its own; every field has a working default, so
`Options{}` is the standalone server.

```go
type Options struct {
	DB         *sql.DB                              // nil: open DBPath
	DBPath     string                               // "": ./piano-alone.db
	Logger     *slog.Logger                         // nil: slog.Default()
	BodyEnd    func(*http.Request) template.HTML    // nil: the piece's own pages
	ClientPath string                               // "": /latest-client answers 404
}

func New(context.Context, Options) (*Server, error)
func (*Server) Handler() http.Handler   // the whole site, rooted at "/"
func (*Server) Run(context.Context) error  // the game loop
func (*Server) Serve(context.Context, addr string) error
func (*Server) Go(context.Context, addr string) error  // Run + Serve: the standalone server
```

A host that owns its own listener runs `Run` and mounts `Handler`.
Schema for a host with its own migration runner: `db.Migrations`,
`db.MigrationsDir`, `db.Baseline`.

### Authorization

There is none in here, deliberately. `Handler` returns every route on
one mux and `Operation` says which of two kinds a request is, so a host
can put one decision in front of the whole thing:

```go
const (
	OpRead    = "read"     // the page, the player, the sockets, finished MIDI
	OpConduct = "conduct"  // schedule, begin, advance, restart, delete, the conductor's socket
)

func Operation(*http.Request) string
```

Anything not plainly a read classifies as a conduct, so a route added
here later arrives in the stricter bucket rather than arriving open.

Inside, the server enforces the half a host cannot see: every socket
message is stamped with the connection it came from, and control
messages are refused from anything but the conductor's connection. So
the disklavier's socket really is receive-only, and a player speaks
only for itself.

### A note on the wire

Messages are `encoding/gob`, which is a choice that has aged. It is
what the two performances ran on and what the recorded renditions were
made with, so it stays. The decoders return errors rather than
panicking, which they did not originally: both sockets and the
schedule endpoint take bytes from the public internet.

## License

[Blue Oak Model License 1.0.0](LICENSE.md).
