package sdi_test

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/rules"
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

func withSDIContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(sdi.V1), tax.AddonForKey(sdi.V1))
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {

		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
	t.Run("missing tax extensions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = nil
		err := rules.Validate(inv)
		require.ErrorContains(t, err, "tax requires 'it-sdi-document-type' and 'it-sdi-format' extensions")
	})

	t.Run("non-EUR currency without exchange rates", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Currency = "USD"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-IT-SDI-V1-BILL-INVOICE-22] invoice must be in EUR or provide exchange rate for conversion")
	})

	t.Run("non-EUR currency with exchange rates", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   "USD",
				To:     "EUR",
				Amount: num.MakeAmount(875967, 6),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestInvoiceNormalization(t *testing.T) {
	ad := tax.AddonForKey(sdi.V1)

	t.Run("supplier fiscal regime", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		ad.Normalizer(inv)
		assert.Equal(t, "RF01", inv.Supplier.Ext[sdi.ExtKeyFiscalRegime].String())
	})

	t.Run("strip +39 from italian supplier telephone", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "+39333123456"}}
		ad.Normalizer(inv)
		require.Len(t, inv.Supplier.Telephones, 1)
		assert.Equal(t, "333123456", inv.Supplier.Telephones[0].Number)
	})

	t.Run("non-italian supplier telephone not normalized", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Country = "FR"
		inv.Supplier.Telephones = []*org.Telephone{{Number: "+39333123456"}}
		ad.Normalizer(inv)
		require.Len(t, inv.Supplier.Telephones, 1)
		assert.Equal(t, "+39333123456", inv.Supplier.Telephones[0].Number)
	})

	t.Run("no telephones nothing happens", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = nil
		ad.Normalizer(inv)
		assert.Nil(t, inv.Supplier.Telephones)
	})

	t.Run("italian supplier telephone without +39 prefix not normalized", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "333123456"}}
		ad.Normalizer(inv)
		require.Len(t, inv.Supplier.Telephones, 1)
		assert.Equal(t, "333123456", inv.Supplier.Telephones[0].Number)
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
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("with supplier missing registration details", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Registration = &org.Registration{
			Entry: "123456",
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier registration office is required")
	})

	t.Run("with invalid tax ID code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "RSSGNN60R30H501U",
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Italian VAT identity code")
	})

	t.Run("missing supplier", func(t *testing.T) {
		// Verify normalizer doesn't panic with nil supplier
		inv := testInvoiceStandard(t)
		ad := tax.AddonForKey(sdi.V1)
		inv.Supplier = nil
		assert.NotPanics(t, func() {
			ad.Normalizer(inv)
		})
	})

	t.Run("valid Latin-1 supplier name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with valid Latin-1 characters including accented characters
		inv.Supplier.Name = "Società di Test SRL àáâãäåæçèéêë"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid supplier name with non-Latin-1 characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with emoji (outside Latin-1 range)
		inv.Supplier.Name = "Test Supplier 😊"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier name must use Latin-1 characters")
	})

	t.Run("invalid supplier name with Greek characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with Greek characters (outside Latin-1 range)
		inv.Supplier.Name = "Test Supplier αβγδε"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier name must use Latin-1 characters")
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
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("missing tax_id", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = append(inv.Customer.Identities, id)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "customer tax ID is required")
	})

	t.Run("missing tax id code and identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "",
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		// ensure contains both errors
		assert.ErrorContains(t, err, fmt.Sprintf("customer requires identity with key '%s'", it.IdentityKeyFiscalCode))
		assert.ErrorContains(t, err, "customer tax ID code is required")
	})

	t.Run("missing address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer addresses are required")
	})

	t.Run("missing customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer is required")
	})

	t.Run("valid Latin-1 customer name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with valid Latin-1 characters including special symbols
		inv.Customer.Name = "Cliente & Cia. S.p.A. ñöüß"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid customer name with Chinese characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Test with Chinese characters (outside Latin-1 range)
		inv.Customer.Name = "测试客户"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer name must use Latin-1 characters")
	})

	t.Run("missing customer name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer name is required")
	})

	t.Run("missing customer people with identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID.Code = ""
		inv.Customer.Name = ""
		inv.Customer.Identities = append(inv.Customer.Identities, id)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer name is required")
		assert.ErrorContains(t, err, "customer people are required when name is empty")
	})

}

func TestSupplierTelephoneValidation(t *testing.T) {
	t.Run("valid italian supplier telephone", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "A1B2C3"}}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid italian supplier telephone too short", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "1234"}}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "Italian telephone number length must be between 5 and 12")
	})

	t.Run("valid italian supplier telephone with symbols", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "+39333123456"}}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("valid italian number, because normalized", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "+393331234567"}}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
		assert.Equal(t, "3331234567", inv.Supplier.Telephones[0].Number)
	})

	t.Run("invalid italian supplier telephone too long", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "1233312345678"}}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "Italian telephone number length must be between 5 and 12")
	})

	t.Run("missing italian supplier telephones", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// No telephones set
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("non-italian supplier telephone not validated", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Country = "FR"
		inv.Supplier.Telephones = []*org.Telephone{{Number: "1234"}} // Too short, but should be ignored
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("italian supplier telephone too short without prefix", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = []*org.Telephone{{Number: "1234"}}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "Italian telephone number length must be between 5 and 12")
	})

	t.Run("no telephones nothing validated", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestTaxValidation(t *testing.T) {
	t.Run("missing tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax is required")
	})
}

func TestChargesValidation(t *testing.T) {
	t.Run("charge with no key", func(t *testing.T) {
		c := &bill.Charge{
			Percent: num.NewPercentage(10, 2),
		}
		err := rules.Validate(c, withSDIContext())
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
		err := rules.Validate(c, withSDIContext())
		assert.ErrorContains(t, err, fmt.Sprintf("fund contribution charge requires '%s' extension", sdi.ExtKeyFundType))
	})

	t.Run("fund contribution charge with valid extension", func(t *testing.T) {
		c := &bill.Charge{
			Key:     sdi.KeyFundContribution,
			Percent: num.NewPercentage(10, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
					Percent:  num.NewPercentage(22, 2),
				},
			},
			Ext: tax.Extensions{
				sdi.ExtKeyFundType: "TC04",
			},
		}
		err := rules.Validate(c, withSDIContext())
		assert.NoError(t, err)
	})

	t.Run("nil charge", func(t *testing.T) {
		var c *bill.Charge
		err := rules.Validate(c, withSDIContext())
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
		err := rules.Validate(c, withSDIContext())
		assert.ErrorContains(t, err, "fund contribution charge must have VAT tax category")
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
		err := rules.Validate(c, withSDIContext())
		assert.ErrorContains(t, err, "fund contribution charge requires a percentage")
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
		assert.NoError(t, rules.Validate(inv))
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
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment instructions are required when terms with due dates are present")
	})

	t.Run("payment terms with no due dates", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{
				Key: "instant",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
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
		assert.NoError(t, rules.Validate(inv))
		assert.Equal(t, "MP08", inv.Payment.Instructions.Ext[sdi.ExtKeyPaymentMeans].String())
	})

}

func TestAddressesValidation(t *testing.T) {
	t.Run("missing addresses", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = nil
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "supplier addresses are required")
		assert.Contains(t, err.Error(), "customer addresses are required")
	})

	t.Run("missing country", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Country = ""
		inv.Customer.Addresses[0].Country = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "country is required")
	})

	t.Run("invalid code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Code = "123456"
		inv.Supplier.Addresses[0].Code = "123456"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "Italian address code must be 5 digits")
	})

	t.Run("codes in foreign country", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Country = "AT"
		inv.Supplier.Addresses[0].Country = "AT"
		inv.Customer.Addresses[0].Code = "1234"
		inv.Supplier.Addresses[0].Code = "1234"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
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
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid supplier address street with emoji", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Street = "Via Test 🏠"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "street must use Latin-1 characters")
	})

	t.Run("invalid supplier postbox  with emoji", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Street = ""
		inv.Supplier.Addresses[0].PostOfficeBox = "post 🏠"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "post office box must use Latin-1 characters")
	})

	t.Run("missing supplier address street and postbox", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Street = ""
		inv.Supplier.Addresses[0].PostOfficeBox = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "either street or post office box must be set")
	})

	t.Run("invalid customer address street with Japanese characters", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Street = "テスト通り"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "street must use Latin-1 characters")
	})

	t.Run("multiple address validation errors", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Set multiple invalid fields to test comprehensive validation
		inv.Supplier.Name = "Test 中文"
		inv.Customer.Name = "Cliente 🎉"
		inv.Supplier.Addresses[0].Street = "Via Test 🏠"
		inv.Customer.Addresses[0].Locality = "Città 한국어"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)

		// Should contain multiple validation errors for Latin-1 violations
		assert.ErrorContains(t, err, "must use Latin-1 characters")

		// Check that relevant fields are mentioned in the error
		errStr := err.Error()
		assert.Contains(t, errStr, "supplier")
		assert.Contains(t, errStr, "customer")
	})
}

func TestRetainedTaxesValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
		Category: "IRPEF",
		Percent:  num.NewPercentage(20, 2),
	})
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), fmt.Sprintf("retained tax combo requires '%s' extension", sdi.ExtKeyRetained))
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
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceLineValidation(t *testing.T) {
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
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.ErrorContains(t, err, "line must have VAT tax category")
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
		// Cannot use inv.Calculate() here because GST is not a valid category in IT regime
		err := rules.Validate(inv, withSDIContext())
		require.ErrorContains(t, err, "line must have VAT tax category")
	})

	t.Run("missing line", func(t *testing.T) {
		// Nil lines should not cause a panic in SDI rules
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{nil}
		// Cannot use Calculate() with nil lines; validate directly
		err := rules.Validate(inv, withSDIContext())
		// SDI addon shouldn't add errors for nil lines
		if err != nil {
			assert.NotContains(t, err.Error(), "IT-SDI-V1-BILL-INVOICE-17")
		}
	})

	t.Run("missing line item", func(t *testing.T) {
		// Nil item should not cause a panic in SDI rules
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item = nil
		// Cannot use Calculate() with nil item; validate directly
		err := rules.Validate(inv, withSDIContext())
		// SDI addon shouldn't add item name errors for nil items
		if err != nil {
			assert.NotContains(t, err.Error(), "IT-SDI-V1-BILL-INVOICE-18")
		}
	})

	t.Run("with invalid item name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item.Name = "Test Item €"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.ErrorContains(t, err, "item name must use Latin-1 characters")
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
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "despatch can only be set when invoice has deferred tag")
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
		require.NoError(t, rules.Validate(inv))
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
		require.NoError(t, rules.Validate(inv))
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
		err := rules.Validate(inv)
		// Code validation is handled by base document reference rules
		assert.ErrorContains(t, err, "document reference code is required")
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
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "despatch issue date is required")
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
		require.NoError(t, rules.Validate(inv))
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
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "despatch issue date is required")
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
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("ordering without despatch", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{
			Code: "ORDER-123",
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("no ordering", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}
