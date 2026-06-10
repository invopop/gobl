package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
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

	t.Run("default type", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Type = ""
		require.NoError(t, dlv.Calculate())
		assert.Equal(t, bill.DeliveryTypeAdvice, dlv.Type)
	})

	t.Run("other type", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Type = bill.DeliveryTypeOther
		require.NoError(t, dlv.Calculate())
		assert.Equal(t, bill.DeliveryTypeOther, dlv.Type)
	})

	t.Run("return tag", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.SetTags(saft.TagReturn)
		dlv.SetAddons(saft.V1)
		require.NoError(t, dlv.Calculate())
		assert.True(t, dlv.HasTags(saft.TagReturn))
	})
}

func TestDeliveryValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		require.NoError(t, dlv.Calculate())
		require.NoError(t, rules.Validate(dlv))
	})

	t.Run("with addons", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.Addons = tax.WithAddons("eu-en16931-v2017") // just for testing
		require.NoError(t, dlv.Calculate())
		require.NoError(t, rules.Validate(dlv))
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

func TestDeliveryTagsValidation(t *testing.T) {
	t.Run("valid tag", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.SetAddons(saft.V1)
		dlv.SetTags(saft.TagReturn)
		assert.NoError(t, dlv.Calculate())
	})

	t.Run("invalid tag", func(t *testing.T) {
		dlv := baseDeliveryWithLines(t)
		dlv.SetTags("invalid-tag")
		err := dlv.Calculate()
		assert.ErrorContains(t, err, "'invalid-tag' undefined")
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
			Addresses: []*org.Address{
				{
					Country: "ES",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
			Addresses: []*org.Address{
				{
					Country: "ES",
				},
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
			"$tags": {
				"items": {
					"$ref": "https://gobl.org/draft-0/cbc/key"
				},
				"type": "array",
				"title": "Tags"
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

	assert.Equal(t, js.Properties.Len(), 5) // from this example

	t.Run("types", func(t *testing.T) {
		prop, ok := js.Properties.Get("type")
		require.True(t, ok)
		assert.Greater(t, len(prop.OneOf), 1)
		it := bill.DeliveryTypes[0]
		assert.Equal(t, it.Key.String(), prop.OneOf[0].Const)
	})

	t.Run("tags", func(t *testing.T) {
		prop, ok := js.Properties.Get("$tags")
		require.True(t, ok)
		require.NotNil(t, prop.Items)
		require.NotEmpty(t, prop.Items.AnyOf)
		// Deliveries have no default tags; only the catch-all "Any" entry is present.
		assert.Equal(t, "Any", prop.Items.AnyOf[0].Title)
	})
}

func TestDeliveryFromToEndpoint(t *testing.T) {
	mkSupplier := func() *org.Party {
		return &org.Party{Endpoints: []*org.Endpoint{{URI: "gobl:supplier.example"}}}
	}
	mkCustomer := func() *org.Party {
		return &org.Party{Endpoints: []*org.Endpoint{{URI: "gobl:customer.example"}}}
	}
	mkDespatcher := func() *org.Party {
		return &org.Party{Endpoints: []*org.Endpoint{{URI: "gobl:despatcher.example"}}}
	}
	mkReceiver := func() *org.Party {
		return &org.Party{Endpoints: []*org.Endpoint{{URI: "gobl:receiver.example"}}}
	}

	despatchTypes := []cbc.Key{
		bill.DeliveryTypeAdvice, bill.DeliveryTypeNote, bill.DeliveryTypeWaybill,
	}
	for _, typ := range despatchTypes {
		t.Run(string(typ)+": supplier→customer (no despatcher/receiver)", func(t *testing.T) {
			dlv := &bill.Delivery{
				Type: typ, Supplier: mkSupplier(), Customer: mkCustomer(),
			}
			assert.Equal(t, "gobl:supplier.example", dlv.FromEndpoint().URI.String())
			assert.Equal(t, "gobl:customer.example", dlv.ToEndpoint().URI.String())
		})
	}

	t.Run("receipt: customer→supplier", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeReceipt, Supplier: mkSupplier(), Customer: mkCustomer(),
		}
		assert.Equal(t, "gobl:customer.example", dlv.FromEndpoint().URI.String())
		assert.Equal(t, "gobl:supplier.example", dlv.ToEndpoint().URI.String())
	})

	t.Run("despatcher and receiver win over supplier/customer", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type:       bill.DeliveryTypeAdvice,
			Supplier:   mkSupplier(),
			Customer:   mkCustomer(),
			Despatcher: mkDespatcher(),
			Receiver:   mkReceiver(),
		}
		assert.Equal(t, "gobl:despatcher.example", dlv.FromEndpoint().URI.String())
		assert.Equal(t, "gobl:receiver.example", dlv.ToEndpoint().URI.String())
	})

	t.Run("receipt prefers receiver→despatcher", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type:       bill.DeliveryTypeReceipt,
			Supplier:   mkSupplier(),
			Customer:   mkCustomer(),
			Despatcher: mkDespatcher(),
			Receiver:   mkReceiver(),
		}
		assert.Equal(t, "gobl:receiver.example", dlv.FromEndpoint().URI.String())
		assert.Equal(t, "gobl:despatcher.example", dlv.ToEndpoint().URI.String())
	})

	t.Run("nil delivery is a no-op", func(t *testing.T) {
		var dlv *bill.Delivery
		assert.Nil(t, dlv.FromEndpoint())
		assert.Nil(t, dlv.ToEndpoint())
	})
}
