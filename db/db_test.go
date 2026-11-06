package db_test

import (
	"bytes"
	"database/sql"
	"errors"
	"io/fs"
	"path/filepath"
	"testing"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	_ "modernc.org/sqlite"

	"monks.co/piano-alone/db"
	"monks.co/piano-alone/game"
)

// legacySchema is what GORM's AutoMigrate left in the deployed
// database, transcribed from it. The tests below open a database of
// this shape to prove the migration adopts it instead of rebuilding
// it, and that the rows it already holds still read.
const legacySchema = "CREATE TABLE `performances` (" +
	"`id` text,`title` text,`composer` text,`score` blob,`date` datetime," +
	"`is_featured` numeric,`is_complete` numeric,`rendition` blob," +
	"`player_count` integer,PRIMARY KEY (`id`))"

// writeLegacy writes the three performances the deployed database
// holds, with their dates in the two spellings that are actually in
// there: one at -05:00 and the rest at +00:00.
func writeLegacy(t *testing.T, path string) {
	t.Helper()
	handle, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	if _, err := handle.Exec(legacySchema); err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct {
		id, title, composer, date string
		featured, complete        int
		players                   int
		rendition                 []byte
	}{
		{"b2cc", "Prelude (from Suite Bergamasque)", "Claude Debussy", "2024-03-23 01:00:00+00:00", 1, 1, 12, []byte("debussy-midi")},
		{"4ba9", "Prelude in C♯ Minor", "Sergei Rachmaninoff", "2024-06-11 21:30:00-05:00", 1, 1, 6, []byte("rachmaninoff-midi")},
		{"2ad2", "Prelude in C♯ Minor", "Sergei Rachmaninoff", "2024-03-23 04:00:00+00:00", 0, 0, 0, nil},
	} {
		if _, err := handle.Exec(
			"INSERT INTO performances (id, title, composer, score, date, is_featured, is_complete, rendition, player_count) VALUES (?,?,?,?,?,?,?,?,?)",
			row.id, row.title, row.composer, []byte("score"), row.date,
			row.featured, row.complete, row.rendition, row.players); err != nil {
			t.Fatal(err)
		}
	}
}

// TestMigrationAdoptsTheDeployedSchema is the one that matters for the
// move into the monorepo: the deployed database holds two performances
// that exist nowhere else, and the first migration must be recorded
// against it rather than run.
func TestMigrationAdoptsTheDeployedSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "performances.db")
	writeLegacy(t, path)

	store, err := db.Open(t.Context(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	ps, err := store.GetFeaturedPerformances(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 2 {
		t.Fatalf("featured performances = %d, want 2", len(ps))
	}
	if got := ps[0].Configuration.Composer; got != "Claude Debussy" {
		t.Errorf("first featured composer = %q, want Claude Debussy", got)
	}
	if got := ps[0].PlayerCount; got != 12 {
		t.Errorf("first featured player count = %d, want 12", got)
	}

	// The two date spellings both read, and neither is silently
	// shifted: the second row is 21:30 at -05:00, which is 02:30 UTC.
	if got := ps[1].Date.UTC(); !got.Equal(time.Date(2024, 6, 12, 2, 30, 0, 0, time.UTC)) {
		t.Errorf("second featured date = %s, want 2024-06-12 02:30 UTC", got)
	}

	midi, err := store.GetMIDIFile(t.Context(), "b2cc")
	if err != nil {
		t.Fatal(err)
	}
	if string(midi) != "debussy-midi" {
		t.Errorf("rendition = %q, want debussy-midi", midi)
	}
}

// TestSeedLeavesAnExistingDatabaseAlone: seeding is for a fresh
// install, and the deployed database is not one. A seed that fired
// there would put a fourth performance in a table of three.
func TestSeedLeavesAnExistingDatabaseAlone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "performances.db")
	writeLegacy(t, path)

	store, err := db.Open(t.Context(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.Seed(t.Context()); err != nil {
		t.Fatal(err)
	}
	ps, err := store.GetScheduledPerformances(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 {
		t.Fatalf("scheduled performances after seeding = %d, want the 1 that was already there", len(ps))
	}
}

func TestSeedGivesAFreshInstallSomethingToPlay(t *testing.T) {
	store := openFresh(t)
	if err := store.Seed(t.Context()); err != nil {
		t.Fatal(err)
	}
	ps, err := store.GetScheduledPerformances(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 {
		t.Fatalf("scheduled performances = %d, want 1", len(ps))
	}
	// Seeding twice is not two performances.
	if err := store.Seed(t.Context()); err != nil {
		t.Fatal(err)
	}
	if ps, err = store.GetScheduledPerformances(t.Context()); err != nil {
		t.Fatal(err)
	} else if len(ps) != 1 {
		t.Fatalf("scheduled performances after seeding twice = %d, want 1", len(ps))
	}
}

// renditionCounted is a one-track rendition that strikes a key half a
// second in and releases it a second in, at 120 bpm, its deltas
// counted at 480 ticks a quarter note under a header declaring header.
// The renditions of the piece's two performances are shaped as
// renditionCounted(t, 960); every one written since, as
// renditionCounted(t, 480).
func renditionCounted(t *testing.T, header smf.MetricTicks) []byte {
	t.Helper()
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(120))
	tr.Add(480, midi.NoteOn(1, 60, 100))
	tr.Add(480, midi.NoteOff(1, 60))
	tr.Close(0)
	f := smf.New()
	f.TimeFormat = header
	if err := f.Add(tr); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := f.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// The renditions of the piece's two performances were written with
// their deltas counted at 480 ticks a quarter note under a header that
// said 960, so anything playing them — the disklavier included —
// played them at twice the speed their audiences played. Migration
// 002 rewrites a rendition's header division from 960 to 480: the
// header's two bytes alone, every event left as it was. Only a
// rendition written before the fix says 960; one written since says
// 480 and is left alone, as is anything in the column that is not a
// MIDI file. The correction run again changes nothing.
func TestRenditionsPlayedAtTwiceTheirSpeedAreCorrected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "performances.db")
	handle, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	if _, err := handle.Exec(legacySchema); err != nil {
		t.Fatal(err)
	}
	stored := map[string][]byte{
		"before": renditionCounted(t, 960),
		"since":  renditionCounted(t, 480),
		"other":  []byte("not a MIDI file"),
		"none":   nil,
	}
	for id, rendition := range stored {
		if _, err := handle.Exec(
			"INSERT INTO performances (id, title, composer, score, date, is_featured, is_complete, rendition, player_count) VALUES (?,?,?,?,?,?,?,?,?)",
			id, "Prelude", "Somebody", []byte("score"), "2024-06-11 21:30:00-05:00", 1, rendition != nil, rendition, 6); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Migrate(t.Context(), handle); err != nil {
		t.Fatal(err)
	}
	store := db.New(handle)

	want := map[string][]byte{
		"before": renditionCounted(t, 480),
		"since":  stored["since"],
		"other":  stored["other"],
		"none":   nil,
	}
	check := func(when string) {
		t.Helper()
		for id, w := range want {
			got, err := store.GetMIDIFile(t.Context(), id)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, w) {
				t.Errorf("%s, the rendition stored %s is %x, want %x", when, id, got, w)
			}
		}
	}
	check("after migrating")

	before, err := store.GetMIDIFile(t.Context(), "before")
	if err != nil {
		t.Fatal(err)
	}
	var at []time.Duration
	rd := smf.ReadTracksFrom(bytes.NewReader(before)).Do(func(ev smf.TrackEvent) {
		var ch, key, vel uint8
		if ev.Message.GetNoteStart(&ch, &key, &vel) || ev.Message.GetNoteEnd(&ch, &key) {
			at = append(at, time.Duration(ev.AbsMicroSeconds)*time.Microsecond)
		}
	})
	if err := rd.Error(); err != nil {
		t.Fatal(err)
	}
	if len(at) != 2 || at[0] != 500*time.Millisecond || at[1] != time.Second {
		t.Errorf("the corrected rendition sounds its key at %v, want at 500ms and 1s", at)
	}

	correction, err := fs.ReadFile(db.Migrations, db.MigrationsDir+"/002_rendition_division.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := handle.Exec(string(correction)); err != nil {
		t.Fatal(err)
	}
	check("after correcting again")
}

// TestTheDebussyIsTitledByItsExcerpt: the Debussy performance was
// titled "(first 4 seconds)", but its score is the first 5.25 seconds of
// the Prelude. The title is corrected; any other performance, and the
// Debussy if it has already been retitled, is left alone.
func TestTheDebussyIsTitledByItsExcerpt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "performances.db")
	handle, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	if _, err := handle.Exec(legacySchema); err != nil {
		t.Fatal(err)
	}
	rows := []struct{ id, title, composer string }{
		{"debussy", "Prelude (from Suite Bergamasque) (first 4 seconds)", "Claude Debussy"},
		{"rachmaninoff", "Prelude in C♯ Minor", "Sergei Rachmaninoff"},
		{"other", "Prelude (from Suite Bergamasque) (first 4 seconds)", "Somebody Else"},
	}
	for _, r := range rows {
		if _, err := handle.Exec(
			"INSERT INTO performances (id, title, composer, score, date, is_featured, is_complete, rendition, player_count) VALUES (?,?,?,?,?,?,?,?,?)",
			r.id, r.title, r.composer, []byte("score"), "2024-06-11 21:30:00-05:00", 1, 0, nil, 6); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Migrate(t.Context(), handle); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"debussy":      "Prelude (from Suite Bergamasque) (first 5 seconds)",
		"rachmaninoff": "Prelude in C♯ Minor",
		"other":        "Prelude (from Suite Bergamasque) (first 4 seconds)",
	}
	for id, w := range want {
		var got string
		if err := handle.QueryRow("SELECT title FROM performances WHERE id = ?", id).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != w {
			t.Errorf("%s is titled %q, want %q", id, got, w)
		}
	}
}

func TestRenditionRoundTrip(t *testing.T) {
	store := openFresh(t)
	perf := &game.Performance{
		Configuration: &game.Configuration{PerformanceID: "abc", Title: "Test", Composer: "Nobody"},
		Date:          time.Date(2026, 8, 21, 17, 0, 0, 0, time.UTC),
	}
	if err := store.SchedulePerformance(t.Context(), perf); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveRendition(t.Context(), "abc", 7, []byte("rendition")); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetPerformance(t.Context(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if !got.IsComplete || got.PlayerCount != 7 || string(got.Rendition) != "rendition" {
		t.Errorf("after SaveRendition: complete=%v players=%d rendition=%q",
			got.IsComplete, got.PlayerCount, got.Rendition)
	}
	if !got.Date.Equal(perf.Date) {
		t.Errorf("date round trip = %s, want %s", got.Date, perf.Date)
	}
}

// A finished performance's rendition is what its audience played and
// exists nowhere else, so it is written once: a second SaveRendition
// against the same id is refused and leaves the first standing. Playing
// the piece again is a new performance.
func TestACompletedRenditionIsNeverOverwritten(t *testing.T) {
	store := openFresh(t)
	perf := &game.Performance{
		Configuration: &game.Configuration{PerformanceID: "abc", Title: "Test", Composer: "Nobody"},
		Date:          time.Date(2026, 8, 21, 17, 0, 0, 0, time.UTC),
	}
	if err := store.SchedulePerformance(t.Context(), perf); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveRendition(t.Context(), "abc", 7, []byte("what they played")); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveRendition(t.Context(), "abc", 2, []byte("something else")); !errors.Is(err, db.ErrComplete) {
		t.Errorf("a second SaveRendition returned %v, want ErrComplete", err)
	}
	got, err := store.GetPerformance(t.Context(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if string(got.Rendition) != "what they played" || got.PlayerCount != 7 {
		t.Errorf("after a refused overwrite: rendition %q, players %d", got.Rendition, got.PlayerCount)
	}
}

// TestMissingPerformanceIsNotAnInternalError: the id comes from the
// URL, so a bad one is the caller's mistake. The handlers tell the two
// apart by this error.
func TestMissingPerformanceIsNotAnInternalError(t *testing.T) {
	store := openFresh(t)
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"GetPerformance", errOf(store.GetPerformance(t.Context(), "nope"))},
		{"GetMIDIFile", errOf(store.GetMIDIFile(t.Context(), "nope"))},
		{"SaveRendition", store.SaveRendition(t.Context(), "nope", 1, []byte("x"))},
	} {
		if tc.err == nil {
			t.Errorf("%s on a missing id returned no error", tc.name)
		} else if !errors.Is(tc.err, db.ErrNotFound) {
			t.Errorf("%s on a missing id returned %v, want ErrNotFound", tc.name, tc.err)
		}
	}
}

func openFresh(t *testing.T) *db.DB {
	t.Helper()
	store, err := db.Open(t.Context(), filepath.Join(t.TempDir(), "performances.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func errOf[T any](_ T, err error) error { return err }
