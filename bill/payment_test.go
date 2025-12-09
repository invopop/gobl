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

func TestPaymentCalculate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		p := testPaymentMinimal(t)
		require.NoError(t, p.Calculate())

		assert.Equal(t, bill.PaymentTypeReceipt, p.Type)
		assert.Equal(t, currency.EUR, p.Currency)
		assert.Equal(t, p.Regime.Country.String(), "ES")
		assert.Equal(t, p.Supplier.TaxID.Code.String(), "B98602642", "should normalize")
	})

	t.Run("missing supplier", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Supplier = nil
		assert.NotPanics(t, func() {
			assert.ErrorContains(t, p.Calculate(), "currency: required, unable to determine")
		})
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Supplier.TaxID = nil
		assert.NotPanics(t, func() {
			assert.ErrorContains(t, p.Calculate(), "currency: required, unable to determine")
		})
	})

	t.Run("with positive amount and refunds", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Lines = append(p.Lines, &bill.PaymentLine{
			Refund: true,
			Amount: num.MakeAmount(5000, 2),
			Document: &org.DocumentRef{
				Type:      "credit-note",
				Series:    "CN1",
				Code:      "0123",
				IssueDate: cal.NewDate(2025, 1, 24),
			},
		})
		require.NoError(t, p.Calculate())
		assert.Equal(t, "50.00", p.Total.String(), "should balance")
	})

	t.Run("with credit", func(t *testing.T) {
		pmt := testPaymentMinimal(t)
		pmt.Lines[0].Refund = true
		pmt.Lines[0].Amount = num.MakeAmount(5000, 2)
		require.NoError(t, pmt.Calculate())
		require.NoError(t, pmt.Validate())
		assert.Equal(t, "-50.00", pmt.Total.String(), "should balance")
	})

	t.Run("with taxes", func(t *testing.T) {
		p := testPaymentWithTax(t)
		require.NoError(t, p.Calculate())
		// assert.Equal(t, "21.00", p.Tax.Sum.String())
	})

	t.Run("with multiple tax lines", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.Lines = append(p.Lines, &bill.PaymentLine{
			Amount: num.MakeAmount(10000, 2),
			Document: &org.DocumentRef{
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
			},
		})
		require.NoError(t, p.Calculate())
		//assert.Len(t, p.Tax.Categories, 1)
		//assert.Len(t, p.Tax.Categories[0].Rates, 2)
		//assert.Equal(t, "31.00", p.Tax.Sum.String())
	})

	t.Run("with partial payments and taxes 50%", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.Lines[0].Payable = num.NewAmount(12100, 2)
		p.Lines[0].Amount = num.MakeAmount(6050, 2) // half paid
		require.NoError(t, p.Calculate())
		assert.Equal(t, "60.50", p.Total.String(), "should be half of the total")
		//assert.Equal(t, "10.50", p.Tax.Sum.String(), "should be half of the tax")
		//assert.Equal(t, "60.50", p.Due.String(), "should be half of the payable amount")
	})

	t.Run("with partial payments and taxes ~25%", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.Lines[0].Installment = 2
		p.Lines[0].Payable = num.NewAmount(12100, 2)
		p.Lines[0].Advances = num.NewAmount(2000, 2) // 20€ already paid
		p.Lines[0].Amount = num.MakeAmount(3025, 2)  // 25% paid, approx
		require.NoError(t, p.Calculate())
		assert.Equal(t, "30.25", p.Total.String(), "should be a quarter of the total")
		//assert.Equal(t, "5.25", p.Tax.Sum.String(), "should be a quarter of the tax")
		//assert.Equal(t, "70.75", p.Due.String(), "should be three quarters of the payable amount, minus 20€")
		assert.Equal(t, 2, p.Lines[0].Installment, "should be the second installment")
	})

	t.Run("missing lines", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Lines = nil
		require.NoError(t, p.Calculate())
		assert.Equal(t, num.AmountZero, p.Total)
	})

	t.Run("line indexes", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Lines = append(p.Lines, &bill.PaymentLine{
			Index: 23,
		})
		require.NoError(t, p.Calculate())
		assert.Equal(t, 1, p.Lines[0].Index)
		assert.Equal(t, 2, p.Lines[1].Index)
	})

	t.Run("without issue date", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.IssueDate = cal.Date{}
		require.NoError(t, p.Calculate())
		tn := cal.TodayIn(p.RegimeDef().TimeLocation())
		assert.Equal(t, p.IssueDate, tn)
		assert.Nil(t, p.IssueTime)
	})

	t.Run("with empty issue time", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.IssueDate = cal.Date{}
		p.IssueTime = new(cal.Time)
		require.NoError(t, p.Calculate())
		tn := cal.ThisSecondIn(p.RegimeDef().TimeLocation())
		assert.Equal(t, p.IssueDate.String(), tn.Date().String())
		assert.Equal(t, p.IssueTime.Hour, tn.Time().Hour)
		assert.Equal(t, p.IssueTime.Minute, tn.Time().Minute)
		assert.Equal(t, p.IssueTime.Second, tn.Time().Second)
	})

	t.Run("with different exchange rates", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.Currency = currency.USD
		p.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.EUR,
				To:     currency.USD,
				Amount: num.MakeAmount(110, 2),
			},
		}
		p.Lines[0].Document.Currency = currency.EUR
		require.NoError(t, p.Calculate())
		assert.Equal(t, "100.00", p.Total.String(), "should convert to USD")
		//assert.Equal(t, "23.10", p.Tax.Sum.String())
	})

	t.Run("with missing exchange rates", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.Currency = currency.MXN
		p.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.EUR,
				To:     currency.USD,
				Amount: num.MakeAmount(110, 2),
			},
		}
		p.Lines[0].Document.Payable = num.NewAmount(10000, 2)
		p.Lines[0].Document.Currency = currency.EUR
		require.ErrorContains(t, p.Calculate(), "lines: (1: (document: (currency: missing exchange rate from EUR to MXN.).).)")
	})

	t.Run("with multiple and different exchange rates", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.Currency = currency.USD
		p.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.EUR,
				To:     currency.USD,
				Amount: num.MakeAmount(110, 2),
			},
		}
		p.Lines = append(p.Lines, &bill.PaymentLine{
			Amount: num.MakeAmount(10000, 2),
			Document: &org.DocumentRef{
				Currency: currency.EUR,
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(20000, 2),
									Percent: num.NewPercentage(10, 2),
								},
							},
						},
					},
				},
			},
		})
		require.NoError(t, p.Calculate())
		data, _ := json.MarshalIndent(p, "", "  ")
		t.Logf("Payment JSON: %s", data)
		assert.Equal(t, "200.00", p.Total.String(), "should convert to USD")
		//assert.Equal(t, "43.00", p.Tax.Sum.String())
	})
}

func TestPaymentValidate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		p := testPaymentMinimal(t)
		require.NoError(t, p.Calculate())
		require.NoError(t, p.Validate())
	})

	t.Run("with error", func(t *testing.T) {
		pmt := testPaymentMinimal(t)
		require.NoError(t, pmt.Calculate())
		pmt.Supplier = nil
		assert.ErrorContains(t, pmt.Validate(), "supplier: cannot be blank")
	})

	t.Run("with addon", func(t *testing.T) {
		pmt := testPaymentMinimal(t)
		pmt.Addons.SetAddons(tbai.V1)
		require.NoError(t, pmt.Calculate())
		require.NoError(t, pmt.Validate())
	})

	t.Run("with nil array entries", func(t *testing.T) {
		pmt := testPaymentMinimal(t)
		pmt.Lines = append(pmt.Lines, nil)
		pmt.Notes = append(pmt.Notes, nil)
		pmt.Preceding = append(pmt.Preceding, nil)
		pmt.ExchangeRates = append(pmt.ExchangeRates, nil)
		pmt.Complements = append(pmt.Complements, nil)
		require.NoError(t, pmt.Calculate())
		err := pmt.Validate()
		assert.ErrorContains(t, err, "lines: (1: is required.)")
		assert.ErrorContains(t, err, "notes: (0: is required.)")
		assert.ErrorContains(t, err, "preceding: (0: is required.)")
		assert.ErrorContains(t, err, "exchange_rates: (0: is required.)")
		assert.ErrorContains(t, err, "complements: (0: is required.)")
	})
}

func testPaymentMinimal(t *testing.T) *bill.Payment {
	t.Helper()
	p := &bill.Payment{
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
		Lines: []*bill.PaymentLine{
			{
				Document: &org.DocumentRef{
					Series:    "F1",
					Code:      "01234",
					IssueDate: cal.NewDate(2025, 1, 24),
				},
				Amount: num.MakeAmount(10000, 2),
			},
		},
	}
	return p
}

func testPaymentWithTax(t *testing.T) *bill.Payment {
	pmt := testPaymentMinimal(t)
	pmt.Lines[0].Document.Tax = &tax.Total{
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
	return pmt
}

func TestPaymentJSONSchemaExtend(t *testing.T) {
	eg := `{
		"properties": {
			"$regime": {
				"$ref": "https://gobl.org/draft-0/tax/regime-code",
				"title": "Tax Regime"
			},
			"$addons": {
				"items": {
            		"$ref": "https://gobl.org/draft-0/tax/addon-list",
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

	pmt := bill.Payment{}
	pmt.JSONSchemaExtend(js)

	assert.Equal(t, js.Properties.Len(), 4) // from this example

	t.Run("types", func(t *testing.T) {
		prop, ok := js.Properties.Get("type")
		require.True(t, ok)
		assert.Greater(t, len(prop.OneOf), 1)
		it := bill.PaymentTypes[0]
		assert.Equal(t, it.Key.String(), prop.OneOf[0].Const)
	})
	t.Run("recommended", func(t *testing.T) {
		assert.Len(t, js.Extras[schema.Recommended], 4)
	})

}
