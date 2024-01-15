package pianists_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"monks.co/piano-alone/pianists"
)

func TestPianistHash(t *testing.T) {
	assert.Equal(t, "Sylviane Deferne", pianists.Hash("andrew"))
	assert.Equal(t, "Arabella Goddard", pianists.Hash("bob"))
}
