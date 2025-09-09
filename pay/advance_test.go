package pay_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvanceNormalize(t *testing.T) {
	a := &pay.Advance{
		Identify:    uuid.Identify{UUID: uuid.Zero},
		Description: "Test advance",
		Percent:     num.NewPercentage(100, 2),
		Ext: tax.Extensions{
			"random": "",
		},
	}
	a.Normalize()
	assert.Empty(t, a.UUID)
	assert.Empty(t, a.Ext)

	a = nil
	assert.NotPanics(t, func() {
		a.Normalize()
	})

}

func TestAdvanceUnmarshal(t *testing.T) {
	a := new(pay.Advance)
	err := json.Unmarshal([]byte(`{"desc":"foo"}`), a)
	require.NoError(t, err)
	assert.Equal(t, "foo", a.Description)
}

func TestAdvanceCalculateFrom(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a := &pay.Advance{
			Percent: num.NewPercentage(10, 2),
		}
		a.CalculateFrom(num.MakeAmount(1000, 2))
		assert.Equal(t, num.NewAmount(100, 2).String(), a.Amount.String())
	})
	t.Run("nil", func(t *testing.T) {
		a := &pay.Advance{
			Percent: nil,
		}
		a.CalculateFrom(num.MakeAmount(1000, 2))
		assert.True(t, a.Amount.IsZero())
	})
}

func TestAdvanceValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a := &pay.Advance{
			Identify:    uuid.Identify{UUID: uuid.Zero},
			Description: "Test advance",
			Percent:     num.NewPercentage(100, 2),
		}
		assert.NoError(t, a.Validate())
	})
	t.Run("invalid", func(t *testing.T) {
		a := &pay.Advance{
			Amount: num.MakeAmount(100, 2),
		}
		assert.ErrorContains(t, a.Validate(), "description: cannot be blank")
	})
	t.Run("valid means key", func(t *testing.T) {
		a := &pay.Advance{
			Description: "Test advance",
			Percent:     num.NewPercentage(100, 2),
			Key:         pay.MeansKeyCard,
		}
		assert.NoError(t, a.Validate())
	})
	t.Run("invalid means key", func(t *testing.T) {
		a := &pay.Advance{
			Description: "Test advance",
			Percent:     num.NewPercentage(100, 2),
			Key:         "invalid",
		}
		assert.ErrorContains(t, a.Validate(), "key: must be or start with a valid key.")
	})
}

func TestAdvanceJSONSchemaExtend(t *testing.T) {
	schema := &jsonschema.Schema{
		Properties: jsonschema.NewProperties(),
	}
	schema.Properties.Set("key", &jsonschema.Schema{
		Type: "string",
	})
	a := &pay.Advance{}
	a.JSONSchemaExtend(schema)
	prop, ok := schema.Properties.Get("key")
	require.True(t, ok)
	assert.Len(t, prop.AnyOf, 15)
	assert.Equal(t, cbc.Key("any"), prop.AnyOf[0].Const)
}
