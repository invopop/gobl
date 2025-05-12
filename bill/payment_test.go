package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
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

	t.Run("with debits and credits", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Lines = append(p.Lines, &bill.PaymentLine{
			Credit: num.NewAmount(5000, 2),
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
		pmt.Lines[0].Credit = num.NewAmount(5000, 2)
		pmt.Lines[0].Debit = nil
		require.NoError(t, pmt.Calculate())

		assert.Equal(t, "-50.00", pmt.Total.String(), "should balance")
	})

	t.Run("with taxes", func(t *testing.T) {
		p := testPaymentWithTax(t)
		require.NoError(t, p.Calculate())
		assert.Equal(t, "21.00", p.Tax.Sum.String())
	})

	t.Run("with multiple tax lines", func(t *testing.T) {
		p := testPaymentWithTax(t)
		p.Lines = append(p.Lines, &bill.PaymentLine{
			Debit: num.NewAmount(10000, 2),
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
		assert.Len(t, p.Tax.Categories, 1)
		assert.Len(t, p.Tax.Categories[0].Rates, 2)
		assert.Equal(t, "31.00", p.Tax.Sum.String())
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
}

func TestPaymentExchangeRates(t *testing.T) {
	t.Run("debit basic", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Currency = currency.EUR
		p.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		p.Lines[0].Currency = currency.USD
		require.NoError(t, p.Calculate())

		assert.Equal(t, "96.00", p.Total.String())
	})

	t.Run("debit missing rate", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Currency = currency.EUR
		p.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		p.Lines[0].Currency = currency.GBP
		require.ErrorContains(t, p.Calculate(), "lines: (0: (currency: no exchange rate found for GBP to EUR.).)")
	})

	t.Run("credit basic", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Currency = currency.EUR
		p.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		p.Lines[0].Currency = currency.USD
		p.Lines[0].Credit = p.Lines[0].Debit
		p.Lines[0].Debit = nil
		require.NoError(t, p.Calculate())

		assert.Equal(t, "-96.00", p.Total.String())
	})

	t.Run("credit missing rate", func(t *testing.T) {
		p := testPaymentMinimal(t)
		p.Currency = currency.EUR
		p.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(96, 2),
			},
		}
		p.Lines[0].Credit = p.Lines[0].Debit
		p.Lines[0].Debit = nil
		p.Lines[0].Currency = currency.GBP
		require.ErrorContains(t, p.Calculate(), "lines: (0: (currency: no exchange rate found for GBP to EUR.).).")
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
				Debit: num.NewAmount(10000, 2),
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

	pmt := bill.Payment{}
	pmt.JSONSchemaExtend(js)

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
		it := bill.PaymentTypes[0]
		assert.Equal(t, it.Key.String(), prop.OneOf[0].Const)
	})
	t.Run("recommended", func(t *testing.T) {
		assert.Len(t, js.Extras[schema.Recommended], 4)
	})

}

func TestPaymentBillable(t *testing.T) {
	pmt := testPaymentMinimal(t)
	pmt.IssueTime = cal.NewTime(10, 30, 0)
	pmt.ValueDate = cal.NewDate(2023, 1, 15)
	pmt.Currency = currency.EUR
	pmt.ExchangeRates = []*currency.ExchangeRate{{}}
	pmt.Complements = []*schema.Object{{}}

	assert.Equal(t, pmt.Series, pmt.GetSeries())
	assert.Equal(t, pmt.Code, pmt.GetCode())
	assert.Equal(t, pmt.IssueDate, pmt.GetIssueDate())
	assert.Equal(t, pmt.IssueTime, pmt.GetIssueTime())
	assert.Equal(t, pmt.ValueDate, pmt.GetValueDate())
	assert.Equal(t, pmt.Preceding, pmt.GetPreceding())
	assert.Equal(t, pmt.Currency, pmt.GetCurrency())
	assert.Equal(t, pmt.ExchangeRates, pmt.GetExchangeRates())
	assert.Equal(t, pmt.Supplier, pmt.GetSupplier())
	assert.Equal(t, pmt.Customer, pmt.GetCustomer())
	assert.Equal(t, pmt.Complements, pmt.GetComplements())
	assert.Nil(t, pmt.GetDiscounts())
	assert.Nil(t, pmt.GetCharges())
	assert.Nil(t, pmt.GetPaymentDetails())
	assert.Nil(t, pmt.GetTax())
	assert.Nil(t, pmt.GetTotals())
	assert.Nil(t, pmt.GetLines())

	pmt.SetCode(cbc.Code("002"))
	assert.Equal(t, cbc.Code("002"), pmt.Code)

	pmt.SetIssueDate(cal.MakeDate(2023, 2, 1))
	assert.Equal(t, cal.MakeDate(2023, 2, 1), pmt.IssueDate)

	pmt.SetIssueTime(cal.NewTime(11, 30, 0))
	assert.Equal(t, cal.NewTime(11, 30, 0), pmt.IssueTime)

	pmt.SetCurrency(currency.USD)
	assert.Equal(t, currency.USD, pmt.Currency)

	newTotals := &bill.Totals{}
	pmt.SetTotals(newTotals)
	// Nothing to assert here
}
