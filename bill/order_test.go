package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderCalculate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ord := baseOrderWithLines(t)
		require.NoError(t, ord.Calculate())
		assert.Nil(t, ord.Totals)
	})
}

func TestOrderValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ord := baseOrderWithLines(t)
		require.NoError(t, ord.Calculate())
		require.NoError(t, ord.Validate())
	})
}

func TestOrderConvertInto(t *testing.T) {
	t.Run("no prices", func(t *testing.T) {
		ord := baseOrderWithLines(t)
		ord.Currency = currency.USD
		ord.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		o2, err := ord.ConvertInto(currency.EUR)
		assert.NoError(t, err)
		assert.Nil(t, o2.Totals)
	})
	t.Run("basic", func(t *testing.T) {
		ord := baseOrderWithLines(t)
		ord.Currency = currency.USD
		ord.Lines[0].Item.Price = num.NewAmount(1000, 2)
		ord.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		o2, err := ord.ConvertInto(currency.EUR)
		assert.NoError(t, err)
		assert.Equal(t, currency.EUR, o2.Currency)
		assert.Equal(t, "87.5970", o2.Lines[0].Total.String())
	})
	t.Run("same currency", func(t *testing.T) {
		ord := baseOrderWithLines(t)
		ord.Currency = currency.USD
		ord.Lines[0].Item.Price = num.NewAmount(1000, 2)
		ord.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		o2, err := ord.ConvertInto(currency.USD)
		assert.NoError(t, err)
		assert.Equal(t, currency.USD, o2.Currency)
		assert.Equal(t, "100.00", o2.Lines[0].Total.String())
	})
	t.Run("missing rate", func(t *testing.T) {
		ord := baseOrderWithLines(t)
		ord.Currency = currency.USD
		ord.Lines[0].Item.Price = num.NewAmount(1000, 2)
		ord.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		_, err := ord.ConvertInto(currency.MXN)
		assert.ErrorContains(t, err, "no exchange rate defined for 'USD' to 'MXN")
	})
}

func baseOrder(t *testing.T, lines ...*bill.Line) *bill.Order {
	t.Helper()
	ord := &bill.Order{
		Series:    "TEST",
		Code:      "00123",
		IssueDate: cal.MakeDate(2022, 6, 13),
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
				Country: "ES",
				Code:    "54387763P",
			},
		},
		Lines: lines,
	}
	return ord
}

func baseOrderWithLines(t *testing.T) *bill.Order {
	ord := baseOrder(t,
		&bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name: "Test Item",
			},
		},
	)
	return ord
}

func TestOrderJSONSchemaExtend(t *testing.T) {
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

	ord := bill.Order{}
	ord.JSONSchemaExtend(js)

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
		it := bill.OrderTypes[0]
		assert.Equal(t, it.Key.String(), prop.OneOf[0].Const)
	})

}

func TestOrderBillable(t *testing.T) {
	ord := baseOrderWithLines(t)
	ord.IssueTime = cal.NewTime(10, 30, 0)
	ord.ValueDate = cal.NewDate(2023, 1, 15)
	ord.Currency = currency.EUR
	ord.ExchangeRates = []*currency.ExchangeRate{{}}
	ord.Complements = []*schema.Object{{}}

	assert.Equal(t, ord.Series, ord.GetSeries())
	assert.Equal(t, ord.Code, ord.GetCode())
	assert.Equal(t, ord.IssueDate, ord.GetIssueDate())
	assert.Equal(t, ord.IssueTime, ord.GetIssueTime())
	assert.Equal(t, ord.ValueDate, ord.GetValueDate())
	assert.Equal(t, ord.Tax, ord.GetTax())
	assert.Equal(t, ord.Preceding, ord.GetPreceding())
	assert.Equal(t, ord.Currency, ord.GetCurrency())
	assert.Equal(t, ord.ExchangeRates, ord.GetExchangeRates())
	assert.Equal(t, ord.Supplier, ord.GetSupplier())
	assert.Equal(t, ord.Customer, ord.GetCustomer())
	assert.Equal(t, ord.Lines, ord.GetLines())
	assert.Equal(t, ord.Discounts, ord.GetDiscounts())
	assert.Equal(t, ord.Charges, ord.GetCharges())
	assert.Equal(t, ord.Payment, ord.GetPaymentDetails())
	assert.Equal(t, ord.Totals, ord.GetTotals())
	assert.Equal(t, ord.Complements, ord.GetComplements())

	ord.SetCode(cbc.Code("002"))
	assert.Equal(t, cbc.Code("002"), ord.Code)

	ord.SetIssueDate(cal.MakeDate(2023, 2, 1))
	assert.Equal(t, cal.MakeDate(2023, 2, 1), ord.IssueDate)

	ord.SetIssueTime(cal.NewTime(11, 30, 0))
	assert.Equal(t, cal.NewTime(11, 30, 0), ord.IssueTime)

	ord.SetCurrency(currency.USD)
	assert.Equal(t, currency.USD, ord.Currency)

	newTotals := &bill.Totals{}
	ord.SetTotals(newTotals)
	assert.Equal(t, newTotals, ord.GetTotals())
}
