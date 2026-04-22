package cfdi_test

import (
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	_ "github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withAddonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(cfdi.V4), tax.AddonForKey(cfdi.V4))
	}
}

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
						Rate:     "general",
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
		require.NoError(t, rules.Validate(inv))
	})
	t.Run("with global period", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyGlobalPeriod: "04",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.ErrorContains(t, err, "extensions must all be present or all absent")
	})
	t.Run("with global month", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyGlobalMonth: "02",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.ErrorContains(t, err, "extensions must all be present or all absent")
	})
	t.Run("with global year", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				cfdi.ExtKeyGlobalYear: "2025",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.ErrorContains(t, err, "extensions must all be present or all absent")
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
		require.ErrorContains(t, rules.Validate(inv), "extensions must all be present or all absent")
	})
}

func TestInvoiceCurrencyValidation(t *testing.T) {
	t.Run("non-MXN currency without exchange rates", func(t *testing.T) {
		inv := validInvoice()
		inv.Currency = "USD"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-MX-CFDI-V4-BILL-INVOICE-27] invoice must be in MXN or provide exchange rate for conversion")
	})

	t.Run("non-MXN currency with exchange rates", func(t *testing.T) {
		inv := validInvoice()
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   "USD",
				To:     "MXN",
				Amount: num.MakeAmount(1800, 2),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestNormalizeInvoice(t *testing.T) {
	t.Run("no tax", func(t *testing.T) {
		inv := validInvoice()
		inv.Addons = tax.WithAddons(cfdi.V4)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
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
		require.NoError(t, rules.Validate(inv))
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
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "must be set with global tag")
		assert.ErrorContains(t, err, "global invoices require")
		assert.ErrorContains(t, err, "payment is required for global invoices")
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
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, cbc.Code("04"), inv.Tax.Ext[cfdi.ExtKeyGlobalPeriod])
	})

}

func TestCustomerValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.TaxID = nil
	assertValidationError(t, inv, "customer tax ID is required")

	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.NoError(t, rules.Validate(inv))
}

func TestCustomerAddressCodeValidation(t *testing.T) {
	inv := validInvoice()
	delete(inv.Customer.Ext, "mx-cfdi-post-code")
	assertValidationError(t, inv, "Mexican customer must have at least one address")

	inv.Customer.Addresses = []*org.Address{{}}
	assertValidationError(t, inv, "customer address postal code is required")

	inv.Customer.Addresses[0].Code = "ABC"
	assertValidationError(t, inv, "customer address postal code format is invalid")

	inv.Customer.Addresses[0].Code = "21000"
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))

	inv.Customer.TaxID.Country = "US"
	inv.Customer.Addresses = nil
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestLineValidation(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Quantity = num.MakeAmount(0, 0)
	assertValidationError(t, inv, "quantity must be greater than 0")

	inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "quantity must be greater than 0")

	inv = validInvoice()

	inv.Lines[0].Item.Price = num.NewAmount(-1, 0)
	// negative price is normalized to negative quantity during Calculate()
	assertValidationError(t, inv, "quantity must be greater than 0")
}

func TestPaymentInstructionsValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.PaymentDetails{
		Instructions: &pay.Instructions{},
	}

	inv.Payment.Instructions.Key = "direct-debit"
	assertValidationError(t, inv, "payment instructions require 'mx-cfdi-payment-means' extension")

	inv.Payment.Instructions.Key = "unexisting"
	assertValidationError(t, inv, "key must be valid")
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
	assertValidationError(t, inv, "payment advance requires 'mx-cfdi-payment-means' extension")

	inv.Payment.Advances[0].Key = "unexisting"
	assertValidationError(t, inv, "key must be valid")

	inv.Payment.Advances[0].Key = ""
	assertValidationError(t, inv, "payment advance requires 'mx-cfdi-payment-means' extension")
}

func TestPaymentTermsValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.PaymentDetails{
		Terms: &pay.Terms{},
	}

	inv.Payment.Terms.Notes = strings.Repeat("x", 1001)
	assertValidationError(t, inv, "notes length must be no more than 1000")

	inv.Payment.Terms.Notes = strings.Repeat("x", 1000)
	require.NoError(t, rules.Validate(inv))
}

func TestUsoCFDIScenarioValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.Ext = tax.Extensions{
		cfdi.ExtKeyFiscalRegime: "601",
		"mx-cfdi-post-code":     "21000",
	}
	assertValidationError(t, inv, "Mexican customer requires 'mx-cfdi-fiscal-regime' and 'mx-cfdi-use' extensions")
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
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing 'sat-uuid' stamp")
	assert.Contains(t, err.Error(), "mx-cfdi-rel-type")

	inv.Type = bill.InvoiceTypeCreditNote
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing 'sat-uuid' stamp")

	inv.Preceding[0].Stamps[0].Provider = "sat-uuid"
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceDiscountValidation(t *testing.T) {
	inv := validInvoice()

	inv.Discounts = []*bill.Discount{
		{
			Percent: num.NewPercentage(20, 2),
		},
	}
	assertValidationError(t, inv, "not supported, use line discounts instead")
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
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
			err: "item price must be greater than 0",
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
			// negative price now normalized to quantity
			err: "quantity must be greater than 0",
		},
		{
			name: "nil price",
			item: &org.Item{
				Name: "Test purchase",
				Ext: tax.Extensions{
					cfdi.ExtKeyProdServ: "12345678",
				},
			},
			err: "item price is required",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(10000, 2),
			},
			err: "item requires 'mx-cfdi-prod-serv' extension",
		},
		{
			name: "empty extension",
			item: &org.Item{
				Name:  "Test purchase",
				Price: num.NewAmount(10000, 2),
				Ext:   tax.Extensions{},
			},
			err: "item requires 'mx-cfdi-prod-serv' extension",
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
			err: "item requires 'mx-cfdi-prod-serv' extension",
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
			err: "product/service code must have 8 digits",
		},
		{
			name: "nil",
			item: nil,
			err:  "item is required",
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
			err: "payment is required for global invoices",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			inv := validInvoice()
			inv.Tags = ts.tags
			inv.Lines[0].Item = ts.item
			require.NoError(t, inv.Calculate())
			err := rules.Validate(inv)
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
		require.NoError(t, rules.Validate(inv))
	})
	t.Run("missing ref", func(t *testing.T) {
		inv := validInvoiceGlobal()
		inv.Lines[0].Item = &org.Item{
			Price: num.NewAmount(10000, 2),
			Name:  "Test purchase",
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "must be set with global tag")
	})
}
