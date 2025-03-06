package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineChargeNormalize(t *testing.T) {
	l := &bill.LineCharge{
		Code:    " FOO--BAR ",
		Percent: num.NewPercentage(200, 3),
		Ext:     tax.Extensions{},
	}
	l.Normalize(nil)
	assert.Equal(t, "20.0%", l.Percent.String())
	assert.Equal(t, "FOO-BAR", l.Code.String())
	assert.Nil(t, l.Ext)
}

func TestLineChargeValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l := &bill.LineCharge{
			Key:    "foo",
			Code:   "BAR",
			Amount: num.MakeAmount(100, 2),
		}
		err := l.Validate()
		assert.NoError(t, err)

		l.Amount = num.MakeAmount(0, 2)
		err = l.Validate()
		assert.ErrorContains(t, err, "amount: must not be zero")
	})
	t.Run("with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code: "IEPS",
			Base: num.NewAmount(3000, 2),
		}
		err := l.Validate()
		assert.ErrorContains(t, err, "percent: cannot be blank")
	})
	t.Run("valid with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:    "IEPS",
			Base:    num.NewAmount(3000, 2),
			Percent: num.NewPercentage(4, 3),
			Amount:  num.MakeAmount(120, 2),
		}
		assert.NoError(t, l.Validate())
	})
	t.Run("valid with rate and quantity", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:     "IEPS",
			Quantity: num.NewAmount(100, 0), // e.g. grams
			Rate:     num.NewAmount(2, 0),   // 1 per gram
			Amount:   num.MakeAmount(200, 2),
		}
		assert.NoError(t, l.Validate())
	})
	t.Run("missing rate with quantity", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:     "IEPS",
			Quantity: num.NewAmount(100, 0), // e.g. grams
			Amount:   num.MakeAmount(200, 2),
		}
		assert.ErrorContains(t, l.Validate(), "rate: cannot be blank with quantity")
	})
	t.Run("quantity with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:     "IEPS",
			Base:     num.NewAmount(3000, 2),
			Percent:  num.NewPercentage(4, 3),
			Quantity: num.NewAmount(100, 0), // e.g. grams
			Amount:   num.MakeAmount(200, 2),
		}
		assert.ErrorContains(t, l.Validate(), "quantity: must be blank with base or percent")
	})
	t.Run("rate with base", func(t *testing.T) {
		l := &bill.LineCharge{
			Code:    "IEPS",
			Base:    num.NewAmount(3000, 2),
			Percent: num.NewPercentage(4, 3),
			Rate:    num.NewAmount(1, 0),
			Amount:  num.MakeAmount(200, 2),
		}
		assert.ErrorContains(t, l.Validate(), "rate: must be blank with base or percent")
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
