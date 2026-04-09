package favat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteHasKeyOrCodeGuard(t *testing.T) {
	assert.True(t, noteHasKeyOrCode("not a note"))
}
