package pianists

import (
	_ "embed"
	"strings"
)

//go:embed pianists.txt
var pianistsTxt string

// Pianists is a list of 1176 pianists, taken from Wikipedia.
var Pianists = strings.Split(pianistsTxt, "\n")

// Hash, given a string, produces a pianist. It will always
// return the same pianist when given the same string.
func Hash(input string) string {
	idx := hashSelect([]byte(input), len(Pianists))
	return Pianists[idx]
}
