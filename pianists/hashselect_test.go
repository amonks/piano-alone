package pianists

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSelectIsAlwaysBelowCap(t *testing.T) {
	assert.Equal(t, hashSelect([]byte{}, 1), 0)
	assert.Equal(t, hashSelect([]byte{0}, 1), 0)
	assert.Equal(t, hashSelect([]byte{1}, 1), 0)
}

func TestHashSelectVaries(t *testing.T) {
	got := map[int]struct{}{}
	for i := 0; i < 256; i++ {
		got[hashSelect([]byte{byte(i)}, 1e3)] = struct{}{}
	}
	assert.Len(t, got, 256)
}
