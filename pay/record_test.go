package pay_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
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
		Ref:         " TRX 2024/0012\t",
		Description: "Test advance",
		Percent:     num.NewPercentage(100, 2),
		DirectDebit: &pay.DirectDebit{
			Ref: " MANDATE-001 ",
		},
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			"random": "",
		}),
	}
	norm.Normalize(a)
	assert.Empty(t, a.UUID)
	assert.Equal(t, "TRX 2024/0012", a.Ref.String())
	assert.Equal(t, "MANDATE-001", a.DirectDebit.Ref.String())
	assert.True(t, a.Ext.IsZero())

	a = nil
	assert.NotPanics(t, func() {
		norm.Normalize(a)
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

func TestRecordCalculateFromTaxes(t *testing.T) {
	zero := num.MakeAmount(0, 2)
	tt := &tax.Total{
		Categories: []*tax.CategoryTotal{
			{
				Code: tax.CategoryVAT,
				Rates: []*tax.RateTotal{
					{Amount: num.MakeAmount(2100, 2)},
				},
			},
		},
	}
	t.Run("with taxes", func(t *testing.T) {
		a := &pay.Record{
			Taxes: []*tax.Filter{{Category: tax.CategoryVAT}},
		}
		a.CalculateFromTaxes(zero, tt)
		assert.Equal(t, "21.00", a.Amount.String())
	})
	t.Run("without taxes", func(t *testing.T) {
		a := &pay.Record{
			Amount: num.MakeAmount(500, 2),
		}
		a.CalculateFromTaxes(zero, tt)
		assert.Equal(t, "5.00", a.Amount.String())
	})
	t.Run("nil", func(t *testing.T) {
		var a *pay.Record
		assert.NotPanics(t, func() {
			a.CalculateFromTaxes(zero, tt)
		})
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
	t.Run("valid with taxes", func(t *testing.T) {
		a := &pay.Record{
			Key:   pay.MeansKeyWaiver,
			Taxes: []*tax.Filter{{Category: tax.CategoryVAT}},
		}
		assert.NoError(t, rules.Validate(a))
	})
	t.Run("taxes with percent", func(t *testing.T) {
		a := &pay.Record{
			Percent: num.NewPercentage(100, 2),
			Taxes:   []*tax.Filter{{Category: tax.CategoryVAT}},
		}
		assert.ErrorContains(t, rules.Validate(a), "taxes must be blank with percent")
	})
	t.Run("invalid tax filter", func(t *testing.T) {
		a := &pay.Record{
			Taxes: []*tax.Filter{{Key: tax.KeyStandard}},
		}
		assert.ErrorContains(t, rules.Validate(a), "tax filter category is required")
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
	assert.Len(t, prop.AnyOf, 18)
	assert.Equal(t, cbc.Key("any"), prop.AnyOf[0].Const)
}
