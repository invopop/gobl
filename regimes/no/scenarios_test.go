package no_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/no"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarios(t *testing.T) {
	t.Parallel()

	t.Run("reverse charge tag", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Tags = tax.Tags{List: []cbc.Key{tax.TagReverseCharge}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		found := false
		for _, n := range inv.Notes {
			if n.Src == tax.TagReverseCharge {
				assert.Contains(t, n.Text, "Omvendt avgiftsplikt")
				found = true
			}
		}
		assert.True(t, found, "expected reverse charge note")
	})

	t.Run("foretaksregisteret tag", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Tags = tax.Tags{List: []cbc.Key{no.TagForetaksregisteret}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		found := false
		for _, n := range inv.Notes {
			if n.Src == no.TagForetaksregisteret {
				assert.Equal(t, "Foretaksregisteret", n.Text)
				found = true
			}
		}
		assert.True(t, found, "expected foretaksregisteret note")
	})

	t.Run("both tags simultaneously", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Tags = tax.Tags{List: []cbc.Key{tax.TagReverseCharge, no.TagForetaksregisteret}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		foundRC := false
		foundFR := false
		for _, n := range inv.Notes {
			if n.Src == tax.TagReverseCharge {
				foundRC = true
			}
			if n.Src == no.TagForetaksregisteret {
				foundFR = true
			}
		}
		assert.True(t, foundRC, "expected reverse charge note")
		assert.True(t, foundFR, "expected foretaksregisteret note")
	})
}
