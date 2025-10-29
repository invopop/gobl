package l10n_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
)

func TestCountries(t *testing.T) {
	t.Parallel()
	c := l10n.Countries().Code(l10n.US)
	assert.NotNil(t, c)
}
