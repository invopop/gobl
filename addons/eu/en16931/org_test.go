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
