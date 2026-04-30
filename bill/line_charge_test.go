package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineChargeNormalize(t *testing.T) {
	l := &bill.LineCharge{
		Code:    " FOO--BAR ",
		Percent: num.NewPercentage(200, 3),
		Ext:     tax.ExtensionsOf(cbc.CodeMap{}),
	}
	l.Normalize(nil)
	assert.Equal(t, "20.0%", l.Percent.String())
	assert.Equal(t, "FOO-BAR", l.Code.String())
	assert.True(t, l.Ext.IsZero())
}

func TestLineChargeValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l := &bill.LineCharge{
			Key:    "foo",
			Code:   "BAR",
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(l)
		assert.NoError(t, err)

		l.Amount = num.MakeAmount(0, 2)
		err = rules.Validate(l)
		assert.NoError(t, err)
	})
	t.Run("with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code: "IEPS",
			Base: num.NewAmount(3000, 2),
		}
		err := rules.Validate(l)
		assert.ErrorContains(t, err, "percent is required when base is set")
	})
	t.Run("valid with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:    "IEPS",
			Base:    num.NewAmount(3000, 2),
			Percent: num.NewPercentage(4, 3),
			Amount:  num.MakeAmount(120, 2),
		}
		assert.NoError(t, rules.Validate(l))
	})
	t.Run("valid with rate and quantity", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:     "IEPS",
			Quantity: num.NewAmount(100, 0), // e.g. grams
			Unit:     "g",
			Rate:     num.NewAmount(2, 0), // 1 per gram
			Amount:   num.MakeAmount(200, 2),
		}
		assert.NoError(t, rules.Validate(l))
	})
	t.Run("missing rate with quantity", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:     "IEPS",
			Quantity: num.NewAmount(100, 0), // e.g. grams
			Amount:   num.MakeAmount(200, 2),
		}
		assert.ErrorContains(t, rules.Validate(l), "rate is required when quantity is set")
	})
	t.Run("missing rate with quantity", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:   "IEPS",
			Unit:   "l",
			Amount: num.MakeAmount(200, 2),
		}
		assert.ErrorContains(t, rules.Validate(l), "unit must be blank without quantity")
	})

	t.Run("quantity with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:     "IEPS",
			Base:     num.NewAmount(3000, 2),
			Percent:  num.NewPercentage(4, 3),
			Quantity: num.NewAmount(100, 0), // e.g. grams
			Amount:   num.MakeAmount(200, 2),
		}
		assert.ErrorContains(t, rules.Validate(l), "quantity must be blank with base or percent")
	})
	t.Run("rate with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:    "IEPS",
			Base:    num.NewAmount(3000, 2),
			Percent: num.NewPercentage(4, 3),
			Rate:    num.NewAmount(1, 0),
			Amount:  num.MakeAmount(200, 2),
		}
		assert.ErrorContains(t, rules.Validate(l), "rate must be blank with base or percent")
	})
}

func TestLineChargeJSONSchema(t *testing.T) {
	eg := `{
		"type": "object",
		"properties": {
			"key": {
				"$ref": "https://gobl.org/draft-0/cbc/key"
			}
		}
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))
	schema := bill.LineCharge{}
	schema.JSONSchemaExtend(js)

	props, ok := js.Properties.Get("key")
	assert.True(t, ok)
	assert.NotNil(t, props)
	assert.Equal(t, 12, len(props.AnyOf))
}

func TestCleanLineCharges(t *testing.T) {
	lines := []*bill.LineCharge{
		{Amount: num.MakeAmount(100, 2)},
		{Amount: num.MakeAmount(0, 2)},
		{Reason: "test", Amount: num.MakeAmount(0, 2)},
		{Code: "ABC", Percent: num.NewPercentage(0, 2)},
	}
	cleaned := bill.CleanLineCharges(lines)
	assert.Len(t, cleaned, 3)
}
