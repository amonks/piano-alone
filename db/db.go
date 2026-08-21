// Package db is piano-alone's store: one table of performances, over
// database/sql.
//
// The handle is the caller's. A host that already manages its own
// SQLite — connection pool, WAL settings, replication — hands one to
// New and keeps ownership of it; the standalone command uses Open to
// get one of its own. Nothing here opens, closes, or re-pragmas a
// handle it was given, because a replicated database's settings belong
// to whoever set up the replication.
package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"monks.co/piano-alone/game"

	_ "modernc.org/sqlite"
)

// Migrations is the schema, for a host that runs migrations itself:
// pass it, MigrationsDir, and Baseline to whatever runner the host
// already has. Baseline names the files an existing database records
// rather than executes — the deployed database predates this runner
// and already holds the table 001 would create.
//
//go:embed migrations
var Migrations embed.FS

const MigrationsDir = "migrations"

var Baseline = []string{"001_initial.sql"}

// dateLayout is how a performance date is stored: the format the rows
// written before this package already hold. New rows are written in
// UTC so that they sort against each other lexically, which is what
// the featured listing's ORDER BY does; the two pre-existing rows keep
// the offsets they were written with.
const dateLayout = "2006-01-02 15:04:05-07:00"

// readLayouts are tried in order when reading a date back. The first
// is what this package writes; the rest are what GORM wrote, whose
// fractional seconds and zone spelling varied with the value.
var readLayouts = []string{
	dateLayout,
	"2006-01-02 15:04:05.999999999-07:00",
	time.RFC3339Nano,
	"2006-01-02 15:04:05",
}

// ErrNotFound is returned for a performance id that is not in the
// table, so a handler can answer 404 rather than 500.
var ErrNotFound = errors.New("no such performance")

type DB struct{ sql *sql.DB }

// New wraps an already-open, already-migrated database. The caller
// keeps ownership of the handle: New starts nothing and DB closes
// nothing.
func New(handle *sql.DB) *DB { return &DB{sql: handle} }

// Open opens SQLite at path with the pure-Go driver and runs the
// migrations, for a caller with no database of its own. A host with
// its own SQLite should use New instead.
func Open(ctx context.Context, path string) (*DB, error) {
	handle, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := Migrate(ctx, handle); err != nil {
		handle.Close()
		return nil, err
	}
	return New(handle), nil
}

// Close closes the underlying handle. Only its owner should call it:
// a DB built by New holds a handle someone else opened.
func (db *DB) Close() error { return db.sql.Close() }

func parseDate(s string) (time.Time, error) {
	for _, layout := range readLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date %q", s)
}

const columns = `id, title, composer, score, date, is_featured, is_complete, rendition, player_count`

// listColumns omits the blobs. The listings render titles and dates,
// and a score is 15 KB per row while a rendition is the whole
// performance — reading them to throw them away is what the original
// SELECT lists avoided too.
const listColumns = `id, title, composer, date, is_featured, is_complete, player_count`

func scanList(rows *sql.Rows) ([]*game.Performance, error) {
	var out []*game.Performance
	for rows.Next() {
		var (
			p    game.Performance
			conf game.Configuration
			date string
		)
		if err := rows.Scan(&conf.PerformanceID, &conf.Title, &conf.Composer,
			&date, &p.IsFeatured, &p.IsComplete, &p.PlayerCount); err != nil {
			return nil, err
		}
		t, err := parseDate(date)
		if err != nil {
			return nil, err
		}
		p.Date, p.Configuration = t, &conf
		out = append(out, &p)
	}
	return out, rows.Err()
}

func (db *DB) GetPerformance(ctx context.Context, id string) (*game.Performance, error) {
	var (
		p    game.Performance
		conf game.Configuration
		date string
	)
	err := db.sql.QueryRowContext(ctx,
		`SELECT `+columns+` FROM performances WHERE id = ?`, id).
		Scan(&conf.PerformanceID, &conf.Title, &conf.Composer, &conf.Score,
			&date, &p.IsFeatured, &p.IsComplete, &p.Rendition, &p.PlayerCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	} else if err != nil {
		return nil, err
	}
	t, err := parseDate(date)
	if err != nil {
		return nil, err
	}
	p.Date, p.Configuration = t, &conf
	return &p, nil
}

func (db *DB) GetScheduledPerformances(ctx context.Context) ([]*game.Performance, error) {
	rows, err := db.sql.QueryContext(ctx,
		`SELECT `+listColumns+` FROM performances WHERE is_complete = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanList(rows)
}

func (db *DB) GetFeaturedPerformances(ctx context.Context) ([]*game.Performance, error) {
	rows, err := db.sql.QueryContext(ctx,
		`SELECT `+listColumns+` FROM performances WHERE is_featured = 1 ORDER BY date ASC, title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanList(rows)
}

func (db *DB) SchedulePerformance(ctx context.Context, p *game.Performance) error {
	if p.Configuration == nil {
		return errors.New("performance has no configuration")
	}
	_, err := db.sql.ExecContext(ctx,
		`INSERT INTO performances (`+columns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Configuration.PerformanceID, p.Configuration.Title, p.Configuration.Composer,
		p.Configuration.Score, p.Date.UTC().Format(dateLayout),
		p.IsFeatured, p.IsComplete, p.Rendition, p.PlayerCount)
	return err
}

// StartPerformance records a performance beginning now, for a caller
// driving the piece without having scheduled it first.
func (db *DB) StartPerformance(ctx context.Context, conf *game.Configuration) error {
	return db.SchedulePerformance(ctx, &game.Performance{Configuration: conf, Date: time.Now()})
}

func (db *DB) DeletePerformance(ctx context.Context, id string) error {
	_, err := db.sql.ExecContext(ctx, `DELETE FROM performances WHERE id = ?`, id)
	return err
}

// SaveRendition stores the merged MIDI of a finished performance. This
// is the piece's one irreplaceable write: the rendition is what an
// audience actually played, and it exists nowhere else.
func (db *DB) SaveRendition(ctx context.Context, id string, playerCount int, rendition []byte) error {
	res, err := db.sql.ExecContext(ctx,
		`UPDATE performances SET is_complete = 1, rendition = ?, player_count = ? WHERE id = ?`,
		rendition, playerCount, id)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err == nil && n == 0 {
		return fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return nil
}

func (db *DB) Feature(ctx context.Context, id string) error {
	_, err := db.sql.ExecContext(ctx, `UPDATE performances SET is_featured = 1 WHERE id = ?`, id)
	return err
}

func (db *DB) GetMIDIFile(ctx context.Context, id string) ([]byte, error) {
	var rendition []byte
	err := db.sql.QueryRowContext(ctx,
		`SELECT rendition FROM performances WHERE id = ?`, id).Scan(&rendition)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	} else if err != nil {
		return nil, err
	}
	return rendition, nil
}
