package saft_test

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoice(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FT", inv.Tax.Ext[saft.ExtKeyInvoiceType].String())
		assert.Equal(t, "NOR", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())

		assert.NoError(t, inv.Validate())
	})

	t.Run("prepaid", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Percent:     num.NewPercentage(1, 0),
					Description: "prepaid",
				},
			},
		}
		inv.Series = "FR SERIES-A"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FR", inv.Tax.Ext[saft.ExtKeyInvoiceType].String())
		assert.NoError(t, inv.Validate())
	})

	t.Run("reduced", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Taxes[0].Rate = tax.RateReduced
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "RED", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.NoError(t, inv.Validate())
	})

	t.Run("intermediate", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Taxes[0].Rate = tax.RateIntermediate
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "INT", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.NoError(t, inv.Validate())
	})

	t.Run("exempt", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Taxes[0].Rate = tax.RateExempt
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "ISE", inv.Lines[0].Taxes[0].Ext[saft.ExtKeyTaxRate].String())
		assert.ErrorContains(t, inv.Validate(), "lines: (0: (taxes: (0: (ext: (pt-saft-exemption: required.).).).).)")

		inv.Lines[0].Taxes[0].Ext[saft.ExtKeyExemption] = "M04"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}

func TestExemptions(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)
	require.NotNil(t, addon)

	exmpts := cbc.GetKeyDefinition(saft.ExtKeyExemption, addon.Extensions)
	require.NotNil(t, exmpts)

	for _, ex := range exmpts.Values {
		tn := fmt.Sprintf("Note for %s exemption code", ex.Code)
		t.Run(tn, func(t *testing.T) {
			inv := validInvoice()
			inv.Lines[0].Taxes[0].Rate = tax.RateExempt
			inv.Lines[0].Taxes[0].Ext = tax.Extensions{
				saft.ExtKeyExemption: ex.Code,
			}
			require.NoError(t, inv.Calculate())

			if assert.Len(t, inv.Notes, 1) {
				assert.Equal(t, org.NoteKeyLegal, inv.Notes[0].Key)
				assert.Equal(t, ex.Code, inv.Notes[0].Code)
				assert.Equal(t, saft.ExtKeyExemption, inv.Notes[0].Src)
				assert.LessOrEqual(t, len(inv.Notes[0].Text), 60, "for use in SAF-T, length must be 60 characters or less")
			}
		})
	}
}
