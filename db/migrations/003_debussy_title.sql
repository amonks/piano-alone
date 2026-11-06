-- The Debussy performance was titled "(first 4 seconds)", but its score
-- is the Prelude's first 5.25 seconds (the last note sounds at 5.12 s),
-- which is what an audience played and what the title names. Played at
-- its corrected speed (002), the rendition runs longer still only
-- because one player struck a key once more after the score ended.
-- Running this again changes nothing.
UPDATE performances
SET title = 'Prelude (from Suite Bergamasque) (first 5 seconds)'
WHERE title = 'Prelude (from Suite Bergamasque) (first 4 seconds)'
  AND composer = 'Claude Debussy';
