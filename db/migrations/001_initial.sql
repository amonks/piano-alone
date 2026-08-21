-- The performances table as GORM's AutoMigrate left it in production,
-- transcribed so a fresh install gets the same shape the existing
-- database already has. The deployed database predates this runner and
-- adopts it as a baseline: the file is recorded, not executed, so the
-- columns are never rebuilt under the two recorded performances.
--
-- Booleans are integers and dates are text in "2006-01-02 15:04:05-07:00",
-- because that is what the rows already hold.
CREATE TABLE IF NOT EXISTS performances (
    id           text PRIMARY KEY,
    title        text,
    composer     text,
    score        blob,
    date         datetime,
    is_featured  numeric,
    is_complete  numeric,
    rendition    blob,
    player_count integer
);
