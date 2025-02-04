package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrgNoteNormalize(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("accepts nil", func(t *testing.T) {
		var n *org.Note
		assert.NotPanics(t, func() {
			ad.Normalizer(n)
		})
	})

	t.Run("accepts empty", func(t *testing.T) {
		n := &org.Note{}
		assert.NotPanics(t, func() {
			ad.Normalizer(n)
		})
	})

	t.Run("converts valid", func(t *testing.T) {
		n := &org.Note{
			Key:  org.NoteKeyGeneral,
			Text: "This is a general note test",
		}
		assert.NotPanics(t, func() {
			ad.Normalizer(n)
		})
		assert.Equal(t, "AAI", n.Ext[untdid.ExtKeyTextSubject].String())
	})
}

func TestOrgItemNormalize(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("accepts nil", func(t *testing.T) {
		var item *org.Item
		assert.NotPanics(t, func() {
			ad.Normalizer(item)
		})
	})

	t.Run("sets default", func(t *testing.T) {
		item := &org.Item{}
		ad.Normalizer(item)
		assert.Equal(t, org.UnitOne, item.Unit)
	})

	t.Run("maintains valid", func(t *testing.T) {
		item := &org.Item{
			Unit: org.UnitHour,
		}
		ad.Normalizer(item)
		assert.Equal(t, org.UnitHour, item.Unit)
	})
}

func TestOrgItemValidate(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("missing unit", func(t *testing.T) {
		item := &org.Item{}
		assert.ErrorContains(t, ad.Validator(item), "unit: cannot be blank (BR-23).")
	})

	t.Run("validates unit", func(t *testing.T) {
		item := &org.Item{
			Unit: org.UnitOne,
		}
		assert.NoError(t, ad.Validator(item))
	})
}
