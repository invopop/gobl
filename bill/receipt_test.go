package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceiptCalculate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		require.NoError(t, r.Calculate())

		assert.Equal(t, bill.ReceiptTypePayment, r.Type)
		assert.Equal(t, currency.EUR, r.Currency)
		assert.Equal(t, r.Regime.Country.String(), "ES")
		assert.Equal(t, r.Supplier.TaxID.Code.String(), "B98602642", "should normalize")
	})

	t.Run("missing supplier", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		r.Supplier = nil
		assert.NotPanics(t, func() {
			r.Calculate()
		})
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		r.Supplier.TaxID = nil
		assert.NotPanics(t, func() {
			r.Calculate()
		})
	})

	t.Run("with debits and credits", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		r.Lines = append(r.Lines, &bill.ReceiptLine{
			Credit: num.NewAmount(5000, 2),
			Document: &org.DocumentRef{
				Type:      "credit-note",
				Series:    "CN1",
				Code:      "0123",
				IssueDate: cal.NewDate(2025, 1, 24),
			},
		})
		require.NoError(t, r.Calculate())

		assert.Equal(t, "50.00", r.Total.String(), "should balance")
	})

	t.Run("with credit", func(t *testing.T) {
		rct := testReceiptPaymentMinimal(t)
		rct.Lines[0].Credit = num.NewAmount(5000, 2)
		rct.Lines[0].Debit = nil
		require.NoError(t, rct.Calculate())

		assert.Equal(t, "-50.00", rct.Total.String(), "should balance")
	})

	t.Run("with taxes", func(t *testing.T) {
		r := testReceiptPaymentWithTax(t)
		require.NoError(t, r.Calculate())
		assert.Equal(t, "21.00", r.Tax.Sum.String())
	})

	t.Run("with multiple tax lines", func(t *testing.T) {
		r := testReceiptPaymentWithTax(t)
		r.Lines = append(r.Lines, &bill.ReceiptLine{
			Debit: num.NewAmount(10000, 2),
			Tax: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code: "VAT",
						Rates: []*tax.RateTotal{
							{
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(10, 2),
							},
						},
					},
				},
			},
		})
		require.NoError(t, r.Calculate())
		assert.Len(t, r.Tax.Categories, 1)
		assert.Len(t, r.Tax.Categories[0].Rates, 2)
		assert.Equal(t, "31.00", r.Tax.Sum.String())
	})
}

func TestReceiptValidate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		require.NoError(t, r.Calculate())
		require.NoError(t, r.Validate())
	})

	t.Run("with error", func(t *testing.T) {
		rct := testReceiptPaymentMinimal(t)
		require.NoError(t, rct.Calculate())
		rct.Supplier = nil
		assert.ErrorContains(t, rct.Validate(), "supplier: cannot be blank")
	})

	t.Run("with addon", func(t *testing.T) {
		rct := testReceiptPaymentMinimal(t)
		rct.Addons.SetAddons(tbai.V1)
		require.NoError(t, rct.Calculate())
		require.NoError(t, rct.Validate())
	})
}

func TestReceiptExchangeRates(t *testing.T) {
	t.Run("debit basic", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		r.Currency = currency.EUR
		r.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		r.Lines[0].Currency = currency.USD
		require.NoError(t, r.Calculate())

		assert.Equal(t, "96.00", r.Total.String())
	})

	t.Run("debit missing rate", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		r.Currency = currency.EUR
		r.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		r.Lines[0].Currency = currency.GBP
		require.ErrorContains(t, r.Calculate(), "lines: (0: (currency: no exchange rate found for GBP to EUR.).)")
	})

	t.Run("credit basic", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		r.Currency = currency.EUR
		r.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		r.Lines[0].Currency = currency.USD
		r.Lines[0].Credit = r.Lines[0].Debit
		r.Lines[0].Debit = nil
		require.NoError(t, r.Calculate())

		assert.Equal(t, "-96.00", r.Total.String())
	})

	t.Run("credit missing rate", func(t *testing.T) {
		r := testReceiptPaymentMinimal(t)
		r.Currency = currency.EUR
		r.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		r.Lines[0].Credit = r.Lines[0].Debit
		r.Lines[0].Debit = nil
		r.Lines[0].Currency = currency.GBP
		require.ErrorContains(t, r.Calculate(), "lines: (0: (currency: no exchange rate found for GBP to EUR.).).")
	})
}

func testReceiptPaymentMinimal(t *testing.T) *bill.Receipt {
	t.Helper()
	r := &bill.Receipt{
		Series:    "P1",
		Code:      "0123",
		IssueDate: cal.MakeDate(2025, 1, 24),
		Method: &pay.Instructions{
			Key: pay.MeansKeyCard,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B-98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		Lines: []*bill.ReceiptLine{
			{
				Document: &org.DocumentRef{
					Series:    "F1",
					Code:      "01234",
					IssueDate: cal.NewDate(2025, 1, 24),
				},
				Debit: num.NewAmount(10000, 2),
			},
		},
	}
	return r
}

func testReceiptPaymentWithTax(t *testing.T) *bill.Receipt {
	rct := testReceiptPaymentMinimal(t)
	rct.Lines[0].Tax = &tax.Total{
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
	}
	return rct
}

func TestReceiptJSONSchemaExtend(t *testing.T) {
	eg := `{
		"properties": {
			"$regime": {
				"$ref": "https://gobl.org/draft-0/cbc/key",
				"title": "Regime"
			},
			"$addons": {
				"items": {
            		"$ref": "https://gobl.org/draft-0/cbc/key",
					"type": "array",
					"title": "Addons",
					"description": "Addons defines a list of keys used to identify tax addons that apply special\nnormalization, scenarios, and validation rules to a document."
				}
			},
			"uuid": {
				"type": "string",
				"format": "uuid",
				"title": "UUID",
				"description": "Universally Unique Identifier."
			},
			"type": {
				"$ref": "https://gobl.org/draft-0/cbc/key",
				"title": "Type",
		        "description": "Type of invoice document subject to the requirements of the local tax regime.",
        		"calculated": true
			}
		}
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	rct := bill.Receipt{}
	rct.JSONSchemaExtend(js)

	assert.Equal(t, js.Properties.Len(), 4) // from this example

	t.Run("regime", func(t *testing.T) {
		prop, ok := js.Properties.Get("$regime")
		require.True(t, ok)
		assert.Greater(t, len(prop.OneOf), 1)
		rd := tax.AllRegimeDefs()[0]
		assert.Equal(t, rd.Code().String(), prop.OneOf[0].Const)
	})
	t.Run("addons", func(t *testing.T) {
		prop, ok := js.Properties.Get("$addons")
		require.True(t, ok)
		assert.Greater(t, len(prop.Items.OneOf), 1)
		ao := tax.AllAddonDefs()[0]
		assert.Equal(t, ao.Key.String(), prop.Items.OneOf[0].Const)
	})
	t.Run("types", func(t *testing.T) {
		prop, ok := js.Properties.Get("type")
		require.True(t, ok)
		assert.Greater(t, len(prop.OneOf), 1)
		it := bill.ReceiptTypes[0]
		assert.Equal(t, it.Key.String(), prop.OneOf[0].Const)
	})
	t.Run("recommended", func(t *testing.T) {
		assert.Len(t, js.Extras[schema.Recommended], 1)
	})

}
