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

func TestDeliveryCalculate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		require.NoError(t, dlv.Calculate())
		assert.Nil(t, dlv.Totals)
	})
}

func TestDeliveryValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		require.NoError(t, dlv.Calculate())
		require.NoError(t, dlv.Validate())
	})

	t.Run("with addons", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Addons = tax.WithAddons("eu-en16931-v2017") // just for testing
		require.NoError(t, dlv.Calculate())
		require.NoError(t, dlv.Validate())
	})
}

func TestDeliveryConvertInto(t *testing.T) {
	t.Run("no prices", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Currency = currency.USD
		dlv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		d2, err := dlv.ConvertInto(currency.EUR)
		assert.NoError(t, err)
		assert.Nil(t, d2.Totals)
	})
	t.Run("basic", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Currency = currency.USD
		dlv.Lines[0].Item.Price = num.NewAmount(1000, 2)
		dlv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		d2, err := dlv.ConvertInto(currency.EUR)
		assert.NoError(t, err)
		assert.Equal(t, currency.EUR, d2.Currency)
		assert.Equal(t, "87.5970", d2.Lines[0].Total.String())
	})
	t.Run("same currency", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Currency = currency.USD
		dlv.Lines[0].Item.Price = num.NewAmount(1000, 2)
		dlv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		d2, err := dlv.ConvertInto(currency.USD)
		assert.NoError(t, err)
		assert.Equal(t, currency.USD, d2.Currency)
		assert.Equal(t, "100.00", d2.Lines[0].Total.String())
	})
	t.Run("missing rate", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Currency = currency.USD
		dlv.Lines[0].Item.Price = num.NewAmount(1000, 2)
		dlv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(875967, 6),
			},
		}
		_, err := dlv.ConvertInto(currency.MXN)
		assert.ErrorContains(t, err, "no exchange rate defined for 'USD' to 'MXN")
	})
}

func baseDelivery(t *testing.T, lines ...*bill.Line) *bill.Delivery {
	t.Helper()
	dlv := &bill.Delivery{
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
	return dlv
}

func baseDeliveryWithLines(t *testing.T) *bill.Delivery {
	dlv := baseDelivery(t,
		&bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name: "Test Item",
			},
		},
	)
	return dlv
}

func TestDeliveryJSONSchemaExtend(t *testing.T) {
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

	dlv := bill.Delivery{}
	dlv.JSONSchemaExtend(js)

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
		it := bill.DeliveryTypes[0]
		assert.Equal(t, it.Key.String(), prop.OneOf[0].Const)
	})

}

func TestDeliveryBillable(t *testing.T) {
	dlv := baseDeliveryWithLines(t)
	dlv.IssueTime = cal.NewTime(10, 30, 0)
	dlv.ValueDate = cal.NewDate(2023, 1, 15)
	dlv.Currency = currency.EUR
	dlv.ExchangeRates = []*currency.ExchangeRate{{}}
	dlv.Complements = []*schema.Object{{}}

	assert.Equal(t, dlv.Series, dlv.GetSeries())
	assert.Equal(t, dlv.Code, dlv.GetCode())
	assert.Equal(t, dlv.IssueDate, dlv.GetIssueDate())
	assert.Equal(t, dlv.IssueTime, dlv.GetIssueTime())
	assert.Equal(t, dlv.ValueDate, dlv.GetValueDate())
	assert.Equal(t, dlv.Tax, dlv.GetTax())
	assert.Equal(t, dlv.Preceding, dlv.GetPreceding())
	assert.Equal(t, dlv.Currency, dlv.GetCurrency())
	assert.Equal(t, dlv.ExchangeRates, dlv.GetExchangeRates())
	assert.Equal(t, dlv.Supplier, dlv.GetSupplier())
	assert.Equal(t, dlv.Customer, dlv.GetCustomer())
	assert.Equal(t, dlv.Lines, dlv.GetLines())
	assert.Equal(t, dlv.Discounts, dlv.GetDiscounts())
	assert.Equal(t, dlv.Charges, dlv.GetCharges())
	assert.Equal(t, dlv.Totals, dlv.GetTotals())
	assert.Equal(t, dlv.Complements, dlv.GetComplements())
	assert.Nil(t, dlv.GetPaymentDetails())

	dlv.SetCode(cbc.Code("002"))
	assert.Equal(t, cbc.Code("002"), dlv.Code)

	dlv.SetIssueDate(cal.MakeDate(2023, 2, 1))
	assert.Equal(t, cal.MakeDate(2023, 2, 1), dlv.IssueDate)

	dlv.SetIssueTime(cal.NewTime(11, 30, 0))
	assert.Equal(t, cal.NewTime(11, 30, 0), dlv.IssueTime)

	dlv.SetCurrency(currency.USD)
	assert.Equal(t, currency.USD, dlv.Currency)

	newTotals := &bill.Totals{}
	dlv.SetTotals(newTotals)
	assert.Equal(t, newTotals, dlv.GetTotals())
}
