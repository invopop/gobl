package verifactu_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoicePartyNormalization(t *testing.T) {
	t.Run("regular Spanish customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "ES",
			Code:    "B12345678",
		}
		require.NoError(t, inv.Calculate())
	})

	t.Run("Spanish customer with identities should not be normalized", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "ES",
			Code:    "B12345678",
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
		}
		require.NoError(t, inv.Calculate())
		// Should not have extension as Spanish NIFs are already handled
		assert.True(t, inv.Customer.Identities[0].Ext.IsZero())
	})

	t.Run("customer without identities", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		// Should not cause any issues
	})

	t.Run("passport identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIdentityTypePassport, inv.Customer.Identities[0].Ext.Get(verifactu.ExtKeyIdentityType))
	})

	t.Run("foreign identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyForeign,
				Code: "FOR123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIdentityTypeForeign, inv.Customer.Identities[0].Ext.Get(verifactu.ExtKeyIdentityType))
	})

	t.Run("resident identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyResident,
				Code: "RES123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIdentityTypeResident, inv.Customer.Identities[0].Ext.Get(verifactu.ExtKeyIdentityType))
	})

	t.Run("other identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyOther,
				Code: "OTH123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIdentityTypeOther, inv.Customer.Identities[0].Ext.Get(verifactu.ExtKeyIdentityType))
	})

	t.Run("unknown identity key not normalized", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  "unknown",
				Code: "UNK123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.True(t, inv.Customer.Identities[0].Ext.IsZero())
	})

	t.Run("multiple identities only normalizes first", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
			{
				Key:  org.IdentityKeyForeign,
				Code: "FOR123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIdentityTypePassport, inv.Customer.Identities[0].Ext.Get(verifactu.ExtKeyIdentityType))
		assert.True(t, inv.Customer.Identities[1].Ext.IsZero())
	})

	t.Run("self-billed", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIssuerTypeCustomer, inv.Tax.Ext.Get(verifactu.ExtKeyIssuerType))
	})

	t.Run("with issuer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{
			Issuer: &org.Party{
				Name: "Test Issuer",
				TaxID: &tax.Identity{
					Country: "ES",
					Code:    "B12345678",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIssuerTypeThirdParty, inv.Tax.Ext.Get(verifactu.ExtKeyIssuerType))
	})
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, inv.Tax.Ext.Get(verifactu.ExtKeyDocType).String(), "F1")
	})
	t.Run("standard invoice without customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-06] ($.customer) customer is required")
	})
	t.Run("missing doc type", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = tax.Extensions{}
		err := rules.Validate(inv)
		require.ErrorContains(t, err, "doc type is required")
	})

	t.Run("note too long", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Notes = []*org.Note{
			{
				Key:  org.NoteKeyGeneral,
				Text: strings.Repeat("a", 501),
			},
		}
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-14] ($.notes[0].text) general note text must be 500 characters or less")
	})

	t.Run("note with wrong key", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLoading,
				Text: strings.Repeat("a", 501),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, inv.Tax.Ext.Get(verifactu.ExtKeyDocType).String(), "F2")
	})

	t.Run("simplified substitution without customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		// Simplified invoice without customer details stays F2
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "F2", inv.Tax.Ext.Get(verifactu.ExtKeyDocType).String())

		require.NoError(t, inv.Correct(bill.Corrective, bill.WithCopyTax()))
		require.NoError(t, rules.Validate(inv))
		// Should get R5 for simplified corrective
		assert.Equal(t, "R5", inv.Tax.Ext.Get(verifactu.ExtKeyDocType).String())
		assert.Equal(t, "S", inv.Tax.Ext.Get(verifactu.ExtKeyCorrectionType).String())
	})

	t.Run("corrective invoice requires preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-01] ($.preceding) preceding documents are required for corrective invoices")
	})
	t.Run("corrective invoice nil preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Preceding = []*org.DocumentRef{nil}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note needs no preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("corrective invoice preceding requires issue date and tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "123",
			},
		}
		assertValidationError(t, inv,
			"[GOBL-ES-VERIFACTU-BILL-INVOICE-02] ($.preceding[0].issue_date) issue date is required",
			"[GOBL-ES-VERIFACTU-BILL-INVOICE-03] ($.preceding[0].tax) preceding invoice tax data is required",
		)
	})

	t.Run("corrective invoice with preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		d := cal.MakeDate(2024, 1, 1)
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "ABC",
				Code:      "122",
				IssueDate: &d,
				Ext: tax.ExtensionsOf(tax.ExtMap{
					verifactu.ExtKeyDocType: "R1",
				}),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		data, _ := json.MarshalIndent(inv, "", "  ")
		t.Log(string(data))
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, inv.Tax.Ext.Get(verifactu.ExtKeyDocType).String(), "R1")
		assert.True(t, inv.Preceding[0].Ext.IsZero())
		assert.Equal(t, "21.00", inv.Preceding[0].Tax.Sum.String())
	})

	t.Run("replacement without preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags("replacement")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("replacement with preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags("replacement")
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "SAMPLE",
				Code:      "003",
				IssueDate: cal.NewDate(2025, 7, 1),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("correction invoice preceding requires issue date and tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "123",
			},
		}
		assertValidationError(t, inv,
			"[GOBL-ES-VERIFACTU-BILL-INVOICE-02] ($.preceding[0].issue_date) issue date is required",
			"[GOBL-ES-VERIFACTU-BILL-INVOICE-03] ($.preceding[0].tax) preceding invoice tax data is required",
		)
	})

	t.Run("customer nil", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
	t.Run("customer with missing ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-07] ($.customer) must have a tax_id or an identity with ext 'es-verifactu-identity-type'")
	})
	t.Run("customer with missing Tax ID code", func(t *testing.T) {
		// VERI*FACTU has no way to handle just a country without an actual code.
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID.Code = ""
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-08] ($.customer.tax_id.code) tax ID must have a code")
	})
	t.Run("customer with identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:     org.IdentityKeyPassport,
				Country: "GB",
				Code:    "AA123456",
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
	t.Run("customer with identity missing country", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
		}
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-17] ($.customer.identities[0].country) country is required when ext 'es-verifactu-identity-type' is not 02 (NIF-VAT)")
	})
	t.Run("customer with NIF-VAT identity without country is valid", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "B12345678",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					verifactu.ExtKeyIdentityType: "02",
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
	t.Run("simplified invoice with customer without tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, inv.Tax.Ext.Get(verifactu.ExtKeyDocType).String(), "F2")
	})
	t.Run("simplified substitution with customer without tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		d := cal.MakeDate(2024, 1, 1)
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "ABC",
				Code:      "122",
				IssueDate: &d,
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, inv.Tax.Ext.Get(verifactu.ExtKeyDocType).String(), "R5")
	})
	t.Run("simplified invoice F2 with customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		// Customer has tax ID - should be normalized to F1 with SimplifiedArt7273
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-04] ($.customer.tax_id) customer tax ID must not be set for simplified invoices")
	})
	t.Run("simplified substitution R5 with customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Type = bill.InvoiceTypeCorrective
		d := cal.MakeDate(2024, 1, 1)
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "ABC",
				Code:      "122",
				IssueDate: &d,
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
		}
		// Customer has tax ID - should be normalized to R1 with SimplifiedArt7273
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-04] ($.customer.tax_id) customer tax ID must not be set for simplified invoices")
	})
	t.Run("simplified invoice F2 with customer identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
		}
		// Customer has identity - should be normalized to F1 with SimplifiedArt7273
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-05] ($.customer) customer identity type extension not allowed for simplified invoices")
	})

	t.Run("invoice with only retained taxes fails", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Replace VAT with IRPF (retained tax)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "IRPF",
				Rate:     "pro",
			},
		}
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-15] ($.lines[0].taxes) must include at least one of VAT, IGIC, or IPSI")
	})

	t.Run("invoice with VAT and IRPF passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Rate:     "standard",
			},
			{
				Category: "IRPF",
				Rate:     "pro",
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier name over 120 chars", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Name = strings.Repeat("a", 121)
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-18] ($.supplier.name) supplier name must be 120 characters or less")
	})

	t.Run("customer name over 120 chars", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Name = strings.Repeat("a", 121)
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-19] ($.customer.name) customer name must be 120 characters or less")
	})

	t.Run("invoice series and code fit within 60 chars", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Series = cbc.Code(strings.Repeat("S", 30))
		inv.Code = cbc.Code(strings.Repeat("C", 30))
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-20] invoice series and code combined must be 60 characters or less")
	})

	t.Run("preceding series and code fit within 60 chars", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Preceding = []*org.DocumentRef{
			{
				Series: cbc.Code(strings.Repeat("S", 30)),
				Code:   cbc.Code(strings.Repeat("C", 30)),
			},
		}
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-21] ($.preceding[0]) preceding series and code combined must be 60 characters or less")
	})

	t.Run("non-ES customer tax ID code over 18 chars", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "NL",
			Code:    cbc.Code(strings.Repeat("X", 19)),
		}
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-22] ($.customer.tax_id.code) non-Spanish customer tax ID code must be 18 characters or less")
	})

	t.Run("non-retained tax rates cannot exceed 12", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = make([]*bill.Line, 0, 13)
		// 13 lines, each with a unique (percent, regime) combination so that
		// calculation groups each into its own rate entry.
		regimes := []cbc.Code{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "14", "15"}
		for i, r := range regimes {
			inv.Lines = append(inv.Lines, &bill.Line{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item",
					Price: num.NewAmount(int64(100+i), 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(int64(10+i), 2),
						Ext: tax.Extensions{
							verifactu.ExtKeyOpClass: "S1",
							verifactu.ExtKeyRegime:  r,
						},
					},
				},
			})
		}
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-23] ($.totals.taxes.categories) non-retained tax rates cannot exceed 12")
	})

	t.Run("non-EUR currency without exchange rates", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetRegime("ES")
		inv.Currency = "USD"
		assertValidationError(t, inv, "[GOBL-ES-VERIFACTU-BILL-INVOICE-16] invoice must be in EUR or provide exchange rate for conversion")
	})

	t.Run("non-EUR currency with exchange rates", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetRegime("ES")
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   "USD",
				To:     "EUR",
				Amount: num.MakeAmount(875967, 6), // 0.875967
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})

}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected ...string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	for _, e := range expected {
		require.ErrorContains(t, err, e)
	}
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Addons: tax.WithAddons(verifactu.V1),
		Code:   "123",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "NL",
				Code:    "000099995B57",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
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
