package mt_test

import (
	"testing"

	_ "github.com/invopop/gobl" // ensure all regimes loaded
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("MT"),
		Series:    "TEST",
		Code:      "0001",
		Currency:  currency.EUR,
		IssueDate: cal.MakeDate(2026, 1, 15),
		Supplier: &org.Party{
			Name: "Test Supplier Ltd",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.MT),
				Code:    "12701906",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer Ltd",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.MT),
				Code:    "12357210",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Consulting services",
					Price: num.NewAmount(100000, 2),
				},
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Rate: tax.KeyStandard},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Parallel()

	t.Run("standard invoice", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier without tax ID code", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "INVOICE-01")
		assert.ErrorContains(t, err, "must have a VAT tax ID code")
	})
}
