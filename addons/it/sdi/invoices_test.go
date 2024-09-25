package sdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Regime:   tax.WithRegime("IT"),
		Addons:   tax.WithAddons(sdi.V1),
		Code:     "123TEST",
		Currency: "EUR",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Type: bill.InvoiceTypeStandard,
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "IT",
				Code:    "12345678903",
			},
			Addresses: []*org.Address{
				{
					Street:   "Via di Test",
					Code:     "12345",
					Locality: "Rome",
					Country:  "IT",
					Number:   "3",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "IT",
				Code:    "13029381004",
			},
			Addresses: []*org.Address{
				{
					Street:   "Piazza di Test",
					Code:     "38342",
					Locality: "Venezia",
					Country:  "IT",
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

func TestInvoiceNormalization(t *testing.T) {
	ad := tax.AddonForKey(sdi.V1)

	t.Run("supplier fiscal regime", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		ad.Normalizer(inv)
		assert.Equal(t, "RF01", inv.Supplier.Ext[sdi.ExtKeyFiscalRegime].String())
	})
}

func TestCustomerValidation(t *testing.T) {
	id := &org.Identity{
		Key:  it.IdentityKeyFiscalCode,
		Code: "RSSGNN60R30H501U",
	}
	t.Run("valid", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "IT",
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
			Country: "IT",
			Code:    "",
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		// ensure contains bother errors
		assert.ErrorContains(t, err, "identities: missing key it-fiscal-code")
		assert.ErrorContains(t, err, "tax_id: (code: cannot be blank.")
	})

	t.Run("payment advances", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.Payment{
			Advances: []*pay.Advance{
				{
					Description: "Paid up front",
					Percent:     num.NewPercentage(100, 3),
					Key:         "card",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("payment terms missing instructions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.Payment{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date:    cal.NewDate(2022, 6, 13),
						Percent: num.NewPercentage(100, 3),
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "payment: (instructions: cannot be blank when terms are present.)")
	})

	t.Run("payment terms with instructions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.Payment{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date:    cal.NewDate(2022, 6, 13),
						Percent: num.NewPercentage(100, 3),
					},
				},
			},
			Instructions: &pay.Instructions{
				Key: "card",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
		assert.Equal(t, "MP08", inv.Payment.Instructions.Ext[sdi.ExtKeyPaymentMeans].String())
	})
}

func TestSupplierValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Supplier.TaxID = &tax.Identity{
		Country: "IT",
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
		assert.Contains(t, err.Error(), "lines: (0: (taxes: (1: (ext: (it-sdi-retained: required.).).).).).")
	}

	inv = testInvoiceStandard(t)
	inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
		Category: "IRPEF",
		Ext: tax.Extensions{
			sdi.ExtKeyRetained: "A",
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
