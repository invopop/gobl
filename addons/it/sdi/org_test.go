package sdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateOrgAttribute(t *testing.T) {
	t.Run("type within 10 characters is valid", func(t *testing.T) {
		attr := &org.Attribute{
			Type: cbc.Code("TARGA"),
			Text: "AB123CD",
		}
		err := rules.Validate(attr, tax.AddonContext(sdi.V1))
		require.NoError(t, err)
	})

	t.Run("type longer than 10 characters is rejected", func(t *testing.T) {
		attr := &org.Attribute{
			Type: cbc.Code("TOOLONGVALUE"), // 12 characters
			Text: "x",
		}
		err := rules.Validate(attr, tax.AddonContext(sdi.V1))
		assert.ErrorContains(t, err, "type cannot be longer than 10 characters")
	})

	t.Run("key-based attribute is not constrained by TipoDato length", func(t *testing.T) {
		attr := &org.Attribute{
			Key:  org.AttributeKeyColor,
			Text: "red",
		}
		err := rules.Validate(attr, tax.AddonContext(sdi.V1))
		require.NoError(t, err)
	})
}
