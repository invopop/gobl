package verifactu_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
		assert.Empty(t, inv.Customer.Identities[0].Ext)
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
		assert.Equal(t, verifactu.ExtCodeIdentityTypePassport, inv.Customer.Identities[0].Ext[verifactu.ExtKeyIdentityType])
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
		assert.Equal(t, verifactu.ExtCodeIdentityTypeForeign, inv.Customer.Identities[0].Ext[verifactu.ExtKeyIdentityType])
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
		assert.Equal(t, verifactu.ExtCodeIdentityTypeResident, inv.Customer.Identities[0].Ext[verifactu.ExtKeyIdentityType])
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
		assert.Equal(t, verifactu.ExtCodeIdentityTypeOther, inv.Customer.Identities[0].Ext[verifactu.ExtKeyIdentityType])
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
		assert.Empty(t, inv.Customer.Identities[0].Ext)
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
		assert.Equal(t, verifactu.ExtCodeIdentityTypePassport, inv.Customer.Identities[0].Ext[verifactu.ExtKeyIdentityType])
		assert.Empty(t, inv.Customer.Identities[1].Ext)
	})

	t.Run("self-billed", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, verifactu.ExtCodeIssuerTypeCustomer, inv.Tax.Ext[verifactu.ExtKeyIssuerType])
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
		assert.Equal(t, verifactu.ExtCodeIssuerTypeThirdParty, inv.Tax.Ext[verifactu.ExtKeyIssuerType])
	})
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Equal(t, inv.Tax.Ext[verifactu.ExtKeyDocType].String(), "F1")
	})
	t.Run("standard invoice without customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.ErrorContains(t, inv.Validate(), "customer: cannot be blank.")
	})
	t.Run("missing doc type", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = nil
		err := inv.Validate()
		require.ErrorContains(t, err, "es-verifactu-doc-type: required")
	})

	t.Run("note too long", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Notes = []*org.Note{
			{
				Key:  org.NoteKeyGeneral,
				Text: strings.Repeat("a", 501),
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.ErrorContains(t, err, "text: the length must be no more than 500")
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
		require.NoError(t, inv.Validate())
	})

	t.Run("simplified invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Equal(t, inv.Tax.Ext[verifactu.ExtKeyDocType].String(), "F2")
	})

	t.Run("simplified substitution", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		require.NoError(t, inv.Calculate())

		require.NoError(t, inv.Correct(bill.Corrective, bill.WithCopyTax(), bill.WithExtension(verifactu.ExtKeyDocType, "F3")))
		require.NoError(t, inv.Validate())
		// Should always set the doc type to R5, even if trying to override as the simplified
		// tag has priority.
		assert.Equal(t, "R5", inv.Tax.Ext[verifactu.ExtKeyDocType].String())
		assert.Equal(t, "S", inv.Tax.Ext[verifactu.ExtKeyCorrectionType].String())
	})

	t.Run("correction invoice requires preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		assertValidationError(t, inv, "preceding: cannot be blank")
	})
	t.Run("correction invoice nil preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{nil}
		require.NoError(t, inv.Calculate())
		ad := tax.AddonForKey(verifactu.V1)
		assert.NoError(t, ad.Validator(inv))
	})

	t.Run("credit-note invoice preceding requires issue date", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "123",
			},
		}
		assertValidationError(t, inv, "preceding: (0: (issue_date: cannot be blank.).")
	})

	t.Run("correction invoice preceding requires issue date and tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "123",
			},
		}
		assertValidationError(t, inv, "preceding: (0: (issue_date: cannot be blank; tax: cannot be blank.).")
	})

	t.Run("correction invoice with preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		d := cal.MakeDate(2024, 1, 1)
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "ABC",
				Code:      "122",
				IssueDate: &d,
				Ext: tax.Extensions{
					verifactu.ExtKeyDocType: "R1",
				},
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
		require.NoError(t, inv.Validate())
		assert.Equal(t, inv.Tax.Ext[verifactu.ExtKeyDocType].String(), "R1")
		assert.Empty(t, inv.Preceding[0].Ext)
		assert.Equal(t, "21.00", inv.Preceding[0].Tax.Sum.String())
	})

	t.Run("replacement without preceding", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags("replacement")
		require.NoError(t, inv.Calculate())
		require.ErrorContains(t, inv.Validate(), "preceding: details of invoice being replaced must be included")
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
		require.NoError(t, inv.Validate())
	})

	t.Run("correction invoice preceding requires issue date and tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "123",
			},
		}
		assertValidationError(t, inv, "preceding: (0: (issue_date: cannot be blank; tax: cannot be blank.).")
	})

	t.Run("customer nil", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		ad := tax.AddonForKey(verifactu.V1)
		assert.NoError(t, ad.Validator(inv))
	})
	t.Run("customer with missing ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "customer: must have a tax_id, or an identity with ext 'es-verifactu-identity-type'")
	})
	t.Run("customer with missing Tax ID code", func(t *testing.T) {
		// VERI*FACTU has no way to handle just a country without an actual code.
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "customer: (tax_id: (code: cannot be blank.).)")
	})
	t.Run("customer with identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.ErrorContains(t, err, expected)
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
