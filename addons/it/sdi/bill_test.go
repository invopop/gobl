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
			Ext: tax.Extensions{
				sdi.ExtKeyDocumentType: "TD01",
				sdi.ExtKeyFormat:       "FPA12",
			},
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
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
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
	t.Run("basic", func(t *testing.T) {

		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
	t.Run("missing tax extensions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = nil
		require.ErrorContains(t, inv.Validate(), "tax: (ext: (it-sdi-document-type: required; it-sdi-format: required.).)")
	})
}

func TestInvoiceNormalization(t *testing.T) {
	ad := tax.AddonForKey(sdi.V1)

	t.Run("supplier fiscal regime", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		ad.Normalizer(inv)
		assert.Equal(t, "RF01", inv.Supplier.Ext[sdi.ExtKeyFiscalRegime].String())
	})
}

func TestSupplierValidation(t *testing.T) {
	t.Run("with supplier registration details", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Registration = &org.Registration{
			Entry:  "123456",
			Office: "Rome",
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("with supplier missing registration details", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Registration = &org.Registration{
			Entry: "123456",
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (registration: (office: cannot be blank.).).")
	})

	t.Run("with invalid tax ID code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "RSSGNN60R30H501U",
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code: contains invalid characters")
	})

	t.Run("missing supplier", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		ad := tax.AddonForKey(sdi.V1)
		inv.Supplier = nil
		ad.Normalizer(inv)
		assert.NoError(t, ad.Validator(inv))
	})

	t.Run("valid Latin-1 supplier name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with valid Latin-1 characters including accented characters
		inv.Supplier.Name = "Società di Test SRL àáâãäåæçèéêë"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("invalid supplier name with non-Latin-1 characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with emoji (outside Latin-1 range)
		inv.Supplier.Name = "Test Supplier 😊"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: contains characters outside of Latin and Latin-1 range.).")
	})

	t.Run("invalid supplier name with Greek characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with Greek characters (outside Latin-1 range)
		inv.Supplier.Name = "Test Supplier αβγδε"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: contains characters outside of Latin and Latin-1 range.).")
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
		assert.ErrorContains(t, err, "identities: missing key 'it-fiscal-code'")
		assert.ErrorContains(t, err, "tax_id: (code: cannot be blank.")
	})

	t.Run("missing address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (addresses: cannot be blank.).")
	})

	t.Run("missing customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: cannot be blank.")
	})

	t.Run("valid Latin-1 customer name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with valid Latin-1 characters including special symbols
		inv.Customer.Name = "Cliente & Cia. S.p.A. ñöüß"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("invalid customer name with Chinese characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with Chinese characters (outside Latin-1 range)
		inv.Customer.Name = "测试客户"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (name: contains characters outside of Latin and Latin-1 range.).")
	})

	t.Run("missing customer name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (name: cannot be blank.).")
	})

	t.Run("missing customer people with identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID.Code = ""
		inv.Customer.Name = ""
		inv.Customer.Identities = append(inv.Customer.Identities, id)
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (name: cannot be blank; people: cannot be blank.).")
	})

}

func TestTaxValidation(t *testing.T) {
	t.Run("missing tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		ad := tax.AddonForKey(sdi.V1)
		ad.Normalizer(inv)
		err := ad.Validator(inv)
		assert.ErrorContains(t, err, "tax: cannot be blank.")
	})
}

func TestChargesValidation(t *testing.T) {
	ad := tax.AddonForKey(sdi.V1)

	t.Run("charge with no key", func(t *testing.T) {
		c := &bill.Charge{
			Percent: num.NewPercentage(10, 2),
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.NoError(t, err)
	})

	t.Run("fund contribution charge missing extension", func(t *testing.T) {
		c := &bill.Charge{
			Key:     sdi.KeyFundContribution,
			Percent: num.NewPercentage(10, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
				},
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.ErrorContains(t, err, "ext: (it-sdi-fund-type: required.)")
	})

	t.Run("fund contribution charge with valid extension", func(t *testing.T) {
		c := &bill.Charge{
			Key:     sdi.KeyFundContribution,
			Percent: num.NewPercentage(10, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "exempt",
				},
			},
			Ext: tax.Extensions{
				sdi.ExtKeyFundType: "TC04",
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.NoError(t, err)
	})

	t.Run("nil charge", func(t *testing.T) {
		var c *bill.Charge
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.NoError(t, err)
	})

	t.Run("fund contribution charge with missing taxes", func(t *testing.T) {
		c := &bill.Charge{
			Key:     sdi.KeyFundContribution,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				sdi.ExtKeyFundType: "TC04",
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.ErrorContains(t, err, "missing category VAT.")
	})

	t.Run("fund contribution charge with missing percentage", func(t *testing.T) {
		c := &bill.Charge{
			Key:    sdi.KeyFundContribution,
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
				},
			},
			Ext: tax.Extensions{
				sdi.ExtKeyFundType: "TC04",
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.ErrorContains(t, err, "percent: cannot be blank")
	})
}

func TestPaymentValidation(t *testing.T) {

	t.Run("payment advances", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
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
		inv.Payment = &bill.PaymentDetails{
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
		assert.ErrorContains(t, err, "payment: (instructions: cannot be blank when terms with due dates are present.).")
	})

	t.Run("payment terms with no due dates", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{
				Key: "instant",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("payment terms with instructions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
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

func TestAddressesValidation(t *testing.T) {
	t.Run("missing addresses", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = nil
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "supplier: (addresses: cannot be blank.)")
		assert.Contains(t, err.Error(), "customer: (addresses: cannot be blank.)")
	})

	t.Run("missing country", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Country = ""
		inv.Customer.Addresses[0].Country = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: (0: (country: cannot be blank.).).)")
		assert.ErrorContains(t, err, "customer: (addresses: (0: (country: cannot be blank.).).)")
	})

	t.Run("invalid code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Code = "123456"
		inv.Supplier.Addresses[0].Code = "123456"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: (0: (code: must be in a valid format.).).)")
		assert.ErrorContains(t, err, "customer: (addresses: (0: (code: must be in a valid format.).).)")
	})

	t.Run("codes in foreign country", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Country = "AT"
		inv.Supplier.Addresses[0].Country = "AT"
		inv.Customer.Addresses[0].Code = "1234"
		inv.Supplier.Addresses[0].Code = "1234"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid Latin-1 address fields", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with valid Latin-1 characters in address
		inv.Supplier.Addresses[0].Street = "Via dell'Università ñ°"
		inv.Supplier.Addresses[0].Locality = "Città di Castello àèìòù"
		inv.Customer.Addresses[0].Street = "Rue de la Paix é"
		inv.Customer.Addresses[0].Locality = "Saint-Étienne ç"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("invalid supplier address street with emoji", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Street = "Via Test 🏠"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: (0: (street: contains characters outside of Latin and Latin-1 range.).).)")
	})

	t.Run("invalid supplier postbox  with emoji", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Street = ""
		inv.Supplier.Addresses[0].PostOfficeBox = "post 🏠"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "contains characters outside of Latin and Latin-1 range")
	})

	t.Run("missing supplier address street and postbox", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Street = ""
		inv.Supplier.Addresses[0].PostOfficeBox = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "either street or post office box must be set")
	})

	t.Run("invalid customer address street with Japanese characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Street = "テスト通り"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (addresses: (0: (street: contains characters outside of Latin and Latin-1 range.).).)")
	})

	t.Run("multiple address validation errors", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Set multiple invalid fields to test comprehensive validation
		inv.Supplier.Name = "Test 中文"
		inv.Customer.Name = "Cliente 🎉"
		inv.Supplier.Addresses[0].Street = "Via Test 🏠"
		inv.Customer.Addresses[0].Locality = "Città 한국어"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()

		// Should contain multiple validation errors for Latin-1 violations
		assert.ErrorContains(t, err, "contains characters outside of Latin and Latin-1 range")

		// Check that all invalid fields are mentioned in the error
		errStr := err.Error()
		assert.Contains(t, errStr, "supplier")
		assert.Contains(t, errStr, "customer")
		assert.Contains(t, errStr, "name")
		assert.Contains(t, errStr, "addresses")
		assert.Contains(t, errStr, "street")
		assert.Contains(t, errStr, "locality")
	})
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

func TestInvoiceLineValidation(t *testing.T) {
	ad := tax.AddonForKey(sdi.V1)
	t.Run("missing item tax addon", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.NewAmount(10000, 2),
			},
			// No taxes!
		})
		ad.Normalizer(inv)
		err := ad.Validator(inv)
		require.ErrorContains(t, err, "lines: (1: (taxes: missing category VAT.).).")
	})

	t.Run("invalid item tax category", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.NewAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "GST",
					Key:      "standard",
				},
			},
		})
		ad.Normalizer(inv)
		err := ad.Validator(inv)
		require.ErrorContains(t, err, "lines: (1: (taxes: missing category VAT.).).")
	})

	t.Run("missing line", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{nil}
		ad.Normalizer(inv)
		require.NoError(t, ad.Validator(inv))
	})

	t.Run("missing line item", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item = nil
		ad.Normalizer(inv)
		require.NoError(t, ad.Validator(inv))
	})

	t.Run("with invalid item name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item.Name = "Test Item €"
		ad.Normalizer(inv)
		err := ad.Validator(inv)
		require.ErrorContains(t, err, "name: contains characters outside of Latin and Latin-1 range.")
	})
}

func TestOrderingValidation(t *testing.T) {
	t.Run("despatch without deferred tag", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				{
					Code:      "12345",
					IssueDate: cal.NewDate(2022, 1, 1),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "can only be set when invoice has deferred tag.")
	})

	t.Run("despatch with deferred tag and valid data", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(sdi.TagDeferred)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				{
					Code:      "12345",
					IssueDate: cal.NewDate(2022, 1, 1),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("despatch with deferred tag and valid additional data", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(sdi.TagDeferred)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				{
					Code:      "12345",
					IssueDate: cal.NewDate(2022, 1, 1),
					Reason:    "Partial shipment",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("despatch with deferred tag but missing code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(sdi.TagDeferred)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				{
					IssueDate: cal.NewDate(2022, 1, 1),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "code: cannot be blank.")
	})

	t.Run("despatch with deferred tag but missing issue date", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(sdi.TagDeferred)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				{
					Code: "12345",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "issue_date: cannot be blank.")
	})

	t.Run("multiple despatch documents with deferred tag", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(sdi.TagDeferred)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				{
					Code:      "12345",
					IssueDate: cal.NewDate(2022, 1, 1),
				},
				{
					Code:      "67890",
					IssueDate: cal.NewDate(2022, 1, 2),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("multiple despatch with one invalid", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(sdi.TagDeferred)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				{
					Code:      "12345",
					IssueDate: cal.NewDate(2022, 1, 1),
				},
				{
					Code: "67890",
					// Missing IssueDate
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "issue_date: cannot be blank.")
	})

	t.Run("nil despatch document", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(sdi.TagDeferred)
		inv.Ordering = &bill.Ordering{
			Despatch: []*org.DocumentRef{
				nil,
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("ordering without despatch", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{
			Code: "ORDER-123",
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("no ordering", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}
