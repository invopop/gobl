package cfdi_test

import (
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	_ "github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Addons:    tax.WithAddons(cfdi.V4),
		Code:      "123",
		Currency:  "MXN",
		IssueDate: cal.MakeDate(2023, 1, 1),
		IssueTime: cal.NewTime(12, 34, 10),
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyIssuePlace: "21000",
			},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			Ext: tax.Extensions{
				"mx-cfdi-post-code":     "21000",
				cfdi.ExtKeyFiscalRegime: "601",
			},
			TaxID: &tax.Identity{
				Country: "MX",
				Code:    "AAA010101AAA",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Ext: tax.Extensions{
				"mx-cfdi-post-code":     "65000",
				cfdi.ExtKeyFiscalRegime: "608",
				cfdi.ExtKeyUse:          "G01",
			},
			TaxID: &tax.Identity{
				Country: "MX",
				Code:    "ZZZ010101ZZZ",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
					Ext: tax.Extensions{
						cfdi.ExtKeyProdServ: "01010101",
					},
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

func validInvoiceGlobal() *bill.Invoice {
	inv := validInvoice()
	inv.Tags = tax.WithTags(cfdi.TagGlobal)
	inv.Lines[0].Item.Ref = "TEST1234"
	inv.Tax.Ext = inv.Tax.Ext.Merge(tax.Extensions{
		cfdi.ExtKeyGlobalPeriod: "04",
		cfdi.ExtKeyGlobalMonth:  "01",
		cfdi.ExtKeyGlobalYear:   "2025",
	})
	inv.Payment = &bill.PaymentDetails{
		Advances: []*pay.Advance{
			{
				Key:         pay.MeansKeyCash,
				Description: "Prepaid",
				Percent:     num.NewPercentage(100, 2),
			},
		},
	}
	return inv
}

func TestValidInvoice(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
	t.Run("with global period", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyGlobalPeriod: "04",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		// Order is not guaranteed, so check for each error separately
		require.ErrorContains(t, err, "mx-cfdi-global-month: required")
		require.ErrorContains(t, err, "mx-cfdi-global-year: required")
	})
	t.Run("with global month", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyGlobalMonth: "02",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.ErrorContains(t, err, "mx-cfdi-global-period: required")
		require.ErrorContains(t, err, "mx-cfdi-global-year: required")
	})
	t.Run("with global year", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyGlobalYear: "2025",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.ErrorContains(t, err, "mx-cfdi-global-month: required")
		require.ErrorContains(t, err, "mx-cfdi-global-period: required")
	})
	t.Run("with global period and month", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyGlobalPeriod: "04",
				cfdi.ExtKeyGlobalMonth:  "02",
			},
		}
		require.NoError(t, inv.Calculate())
		require.ErrorContains(t, inv.Validate(), "tax: (ext: (mx-cfdi-global-year: required.).)")
	})
}

func TestNormalizeInvoice(t *testing.T) {
	t.Run("no tax", func(t *testing.T) {
		inv := validInvoice()
		inv.Addons = tax.WithAddons(cfdi.V4)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, cbc.Code("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
	t.Run("with supplier address code", func(t *testing.T) {
		inv := validInvoice()
		inv.Addons = tax.WithAddons(cfdi.V4)
		delete(inv.Supplier.Ext, "mx-cfdi-post-code")
		inv.Supplier.Addresses = append(inv.Supplier.Addresses,
			&org.Address{
				Locality: "Mexico",
				Code:     "21000",
			},
		)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, cbc.Code("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
	t.Run("with global tag, invalid", func(t *testing.T) {
		inv := validInvoice()
		inv.Tags = tax.WithTags(cfdi.TagGlobal)
		require.NoError(t, inv.Calculate())
		require.Nil(t, inv.Customer)
	})
	t.Run("should set time and date", func(t *testing.T) {
		// These tests can fail very rarely if run on the exact transition of the milliseconds
		tz, err := time.LoadLocation("America/Mexico_City")
		require.NoError(t, err)
		inv := validInvoice()
		inv.IssueTime = nil
		tn := time.Now().In(tz)
		require.NoError(t, inv.Calculate())
		assert.NotNil(t, inv.IssueTime)
		assert.Equal(t, tn.Format("2006-01-02"), inv.IssueDate.String())
		assert.Equal(t, tn.Format("15:04:05"), inv.IssueTime.String())
	})
}

func TestInvoiceGlobalTagValidation(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		inv := validInvoice()
		inv.Tags = tax.WithTags(cfdi.TagGlobal)
		require.NoError(t, inv.Calculate())
		require.Nil(t, inv.Customer)
		err := inv.Validate()
		assert.ErrorContains(t, err, "lines: (0: (item: (ref: must be set with global tag.).).)")
		assert.ErrorContains(t, err, "tax: (ext: (mx-cfdi-global-month: required; mx-cfdi-global-period: required; mx-cfdi-global-year: required.).)")
		assert.ErrorContains(t, err, "payment: cannot be blank;")
	})
	t.Run("success", func(t *testing.T) {
		inv := validInvoice()
		inv.Tags = tax.WithTags(cfdi.TagGlobal)
		inv.Lines[0].Item.Ref = "TEST1234"
		inv.Tax.Ext = inv.Tax.Ext.Merge(tax.Extensions{
			cfdi.ExtKeyGlobalPeriod: "04",
			cfdi.ExtKeyGlobalMonth:  "01",
			cfdi.ExtKeyGlobalYear:   "2025",
		})
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Key:         pay.MeansKeyCash,
					Description: "Prepaid",
					Percent:     num.NewPercentage(100, 2),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.Nil(t, inv.Customer)
		require.NoError(t, inv.Validate())
		assert.Equal(t, cbc.Code("04"), inv.Tax.Ext[cfdi.ExtKeyGlobalPeriod])
	})

}

func TestCustomerValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.TaxID = nil
	assertValidationError(t, inv, "customer: (tax_id: cannot be blank.)")

	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())
}

func TestCustomerAddressCodeValidation(t *testing.T) {
	inv := validInvoice()
	delete(inv.Customer.Ext, "mx-cfdi-post-code")
	assertValidationError(t, inv, "customer: (addresses: cannot be blank.)")

	inv.Customer.Addresses = []*org.Address{{}}
	assertValidationError(t, inv, "customer: (addresses: (0: (code: cannot be blank.).).)")

	inv.Customer.Addresses[0].Code = "ABC"
	assertValidationError(t, inv, "customer: (addresses: (0: (code: must be in a valid format.).).)")

	inv.Customer.Addresses[0].Code = "21000"
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	inv.Customer.TaxID.Country = "US"
	inv.Customer.Addresses = nil
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestLineValidation(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Quantity = num.MakeAmount(0, 0)
	assertValidationError(t, inv, "lines: (0: (quantity: must be greater than 0.).)")

	inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (quantity: must be greater than 0; total: must be no less than 0.).)")

	inv = validInvoice()

	inv.Lines[0].Item.Price = num.NewAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (item: (price: must be greater than 0.); total: must be no less than 0.).)")
}

func TestPaymentInstructionsValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{},
	}

	inv.Payment.Instructions.Key = "direct-debit"
	assertValidationError(t, inv, "payment: (instructions: (ext: (mx-cfdi-payment-means: required.).).)")

	inv.Payment.Instructions.Key = "unexisting"
	assertValidationError(t, inv, "payment: (instructions: (key: must be or start with a valid key.).)")
}

func TestPaymentAdvancesValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.PaymentDetails{
		Advances: []*pay.Advance{
			{
				Description: "A prepayment",
			},
		},
	}

	inv.Payment.Advances[0].Key = "direct-debit"
	assertValidationError(t, inv, "payment: (advances: (0: (ext: (mx-cfdi-payment-means: required.).).).)")

	inv.Payment.Advances[0].Key = "unexisting"
	assertValidationError(t, inv, "payment: (advances: (0: (key: must be or start with a valid key.).).)")

	inv.Payment.Advances[0].Key = ""
	assertValidationError(t, inv, "payment: (advances: (0: (ext: (mx-cfdi-payment-means: required.).).).)")
}

func TestPaymentTermsValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.PaymentDetails{
		Terms: &pay.Terms{},
	}

	inv.Payment.Terms.Notes = strings.Repeat("x", 1001)
	assertValidationError(t, inv, "payment: (terms: (notes: the length must be no more than 1000.).)")

	inv.Payment.Terms.Notes = strings.Repeat("x", 1000)
	require.NoError(t, inv.Validate())
}

func TestUsoCFDIScenarioValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.Ext = tax.Extensions{
		cfdi.ExtKeyFiscalRegime: "601",
		"mx-cfdi-post-code":     "21000",
	}
	assertValidationError(t, inv, "ext: (mx-cfdi-use: required.)")
}

func TestPrecedingValidation(t *testing.T) {
	inv := validInvoice()

	inv.Preceding = []*org.DocumentRef{
		{
			Code: "123",
			Stamps: []*head.Stamp{
				{
					Provider: "unexpected",
					Value:    "1234",
				},
			},
		},
	}
	assertValidationError(t, inv, "preceding: (0: (stamps: missing sat-uuid stamp.).); tax: (ext: (mx-cfdi-rel-type: required.).)")

	inv.Type = bill.InvoiceTypeCreditNote
	assertValidationError(t, inv, "preceding: (0: (stamps: missing sat-uuid stamp.).)")

	inv.Preceding[0].Stamps[0].Provider = "sat-uuid"
	require.NoError(t, inv.Validate())
}

func TestInvoiceDiscountValidation(t *testing.T) {
	inv := validInvoice()

	inv.Discounts = []*bill.Discount{
		{
			Percent: num.NewPercentage(20, 2),
		},
	}
	assertValidationError(t, inv, "discounts: not supported, use line discounts instead")
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}

func TestInvoiceLineItemValidation(t *testing.T) {
	tests := []struct {
		name string
		tags tax.Tags
		item *org.Item
		err  string
	}{
		{
			name: "valid item",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(10000, 2),
				Ext: tax.Extensions{
					cfdi.ExtKeyProdServ: "12345678",
				},
			},
		},
		{
			name: "zero price",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(0, 2),
				Ext: tax.Extensions{
					cfdi.ExtKeyProdServ: "12345678",
				},
			},
			err: "price: must be greater than 0",
		},
		{
			name: "negative price",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(-5000, 2),
				Ext: tax.Extensions{
					cfdi.ExtKeyProdServ: "12345678",
				},
			},
			err: "price: must be greater than 0",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(10000, 2),
			},
			err: "ext: (mx-cfdi-prod-serv: required.)",
		},
		{
			name: "empty extension",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(10000, 2),
				Ext:   tax.Extensions{},
			},
			err: "ext: (mx-cfdi-prod-serv: required.)",
		},
		{
			name: "invalid extension key",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(10000, 2),
				Ext: tax.Extensions{
					"random": "12345678",
				},
			},
			err: "ext: (random: undefined.)",
		},
		{
			name: "invalid code format",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(10000, 2),
				Ext: tax.Extensions{
					cfdi.ExtKeyProdServ: "AbC2",
				},
			},
			err: "ext: (mx-cfdi-prod-serv: must have 8 digits.)",
		},
		{
			name: "nil",
			item: nil,
			err:  "lines: (0: (item: cannot be blank.).",
		},
		{
			// see below for specific global tag tests
			name: "with global tag",
			tags: tax.WithTags(cfdi.TagGlobal),
			item: &org.Item{
				Ref:   "TEST1234",
				Price: num.NewAmount(10000, 2),
				Name:  "Test purchase",
			},
			err: "payment: cannot be blank; tax: (ext: (mx-cfdi-global-month: required; mx-cfdi-global-period: required; mx-cfdi-global-year: required.).)",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			inv := validInvoice()
			inv.Tags = ts.tags
			inv.Lines[0].Item = ts.item
			require.NoError(t, inv.Calculate())
			err := inv.Validate()
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}

func TestInvoiceLineItemGlobalValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		inv := validInvoiceGlobal()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
	t.Run("missing ref", func(t *testing.T) {
		inv := validInvoiceGlobal()
		inv.Lines[0].Item = &org.Item{
			Price: num.NewAmount(10000, 2),
			Name:  "Test purchase",
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "lines: (0: (item: (ref: must be set with global tag.).).)")
	})
}
