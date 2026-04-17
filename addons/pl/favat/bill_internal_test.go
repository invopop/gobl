package favat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteHasKeyOrCodeGuard(t *testing.T) {
	assert.False(t, noteHasKeyOrCode("not a note"))
}
