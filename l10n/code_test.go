package l10n_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
)

func TestCodeIn(t *testing.T) {
	c := l10n.Code("MAD")

	assert.True(t, c.In("A", "MAD"))
	assert.False(t, c.In("A", "V"))
}
