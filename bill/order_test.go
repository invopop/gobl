package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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

	t.Run("with nil array entries", func(t *testing.T) {
		ord := baseOrderWithLines(t)
		ord.ExchangeRates = append(ord.ExchangeRates, nil)
		ord.Preceding = append(ord.Preceding, nil)
		ord.Lines = append(ord.Lines, nil)
		ord.Discounts = append(ord.Discounts, nil)
		ord.Charges = append(ord.Charges, nil)
		ord.Notes = append(ord.Notes, nil)
		ord.Complements = append(ord.Complements, nil)
		ord.Attachments = append(ord.Attachments, nil)
		require.NoError(t, ord.Calculate())
		err := ord.Validate()
		assert.ErrorContains(t, err, "exchange_rates: (0: is required.)")
		assert.ErrorContains(t, err, "preceding: (0: is required.)")
		assert.ErrorContains(t, err, "lines: (1: is required.)")
		assert.ErrorContains(t, err, "discounts: (0: is required.)")
		assert.ErrorContains(t, err, "charges: (0: is required.)")
		assert.ErrorContains(t, err, "notes: (0: is required.)")
		assert.ErrorContains(t, err, "complements: (0: is required.)")
		assert.ErrorContains(t, err, "attachments: (0: is required.)")
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
