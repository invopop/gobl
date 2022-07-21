package l10n_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
)

func TestCodeIn(t *testing.T) {
	c := l10n.ES

	assert.True(t, c.In(l10n.PT, l10n.ES))
	assert.False(t, c.In(l10n.GB, l10n.PT))
}
