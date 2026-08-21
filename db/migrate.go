package db

import (
	"context"
	"database/sql"
	"time"

	"monks.co/pkg/migrate"

	"monks.co/piano-alone/game"
	"monks.co/piano-alone/id"
	"monks.co/piano-alone/songs"
)

// Migrate brings a handle's schema up to date. A host that has its own
// migration runner should use Migrations and Baseline instead of
// calling this.
func Migrate(ctx context.Context, handle *sql.DB) error {
	return migrate.Run(ctx, migrate.Config{
		DB:       handle,
		FS:       Migrations,
		Dir:      MigrationsDir,
		Baseline: Baseline,
	})
}

// Seed puts one unplayed performance in an empty table, so a fresh
// install has something for a conductor to begin. It is the standalone
// command's, not a host's: opening a database used to insert this row
// as a side effect, and a store that writes when you read it is a
// surprise a host does not need. An existing database — which is any
// database that has ever held a performance — is left alone.
func (db *DB) Seed(ctx context.Context) error {
	var n int
	if err := db.sql.QueryRowContext(ctx, `SELECT count(*) FROM performances`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return db.SchedulePerformance(ctx, &game.Performance{
		Configuration: &game.Configuration{
			PerformanceID: id.Random128(),
			Score:         songs.PreludeOpus3No2Bytes,
			Title:         "Prelude in C♯ Minor",
			Composer:      "Sergei Rachmaninoff",
		},
		Date: time.Now(),
	})
}
