package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComboNormalize(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "general",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Rate, tax.RateGeneral)
	})
	t.Run("migrate standard rate to general", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "standard",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Rate, tax.RateGeneral)
	})
	t.Run("migrate standard with extension rate to general", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "standard+eqs",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Rate.String(), "general+eqs")
	})
	t.Run("migrate zero rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "zero",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Key, tax.KeyZero)
		assert.Equal(t, c.Percent.String(), "0%")
		assert.Empty(t, c.Rate)
	})
	t.Run("assign zero percent", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Key:      tax.KeyZero,
		}
		c.Normalize(nil)
		assert.Equal(t, "0%", c.Percent.String())
	})
	t.Run("remove exempt rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "exempt",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Empty(t, c.Key)
		assert.Empty(t, c.Rate)
	})
	t.Run("remove exempt rate and key", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Key:      tax.KeyExempt,
			Rate:     "exempt",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Empty(t, c.Key)
		assert.Empty(t, c.Rate)
	})
	t.Run("migrate exempt+reverse-charge rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "exempt+reverse-charge",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Key, tax.KeyReverseCharge)
		assert.Empty(t, c.Rate)
	})
	t.Run("migrate exempt+export rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "exempt+export",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Key, tax.KeyExport)
		assert.Empty(t, c.Rate)
	})
	t.Run("migrate exempt+export+eea rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Rate:     "exempt+export+eea",
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Key, tax.KeyIntraCommunity)
		assert.Empty(t, c.Rate)
	})
	t.Run("migrate zero percentage", func(t *testing.T) {
		c := &tax.Combo{
			Category: "VAT",
			Percent:  num.NewPercentage(0, 2),
		}
		c.Normalize(nil)
		assert.Equal(t, c.Category, tax.CategoryVAT)
		assert.Equal(t, c.Key, tax.KeyZero)
		assert.Empty(t, c.Rate)
		assert.Equal(t, c.Percent.String(), "0%")
	})

}

func TestComboUnmarshal(t *testing.T) {
	t.Run("with tags", func(t *testing.T) {
		data := []byte(`{"cat":"VAT","tags":["standard"],"percent":"20%"}`)
		var c tax.Combo
		err := json.Unmarshal(data, &c)
		require.NoError(t, err)
		assert.Equal(t, c.Category, cbc.Code("VAT"))
		assert.Equal(t, c.Rate, cbc.Key("standard"))
	})
	t.Run("with tags and rate", func(t *testing.T) {
		data := []byte(`{"cat":"VAT","rate":"general","tags":["standard"],"percent":"20%"}`)
		var c tax.Combo
		err := json.Unmarshal(data, &c)
		require.NoError(t, err)
		assert.Equal(t, c.Category, cbc.Code("VAT"))
		assert.Equal(t, c.Rate, tax.RateGeneral) // don't override
	})
}

func TestComboJSONSchemaExtend(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		c := tax.Combo{}
		s := new(jsonschema.Schema)
		c.JSONSchemaExtend(s)

		assert.NotEmpty(t, s.AnyOf)
		require.Len(t, s.AnyOf, 1)
		p, ok := s.AnyOf[0].If.Properties.Get("cat")
		require.True(t, ok)
		assert.Equal(t, p.Const, "VAT")
	})
}
