package tax_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNoteValidation(t *testing.T) {
	t.Run("valid note", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "exempt",
			Text:     "Exempt under Article 132",
		}
		assert.NoError(t, rules.Validate(n))
	})

	t.Run("with extensions", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "reverse-charge",
			Text:     "Reverse charge applies",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"untdid-tax-category": "AE",
			}),
		}
		assert.NoError(t, rules.Validate(n))
	})

	t.Run("text only", func(t *testing.T) {
		n := &tax.Note{
			Text: "Some exemption reason",
		}
		assert.NoError(t, rules.Validate(n))
	})

	t.Run("missing text", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "exempt",
		}
		err := rules.Validate(n)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "text")
	})

	t.Run("empty note", func(t *testing.T) {
		n := &tax.Note{}
		err := rules.Validate(n)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "text")
	})

	t.Run("free-form key for category", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "reverse-charge",
			Text:     "Some reason",
		}
		assert.NoError(t, rules.Validate(n))
	})
}

func TestNoteNormalize(t *testing.T) {
	t.Run("nil note", func(t *testing.T) {
		var n *tax.Note
		assert.NotPanics(t, func() {
			n.Normalize(nil)
		})
	})

	t.Run("cleans extensions", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "exempt",
			Text:     "Exempt",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"untdid-tax-category": "E",
				"empty-key":           "",
			}),
		}
		n.Normalize(nil)
		assert.Equal(t, "E", n.Ext.Get("untdid-tax-category").String())
		assert.False(t, n.Ext.Has("empty-key"))
	})
}
