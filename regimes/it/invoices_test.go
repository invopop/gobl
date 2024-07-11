package it_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
			PricesInclude: tax.CategoryVAT,
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
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestCustomerValidation(t *testing.T) {
	id := &org.Identity{
		Key:  it.IdentityKeyFiscalCode,
		Code: "RSSGNN60R30H501U",
	}
	t.Run("valid", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: l10n.IT,
			Code:    "",
		}
		inv.Customer.Identities = append(inv.Customer.Identities, id)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("missing tax_id", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = append(inv.Customer.Identities, id)
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "customer: (tax_id: cannot be blank.)")
	})

	t.Run("missing tax id code and identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: l10n.IT,
			Code:    "",
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		// ensure contains bother errors
		assert.ErrorContains(t, err, "identities: missing key it-fiscal-code")
		assert.ErrorContains(t, err, "tax_id: (code: cannot be blank.")
	})
}

func TestSupplierValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Supplier.TaxID = &tax.Identity{
		Country: l10n.IT,
		Code:    "RSSGNN60R30H501U",
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code: contains invalid characters")
}

func TestSupplierAddressesValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Supplier.Addresses = nil
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "addresses: cannot be blank.")

	inv = testInvoiceStandard(t)
	inv.Supplier.Addresses[0].Code = "123456"
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "supplier: (addresses: (0: (code: must be in a valid format.).).)")

	inv = testInvoiceStandard(t)
	inv.Customer.Addresses[0].Code = "123456"
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "customer: (addresses: (0: (code: must be in a valid format.).).)")
}

func TestRetainedTaxesValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
		Category: "IRPEF",
		Percent:  num.NewPercentage(20, 2),
	})
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "lines: (0: (taxes: (1: (ext: (it-sdi-retained-tax: required.).).).).).")
	}

	inv = testInvoiceStandard(t)
	inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
		Category: "IRPEF",
		Ext: tax.Extensions{
			it.ExtKeySDIRetainedTax: "A",
		},
		Percent: num.NewPercentage(20, 2),
	})
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestInvoiceLineTaxes(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Lines = append(inv.Lines, &bill.Line{
		Quantity: num.MakeAmount(10, 0),
		Item: &org.Item{
			Name:  "Test Item",
			Price: num.MakeAmount(10000, 2),
		},
		// No taxes!
	})
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.EqualError(t, err, "lines: (1: (taxes: missing category VAT.).).")
}
