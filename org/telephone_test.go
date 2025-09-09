package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestTelephoneNormalizeAndValidate(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tel *org.Telephone
		assert.NotPanics(t, func() {
			tel.Normalize()
		})
	})
	t.Run("basic normalization", func(t *testing.T) {
		tel := &org.Telephone{
			Number: "  +123 456 7890  ",
		}
		tel.Normalize()
		assert.Equal(t, "+123 456 7890", tel.Number)
		assert.NoError(t, tel.Validate())
	})

	t.Run("empty number", func(t *testing.T) {
		tel := &org.Telephone{
			Number: "   ",
		}
		tel.Normalize()
		assert.Equal(t, "", tel.Number)
		assert.ErrorContains(t, tel.Validate(), "num: cannot be blank")
	})

	t.Run("allow complex numbers", func(t *testing.T) {
		tel := &org.Telephone{
			Number: "+1 (123) 456-7890 ext. 123",
		}
		tel.Normalize()
		assert.Equal(t, "+1 (123) 456-7890 ext. 123", tel.Number)
		assert.NoError(t, tel.Validate())
	})
}
