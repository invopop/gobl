package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceCustomerVATStatusNormalization(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("customer without tax ID sets final consumer", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = nil
		inv.Customer.Ext = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.Extensions{
					arca.ExtKeyIdentityType: "96", // DNI
				},
			},
		}
		ad.Normalizer(inv)
		assert.Equal(t, "5", inv.Customer.Ext[arca.ExtKeyVATStatus].String())
	})

	t.Run("customer with AR tax ID sets registered VAT responsible", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.Ext = nil
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Customer.Ext[arca.ExtKeyVATStatus].String())
	})

	t.Run("customer with foreign tax ID sets foreign customer", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "US",
			Code:    "123456789",
		}
		inv.Customer.Ext = nil
		ad.Normalizer(inv)
		assert.Equal(t, "9", inv.Customer.Ext[arca.ExtKeyVATStatus].String())
	})

	t.Run("customer with existing VAT status is not overwritten", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.Ext = tax.Extensions{
			arca.ExtKeyVATStatus: "6", // Monotributo Responsible
		}
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Customer.Ext[arca.ExtKeyVATStatus].String())
	})

	t.Run("nil customer does not panic", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Customer = nil
		assert.NotPanics(t, func() {
			ad.Normalizer(inv)
		})
	})
}

func TestInvoiceConceptNormalization(t *testing.T) {
	t.Run("only goods sets concept to products", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1", inv.Tax.Ext[arca.ExtKeyConcept].String())
	})

	t.Run("only services sets concept to services", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2", inv.Tax.Ext[arca.ExtKeyConcept].String())
	})

	t.Run("default item key (empty) treated as services", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Item.Key is empty by default, treated as services
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2", inv.Tax.Ext[arca.ExtKeyConcept].String())
	})

	t.Run("mixed goods and services sets concept to products and services", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Service Item",
				Price: num.NewAmount(5000, 2),
				Key:   org.ItemKeyServices,
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
		})
		// For mixed goods/services, we need ordering and payment
		inv.Ordering = testOrdering()
		inv.Payment = testPayment()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3", inv.Tax.Ext[arca.ExtKeyConcept].String())
	})

	t.Run("nil item defaults to services", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item:     nil, // nil item defaults to services
		})
		// With goods + nil (service), we need ordering and payment for mixed concept
		inv.Ordering = testOrdering()
		inv.Payment = testPayment()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3", inv.Tax.Ext[arca.ExtKeyConcept].String()) // mixed goods and services
	})

	t.Run("only nil items sets concept to services", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     nil, // nil item defaults to services
			},
		}
		inv.Ordering = testOrdering()
		inv.Payment = testPayment()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2", inv.Tax.Ext[arca.ExtKeyConcept].String()) // services
	})

	t.Run("existing tax extensions are merged", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				arca.ExtKeyDocType: "001",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1", inv.Tax.Ext[arca.ExtKeyConcept].String())
		assert.Equal(t, "001", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("empty lines does not set concept", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = nil
		require.NoError(t, inv.Calculate())
		assert.Empty(t, inv.Tax.Ext[arca.ExtKeyConcept])
	})

	t.Run("nil tax is initialized", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, "1", inv.Tax.Ext[arca.ExtKeyConcept].String())
	})
}

func TestInvoiceSeriesValidation(t *testing.T) {
	t.Run("missing series fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = ""
		assertValidationError(t, inv, "series: cannot be blank")
	})

	t.Run("non-numeric series fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "ABC"
		assertValidationError(t, inv, "series: must be a number")
	})

	t.Run("series with leading zeros is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "00001"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("series below range fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "0"
		assertValidationError(t, inv, "series: must be between 1 and 99998")
	})

	t.Run("series above range fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "99999"
		assertValidationError(t, inv, "series: must be between 1 and 99998")
	})

	t.Run("series at lower bound is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "1"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("series at upper bound is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "99998"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("series with very large number overflows", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "9999999999999999999999999" // number too large for strconv
		assertValidationError(t, inv, "series: must be a number")
	})
}

func TestInvoiceTaxValidation(t *testing.T) {
	t.Run("valid standard invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing tax fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		inv.Tax = nil
		require.ErrorContains(t, inv.Validate(), "tax: cannot be blank")
	})

	t.Run("missing doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = nil
		require.ErrorContains(t, inv.Validate(), "ar-arca-doc-type: required")
	})
}

func TestInvoiceCustomerValidation(t *testing.T) {
	t.Run("customer required for invoice type A", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer = nil
		assertValidationError(t, inv, "customer: cannot be blank")
	})

	t.Run("customer not required for invoice type B (006)", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("customer not required for debit note type B (007)", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext[arca.ExtKeyDocType] = "007"
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("customer not required for credit note type B (008)", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext[arca.ExtKeyDocType] = "008"
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("customer with tax ID is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("customer with identity and ext is valid", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.Extensions{
					arca.ExtKeyIdentityType: "96", // DNI
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("customer without tax ID or identity ext fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		assertValidationError(t, inv, "customer: must have a tax_id, or an identity with ext 'ar-arca-identity-type'")
	})

	t.Run("customer with identity but no ext fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678", // No ext
			},
		}
		assertValidationError(t, inv, "customer: must have a tax_id, or an identity with ext 'ar-arca-identity-type'")
	})

	t.Run("customer with tax ID but missing code fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID.Code = ""
		assertValidationError(t, inv, "customer: (tax_id: (code: cannot be blank.).")
	})
}

func TestInvoiceServiceRequirements(t *testing.T) {
	t.Run("services require ordering with period", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Ordering = nil
		assertValidationError(t, inv, "ordering: cannot be blank")
	})

	t.Run("services require ordering period", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Ordering = &bill.Ordering{}
		assertValidationError(t, inv, "ordering: (period: cannot be blank.).")
	})

	t.Run("services with valid ordering passes", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("services require payment terms", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Payment = nil
		assertValidationError(t, inv, "payment: cannot be blank")
	})

	t.Run("services require payment terms object", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Payment = &bill.PaymentDetails{}
		assertValidationError(t, inv, "payment: (terms: cannot be blank.).")
	})

	t.Run("services require payment due dates", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{},
		}
		assertValidationError(t, inv, "payment: (terms: (due_dates: cannot be blank.).).")
	})

	t.Run("products do not require ordering or payment", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Ordering = nil
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("mixed goods and services requires ordering and payment", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Service Item",
				Price: num.NewAmount(5000, 2),
				Key:   org.ItemKeyServices,
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
		})
		inv.Ordering = nil
		assertValidationError(t, inv, "ordering: cannot be blank")
	})
}

func TestCreditNoteValidation(t *testing.T) {
	t.Run("valid credit note type A", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Equal(t, "003", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("valid credit note type B", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Equal(t, "008", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})
}

func TestDebitNoteValidation(t *testing.T) {
	t.Run("valid debit note type A", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Equal(t, "002", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("valid debit note type B", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeDebitNote
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Equal(t, "007", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})
}

func TestValidateFunctionsWithNilValues(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("validate with nil tax", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax = nil
		// Call validator directly
		err := ad.Validator(inv)
		assert.ErrorContains(t, err, "tax: cannot be blank")
	})

	t.Run("validate ordering with nil for services", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		// Need to set concept to "2" (services) to trigger ordering validation
		inv.Tax.Ext[arca.ExtKeyConcept] = "2"
		inv.Ordering = nil
		err := ad.Validator(inv)
		assert.ErrorContains(t, err, "ordering: cannot be blank")
	})

	t.Run("validate payment with nil for services", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		// Need to set concept to "2" (services) to trigger payment validation
		inv.Tax.Ext[arca.ExtKeyConcept] = "2"
		inv.Payment = nil
		err := ad.Validator(inv)
		assert.ErrorContains(t, err, "payment: cannot be blank")
	})
}

func TestCorrectionDefinitions(t *testing.T) {
	t.Run("correction definitions exist for credit and debit notes", func(t *testing.T) {
		ad := tax.AddonForKey(arca.V4)
		require.NotNil(t, ad.Corrections)
		// Check that invoice correction definition exists
		def := ad.Corrections.Def(bill.ShortSchemaInvoice)
		require.NotNil(t, def)
		assert.True(t, def.HasType(bill.InvoiceTypeCreditNote))
		assert.True(t, def.HasType(bill.InvoiceTypeDebitNote))
		assert.True(t, def.HasExtension(arca.ExtKeyDocType))
	})
}

// Helper functions

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.ErrorContains(t, err, expected)
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Addons: tax.WithAddons(arca.V4),
		Series: "1",
		Code:   "123",
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				arca.ExtKeyDocType: "001",
			},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30500010912", // Valid company CUIT
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "20172543597", // Valid individual CUIL
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
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

func testInvoiceWithGoods(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceStandard(t)
	inv.Lines[0].Item.Key = org.ItemKeyGoods
	return inv
}

func testInvoiceWithServices(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceStandard(t)
	inv.Lines[0].Item.Key = org.ItemKeyServices
	inv.Ordering = testOrdering()
	inv.Payment = testPayment()
	return inv
}

func testInvoiceSimplified(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceWithGoods(t)
	inv.SetTags(tax.TagSimplified)
	inv.Tax.Ext[arca.ExtKeyDocType] = "006"
	return inv
}

func testOrdering() *bill.Ordering {
	return &bill.Ordering{
		Period: &cal.Period{
			Start: cal.MakeDate(2024, 1, 1),
			End:   cal.MakeDate(2024, 1, 31),
		},
	}
}

func testPayment() *bill.PaymentDetails {
	return &bill.PaymentDetails{
		Terms: &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date:   cal.NewDate(2024, 2, 15),
					Amount: num.MakeAmount(10000, 2),
				},
			},
		},
	}
}
