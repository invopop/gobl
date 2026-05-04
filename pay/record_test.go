package pay_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordNormalize(t *testing.T) {
	a := &pay.Record{
		Identify:    uuid.Identify{UUID: uuid.Zero},
		Description: "Test advance",
		Percent:     num.NewPercentage(100, 2),
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			"random": "",
		}),
	}
	a.Normalize(nil)
	assert.Empty(t, a.UUID)
	assert.True(t, a.Ext.IsZero())

	a = nil
	assert.NotPanics(t, func() {
		a.Normalize(nil)
	})

}

func TestRecordUnmarshal(t *testing.T) {
	a := new(pay.Record)
	err := json.Unmarshal([]byte(`{"desc":"foo"}`), a)
	require.NoError(t, err)
	assert.Equal(t, "foo", a.Description)
}

func TestRecordCalculateFrom(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a := &pay.Record{
			Percent: num.NewPercentage(10, 2),
		}
		a.CalculateFrom(num.MakeAmount(1000, 2))
		assert.Equal(t, num.NewAmount(100, 2).String(), a.Amount.String())
	})
	t.Run("nil", func(t *testing.T) {
		a := &pay.Record{
			Percent: nil,
		}
		a.CalculateFrom(num.MakeAmount(1000, 2))
		assert.True(t, a.Amount.IsZero())
	})
}

func TestRecordValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a := &pay.Record{
			Identify:    uuid.Identify{UUID: uuid.Zero},
			Description: "Test advance",
			Percent:     num.NewPercentage(100, 2),
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("valid without description", func(t *testing.T) {
		a := &pay.Record{
			Amount: num.MakeAmount(100, 2),
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("valid means key", func(t *testing.T) {
		a := &pay.Record{
			Description: "Test advance",
			Percent:     num.NewPercentage(100, 2),
			Key:         pay.MeansKeyCard,
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("invalid means key", func(t *testing.T) {
		a := &pay.Record{
			Description: "Test advance",
			Percent:     num.NewPercentage(100, 2),
			Key:         "invalid",
		}
		assert.ErrorContains(t, rules.Validate(a), "key must be valid")
	})
}

func TestRecordJSONSchemaExtend(t *testing.T) {
	schema := &jsonschema.Schema{
		Properties: jsonschema.NewProperties(),
	}
	schema.Properties.Set("key", &jsonschema.Schema{
		Type: "string",
	})
	a := &pay.Record{}
	a.JSONSchemaExtend(schema)
	prop, ok := schema.Properties.Get("key")
	require.True(t, ok)
	assert.Len(t, prop.AnyOf, 17)
	assert.Equal(t, cbc.Key("any"), prop.AnyOf[0].Const)
}
