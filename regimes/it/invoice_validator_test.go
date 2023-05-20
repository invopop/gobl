package it_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Code:     "123TEST",
		Currency: "EUR",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.IT,
				Code:    "12345678903",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: l10n.IT,
				Type:    it.TaxIdentityTypeBusiness,
				Code:    "13029381004",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}
	return i
}

func TestInvoiceValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	inv = testInvoiceStandard(t)
	inv.Customer.TaxID = &tax.Identity{
		Country: l10n.IT,
		Type:    it.TaxIdentityTypeIndividual,
		Code:    "RSSGNN60R30H501U",
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	inv = testInvoiceStandard(t)
	inv.Supplier.TaxID = &tax.Identity{
		Country: l10n.IT,
		Type:    it.TaxIdentityTypeIndividual,
		Code:    "RSSGNN60R30H501U",
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type: must be a valid value")

}
