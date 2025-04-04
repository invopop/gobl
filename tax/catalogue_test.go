package tax_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestAllCatalogueDefs(t *testing.T) {
	cds := tax.AllCatalogueDefs()
	assert.GreaterOrEqual(t, len(cds), 1)
	match := true
	for _, cd := range cds {
		if cd.Key == "untdid" {
			match = true
			break
		}
	}
	assert.True(t, match)
}

func TestRegisterCatalogueDef(t *testing.T) {
	t.Run("missing source file", func(t *testing.T) {
		assert.Panics(t, func() {
			tax.RegisterCatalogueDef("missing.json")
		})
	})
}
