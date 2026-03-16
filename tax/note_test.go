package tax_test

import (
	"context"
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNoteValidation(t *testing.T) {
	ctx := tax.RegimeDefFor("DE").WithContext(context.Background())

	t.Run("valid note", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "exempt",
			Text:     "Exempt under Article 132",
		}
		assert.NoError(t, n.ValidateWithContext(ctx))
	})

	t.Run("with extensions", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "reverse-charge",
			Text:     "Reverse charge applies",
			Ext: tax.Extensions{
				"untdid-tax-category": "AE",
			},
		}
		assert.NoError(t, n.ValidateWithContext(ctx))
	})

	t.Run("text only", func(t *testing.T) {
		n := &tax.Note{
			Text: "Some exemption reason",
		}
		assert.NoError(t, n.ValidateWithContext(ctx))
	})

	t.Run("missing text", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "exempt",
		}
		err := n.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "text: cannot be blank")
	})

	t.Run("empty note", func(t *testing.T) {
		n := &tax.Note{}
		err := n.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "text: cannot be blank")
	})

	t.Run("invalid category", func(t *testing.T) {
		n := &tax.Note{
			Category: "INVALID",
			Text:     "Some reason",
		}
		err := n.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "cat")
	})

	t.Run("invalid key for category", func(t *testing.T) {
		n := &tax.Note{
			Category: "VAT",
			Key:      "not-a-real-key",
			Text:     "Some reason",
		}
		err := n.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "key")
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
			Ext: tax.Extensions{
				"untdid-tax-category": "E",
				"empty-key":           "",
			},
		}
		n.Normalize(nil)
		assert.Equal(t, "E", n.Ext.Get("untdid-tax-category").String())
		assert.False(t, n.Ext.Has("empty-key"))
	})
}
