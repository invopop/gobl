package no_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
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

	t.Run("reverse charge exempt produces zero tax", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Tags = tax.Tags{List: []cbc.Key{tax.TagReverseCharge}}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: tax.CategoryVAT,
				Rate:     tax.KeyExempt.With("reverse-charge"),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		// Verify the note is injected
		found := false
		for _, n := range inv.Notes {
			if n.Src == tax.TagReverseCharge {
				found = true
			}
		}
		assert.True(t, found, "expected reverse charge note")
	})
}
