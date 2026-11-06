-- A rendition written before this migration counted its deltas at 480
-- ticks a quarter note under a header that declared 960, so anything
-- playing it, the disklavier included, played it at twice the speed
-- its audience played. This rewrites the header's division (the MThd
-- chunk's last two bytes, the file's thirteenth and fourteenth) from
-- 960 to 480 and leaves every other byte as it was.
--
-- A rendition is written at 480 ticks a quarter note, its header and
-- its deltas both, from the release this migration ships in, so a
-- header that declares 960 is one written before it. A rendition
-- written since, and anything in the column that is not a standard
-- MIDI file, is left alone, and running this again changes nothing.
-- The blob is spliced as hex because SQLite has no byte-level splice,
-- and unhex returns a blob.
UPDATE performances
SET rendition = unhex(substr(hex(rendition), 1, 24) || '01E0' || substr(hex(rendition), 29))
WHERE substr(rendition, 1, 8) = X'4D54686400000006'
  AND substr(rendition, 13, 2) = X'03C0';
