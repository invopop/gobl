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

func TestLineDiscountNormalize(t *testing.T) {
	l := &bill.LineDiscount{
		Code:    " FOO--BAR ",
		Percent: num.NewPercentage(200, 3),
		Ext:     tax.Extensions{},
	}
	l.Normalize(nil)
	assert.Equal(t, "20.0%", l.Percent.String())
	assert.Equal(t, "FOO-BAR", l.Code.String())
	assert.Nil(t, l.Ext)
}

func TestLineDiscountValidation(t *testing.T) {
	l := &bill.LineDiscount{
		Key:    "foo",
		Code:   "BAR",
		Amount: num.MakeAmount(100, 2),
	}
	err := l.Validate()
	assert.NoError(t, err)

	l.Amount = num.MakeAmount(0, 2)
	err = l.Validate()
	assert.ErrorContains(t, err, "amount: must not be zero")
}

func TestLineDiscountJSONSchema(t *testing.T) {
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
	schema := bill.LineDiscount{}
	schema.JSONSchemaExtend(js)

	props, ok := js.Properties.Get("key")
	assert.True(t, ok)
	assert.NotNil(t, props)
	assert.True(t, len(props.AnyOf) > 10)
}

func TestCleanLineDiscounts(t *testing.T) {
	discounts := []*bill.LineDiscount{
		{Amount: num.MakeAmount(100, 2)},
		{Amount: num.MakeAmount(0, 2)},
		{Reason: "test", Amount: num.MakeAmount(0, 2)},
		{Code: "ABC", Percent: num.NewPercentage(0, 2)},
	}
	cleaned := bill.CleanLineDiscounts(discounts)
	assert.Len(t, cleaned, 3)
}
