package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrgNoteNormalize(t *testing.T) {
	t.Run("accepts nil", func(t *testing.T) {
		var n *org.Note
		assert.NotPanics(t, func() {
			norm.Normalize(n, tax.AddonContext(en16931.V2017))
		})
	})

	t.Run("accepts empty", func(t *testing.T) {
		n := &org.Note{}
		assert.NotPanics(t, func() {
			norm.Normalize(n, tax.AddonContext(en16931.V2017))
		})
	})

	t.Run("converts valid", func(t *testing.T) {
		n := &org.Note{
			Key:  org.NoteKeyGeneral,
			Text: "This is a general note test",
		}
		assert.NotPanics(t, func() {
			norm.Normalize(n, tax.AddonContext(en16931.V2017))
		})
		assert.Equal(t, "AAI", n.Ext.Get(untdid.ExtKeyTextSubject).String())
	})
}

func TestOrgItemNormalize(t *testing.T) {
	t.Run("accepts nil", func(t *testing.T) {
		var item *org.Item
		assert.NotPanics(t, func() {
			norm.Normalize(item, tax.AddonContext(en16931.V2017))
		})
	})

	t.Run("sets default", func(t *testing.T) {
		item := &org.Item{}
		norm.Normalize(item, tax.AddonContext(en16931.V2017))
		assert.Equal(t, org.UnitOne, item.Unit)
	})

	t.Run("maintains valid", func(t *testing.T) {
		item := &org.Item{
			Unit: org.UnitHour,
		}
		norm.Normalize(item, tax.AddonContext(en16931.V2017))
		assert.Equal(t, org.UnitHour, item.Unit)
	})
}

func TestOrgIdentityNormalize(t *testing.T) {
	t.Run("accepts nil", func(t *testing.T) {
		var id *org.Identity
		assert.NotPanics(t, func() {
			norm.Normalize(id, tax.AddonContext(en16931.V2017))
		})
	})
	t.Run("accepts empty", func(t *testing.T) {
		id := &org.Identity{}
		assert.NotPanics(t, func() {
			norm.Normalize(id, tax.AddonContext(en16931.V2017))
		})
	})
	t.Run("normalizes key", func(t *testing.T) {
		id := &org.Identity{
			Key:  "gln",
			Code: "1234567890123",
		}
		norm.Normalize(id, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "gln", id.Key.String())
		assert.Equal(t, "1234567890123", id.Code.String())
		assert.Equal(t, "0088", id.Ext.Get(iso.ExtKeySchemeID).String())
	})
	t.Run("overrides key", func(t *testing.T) {
		id := &org.Identity{
			Key:  "gln",
			Code: "1234567890123",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "9999",
			}),
		}
		norm.Normalize(id, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "gln", id.Key.String())
		assert.Equal(t, "1234567890123", id.Code.String())
		assert.Equal(t, "0088", id.Ext.Get(iso.ExtKeySchemeID).String())
	})
	t.Run("normalizes type", func(t *testing.T) {
		id := &org.Identity{
			Type: "SIREN",
			Code: "1234567890123",
		}
		norm.Normalize(id, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "SIREN", id.Type.String())
		assert.Equal(t, "1234567890123", id.Code.String())
		assert.Equal(t, "0002", id.Ext.Get(iso.ExtKeySchemeID).String())
	})
	t.Run("overrides key", func(t *testing.T) {
		id := &org.Identity{
			Type: "SIREN",
			Code: "1234567890123",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "9999",
			}),
		}
		norm.Normalize(id, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "SIREN", id.Type.String())
		assert.Equal(t, "1234567890123", id.Code.String())
		assert.Equal(t, "0002", id.Ext.Get(iso.ExtKeySchemeID).String())
	})
}

func TestOrgInboxNormalize(t *testing.T) {

	t.Run("accepts nil", func(t *testing.T) {
		var i *org.Inbox
		assert.NotPanics(t, func() {
			norm.Normalize(i, tax.AddonContext(en16931.V2017))
		})
	})
	t.Run("accepts empty", func(t *testing.T) {
		i := &org.Inbox{}
		assert.NotPanics(t, func() {
			norm.Normalize(i, tax.AddonContext(en16931.V2017))
		})
	})
	t.Run("normalizes scheme and code", func(t *testing.T) {
		i := &org.Inbox{
			Code: "0004:BAR",
		}
		norm.Normalize(i, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "0004", i.Scheme.String())
		assert.Equal(t, "BAR", i.Code.String())
	})
}

func TestOrgPartyNormalizePeppolEndpoint(t *testing.T) {
	t.Run("accepts nil", func(t *testing.T) {
		var p *org.Party
		assert.NotPanics(t, func() {
			norm.Normalize(p, tax.AddonContext(en16931.V2017))
		})
	})

	t.Run("accepts party with no inboxes", func(t *testing.T) {
		p := &org.Party{Name: "Acme"}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		assert.Empty(t, p.Endpoints, "no endpoints added when there is nothing to copy")
	})

	t.Run("copies peppol inbox into endpoint", func(t *testing.T) {
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{Key: org.InboxKeyPeppol, Scheme: "9920", Code: "x3157928m"},
			},
		}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		require.Len(t, p.Endpoints, 1)
		assert.Equal(t,
			"iso6523-actorid-upis::9920:x3157928m",
			p.Endpoints[0].URI.String(),
		)
		// Source inbox left intact for back-compat with existing consumers.
		require.Len(t, p.Inboxes, 1)
		assert.Equal(t, org.InboxKeyPeppol, p.Inboxes[0].Key)
		assert.Equal(t, "9920", p.Inboxes[0].Scheme.String())
		assert.Equal(t, "x3157928m", p.Inboxes[0].Code.String())
	})

	t.Run("copies the inbox label across", func(t *testing.T) {
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{Key: org.InboxKeyPeppol, Label: "Primary", Scheme: "9920", Code: "x3157928m"},
			},
		}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		require.Len(t, p.Endpoints, 1)
		assert.Equal(t, "Primary", p.Endpoints[0].Label)
	})

	t.Run("normalizes a combined code+scheme value before copying", func(t *testing.T) {
		// The inbox normalizer splits "9920:x3157928m" into Scheme + Code for
		// peppol-keyed inboxes; the norm engine walks children before parents,
		// so by the time the party normalizer runs it sees already-split
		// values from a single Normalize pass.
		in := &org.Inbox{Key: org.InboxKeyPeppol, Code: "9920:x3157928m"}
		p := &org.Party{Inboxes: []*org.Inbox{in}}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		require.Len(t, p.Endpoints, 1)
		assert.Equal(t, "iso6523-actorid-upis::9920:x3157928m", p.Endpoints[0].URI.String())
	})

	t.Run("skips non-peppol inboxes", func(t *testing.T) {
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{Key: "other", Scheme: "9920", Code: "x3157928m"},
			},
		}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		assert.Empty(t, p.Endpoints)
	})

	t.Run("skips peppol inboxes with missing scheme or code", func(t *testing.T) {
		// Without a Scheme we cannot build a valid iso6523 identifier.
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{Key: org.InboxKeyPeppol, Code: "x3157928m"},
			},
		}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		assert.Empty(t, p.Endpoints)

		// Same when Code is missing.
		p = &org.Party{
			Inboxes: []*org.Inbox{
				{Key: org.InboxKeyPeppol, Scheme: "9920"},
			},
		}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		assert.Empty(t, p.Endpoints)
	})

	t.Run("does not duplicate an existing iso6523 endpoint", func(t *testing.T) {
		p := &org.Party{
			Endpoints: []*org.Endpoint{
				{URI: "iso6523-actorid-upis::9920:already-here"},
			},
			Inboxes: []*org.Inbox{
				{Key: org.InboxKeyPeppol, Scheme: "9920", Code: "x3157928m"},
			},
		}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		require.Len(t, p.Endpoints, 1, "existing iso6523 endpoint blocks the copy")
		assert.Equal(t, "iso6523-actorid-upis::9920:already-here", p.Endpoints[0].URI.String())
	})

	t.Run("only copies the first peppol inbox", func(t *testing.T) {
		// Multiple peppol inboxes is unusual (BR-34/BR-49 says one inbox
		// per party for EN16931) but the normalizer must still produce
		// exactly one endpoint, not one per inbox.
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{Key: org.InboxKeyPeppol, Scheme: "9920", Code: "first"},
				{Key: org.InboxKeyPeppol, Scheme: "9920", Code: "second"},
			},
		}
		norm.Normalize(p, tax.AddonContext(en16931.V2017))
		require.Len(t, p.Endpoints, 1)
		assert.Equal(t, "iso6523-actorid-upis::9920:first", p.Endpoints[0].URI.String())
	})
}

func TestOrgItemValidate(t *testing.T) {
	t.Run("missing unit", func(t *testing.T) {
		item := &org.Item{Name: "Test"}
		err := rules.Validate(item, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "unit is required (BR-23)")
	})

	t.Run("validates unit", func(t *testing.T) {
		item := &org.Item{
			Name: "Test",
			Unit: org.UnitOne,
		}
		err := rules.Validate(item, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("negative price", func(t *testing.T) {
		item := &org.Item{
			Name:  "Test",
			Unit:  org.UnitOne,
			Price: num.NewAmount(-100, 0),
		}
		err := rules.Validate(item, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "zero or positive")
	})

	t.Run("0 price", func(t *testing.T) {
		item := &org.Item{
			Name:  "Test",
			Unit:  org.UnitOne,
			Price: num.NewAmount(0, 0),
		}
		err := rules.Validate(item, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})
}

func TestOrgAttachmentValidation(t *testing.T) {
	t.Run("blank", func(t *testing.T) {
		a := &org.Attachment{}
		err := rules.Validate(a, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "code is required")
	})
	t.Run("with code", func(t *testing.T) {
		a := &org.Attachment{
			Code: "123",
			URL:  "https://example.com/attachment.pdf",
		}
		err := rules.Validate(a, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})
}

func TestOrgPartyValidate(t *testing.T) {
	t.Run("no inboxes", func(t *testing.T) {
		p := &org.Party{}
		err := rules.Validate(p, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("one inbox", func(t *testing.T) {
		p := &org.Party{
			Inboxes: []*org.Inbox{
				{
					Scheme: "scheme1",
					Code:   "code1",
				},
			},
		}
		err := rules.Validate(p, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
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
		}
		err := rules.Validate(p, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "cannot have more than one inbox (BT-34, BT-49)")
	})
}

func TestOrgInboxValidate(t *testing.T) {
	t.Run("missing scheme and code", func(t *testing.T) {
		i := &org.Inbox{}
		// Not specific for addon, but this is important to check
		assert.ErrorContains(t, rules.Validate(i), "inbox requires a code, url, or email")
	})

	t.Run("missing scheme", func(t *testing.T) {
		i := &org.Inbox{
			Code: "code1",
		}
		err := rules.Validate(i, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "scheme cannot be blank when code is set (BR-62, BR-63)")
	})

	t.Run("missing code", func(t *testing.T) {
		i := &org.Inbox{
			Scheme: "scheme1",
		}
		err := rules.Validate(i, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "code cannot be blank when scheme is set")
	})

	t.Run("valid inbox", func(t *testing.T) {
		i := &org.Inbox{
			Scheme: "scheme1",
			Code:   "code1",
		}
		err := rules.Validate(i, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})
}
