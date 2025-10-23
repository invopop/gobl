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

func TestOrgInboxNormalize(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)

	t.Run("accepts nil", func(t *testing.T) {
		var i *org.Inbox
		assert.NotPanics(t, func() {
			ad.Normalizer(i)
		})
	})
	t.Run("accepts empty", func(t *testing.T) {
		i := &org.Inbox{}
		assert.NotPanics(t, func() {
			ad.Normalizer(i)
		})
	})
	t.Run("normalizes scheme and code", func(t *testing.T) {
		i := &org.Inbox{
			Code: "0004:BAR",
		}
		ad.Normalizer(i)
		assert.Equal(t, "0004", i.Scheme.String())
		assert.Equal(t, "BAR", i.Code.String())
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

func TestOrgAttachmentValidation(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("blank", func(t *testing.T) {
		a := &org.Attachment{}
		assert.ErrorContains(t, ad.Validator(a), "code: cannot be blank.")
	})
	t.Run("with code", func(t *testing.T) {
		a := &org.Attachment{
			Code: "123",
		}
		assert.NoError(t, ad.Validator(a))
	})
}

func TestOrgPartyValidate(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("no inboxes", func(t *testing.T) {
		p := &org.Party{
			Addresses: []*org.Address{
				{
					Country: "FR",
				},
			},
		}
		assert.NoError(t, ad.Validator(p))
	})

	t.Run("one inbox", func(t *testing.T) {
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: "scheme1",
					Code:   "code1",
				},
			},
			Addresses: []*org.Address{
				{
					Country: "FR",
				},
			},
		}
		assert.NoError(t, ad.Validator(p))
	})

	t.Run("multiple inboxes", func(t *testing.T) {
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: "scheme1",
					Code:   "code1",
				},
				{
					Scheme: "scheme2",
					Code:   "code2",
				},
			},
			Addresses: []*org.Address{
				{
					Country: "FR",
				},
			},
		}
		assert.ErrorContains(t, ad.Validator(p), "inboxes: cannot have more than one inbox (BT-34, BT-49).")
	})
	t.Run("missing addresses", func(t *testing.T) {
		p := &org.Party{}
		assert.ErrorContains(t, ad.Validator(p), "addresses: cannot be blank.")
	})
}

func TestOrgInboxValidate(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("missing scheme and code", func(t *testing.T) {
		i := &org.Inbox{}
		// Not specific for addon, but this is important to check
		assert.ErrorContains(t, i.Validate(), "code: cannot be blank without url or email")
	})

	t.Run("missing scheme", func(t *testing.T) {
		i := &org.Inbox{
			Code: "code1",
		}
		assert.ErrorContains(t, ad.Validator(i), "scheme: cannot be blank with code (BR-62, BR-63)")
	})

	t.Run("missing code", func(t *testing.T) {
		i := &org.Inbox{
			Scheme: "scheme1",
		}
		assert.ErrorContains(t, ad.Validator(i), "code: cannot be blank")
	})

	t.Run("valid inbox", func(t *testing.T) {
		i := &org.Inbox{
			Scheme: "scheme1",
			Code:   "code1",
		}
		assert.NoError(t, ad.Validator(i))
	})
}
