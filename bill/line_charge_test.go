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
