package it_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
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
		Type: bill.InvoiceTypeStandard,
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.IT,
				Code:    "12345678903",
			},
			Addresses: []*org.Address{
				{
					Street:   "Via di Test",
					Code:     "12345",
					Locality: "Rome",
					Country:  l10n.IT,
					Number:   "3",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: l10n.IT,
				Type:    it.TaxIdentityTypeBusiness,
				Code:    "13029381004",
			},
			Addresses: []*org.Address{
				{
					Street:   "Piazza di Test",
					Code:     "38342",
					Locality: "Venezia",
					Country:  l10n.IT,
					Number:   "1",
				},
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
	ctx := context.Background()
	require.NoError(t, inv.Calculate(ctx))
	require.NoError(t, inv.Validate())
}

func TestCustomerValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Customer.TaxID = &tax.Identity{
		Country: l10n.IT,
		Type:    it.TaxIdentityTypeIndividual,
		Code:    "RSSGNN60R30H501U",
	}
	ctx := context.Background()
	require.NoError(t, inv.Calculate(ctx))
	require.NoError(t, inv.Validate())

}

func TestSupplierValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Supplier.TaxID = &tax.Identity{
		Country: l10n.IT,
		Type:    it.TaxIdentityTypeIndividual,
		Code:    "RSSGNN60R30H501U",
	}
	ctx := context.Background()
	require.NoError(t, inv.Calculate(ctx))
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type: must be a valid value")
}

func TestRetainedTaxesValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
		Category: "IRPEF",
		Percent:  num.NewPercentage(20, 2),
	})
	ctx := context.Background()
	require.NoError(t, inv.Calculate(ctx))
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lines: (0: (taxes: 1: rate: cannot be blank..).).")

	inv = testInvoiceStandard(t)
	inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
		Category: "IRPEF",
		Rate:     cbc.Key("self-employed-habitual"),
		Percent:  num.NewPercentage(20, 2),
	})
	require.NoError(t, inv.Calculate(ctx))
	require.NoError(t, inv.Validate())
}
