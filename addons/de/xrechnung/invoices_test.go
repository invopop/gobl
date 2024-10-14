package xrechnung_test

import (
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Type:   cbc.Key("standard"),
		Addons: tax.WithAddons(xrechnung.V3),
		Series: "A",
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "123456789",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "Cursor Subscription",
					Price: num.MakeAmount(1000, 3),
				},
			},
		},
	}
	return inv
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (tax_id: cannot be blank.)")
	})

	t.Run("with exemption reason", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext = nil
		assertValidationError(t, inv, "de-xrechnung-exemption: required")
	})

	t.Run("without series", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Series = ""
		assertValidationError(t, inv, "series: cannot be blank")
	})
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expectedError string) {
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	assert.ErrorContains(t, err, expectedError)
}
