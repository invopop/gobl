package kr_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/kr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	regime := kr.New()
	require.NotNil(t, regime)
	require.NoError(t, rules.Validate(regime))
	assert.Equal(t, l10n.KR, regime.Country.Code())
	assert.Equal(t, "South Korea", regime.Name.String())
	assert.Equal(t, "KRW", regime.Currency.String())
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
	assert.NotEmpty(t, regime.Description)
	assert.Len(t, regime.Corrections, 1)
}

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("KR"),
		Series: "TEST",
		Code:   "0001",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "KR",
				Code:    "1208147521",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "KR",
				Code:    "2208162517",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Consulting services",
					Price: num.NewAmount(10000, 0),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice calculates 10% VAT", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, "1000", inv.Totals.Tax.String())
	})

	t.Run("invalid supplier BRN is rejected", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = "1208147520" // bad check digit
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "invalid Korean business registration number")
	})
}
