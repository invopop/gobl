package en16931

import (
	"regexp"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/dk"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Map of GOBL Note keys to the corresponding UNTDID 4451 code.
var orgNoteTextSubjectMap = map[cbc.Key]cbc.Code{
	org.NoteKeyGoods:          "AAA",
	org.NoteKeyPayment:        "PMT",
	org.NoteKeyPaymentMethod:  "PMD",
	org.NoteKeyPaymentTerm:    "AAB",
	org.NoteKeyGeneral:        "AAI", // General information
	org.NoteKeyLegal:          "ABY",
	org.NoteKeyDangerousGoods: "AAC",
	org.NoteKeyAck:            "AAE",
	org.NoteKeyRate:           "AAF",
	org.NoteKeyReason:         "ACD",
	org.NoteKeyDispute:        "ACE",
	org.NoteKeyCustomer:       "CUR",
	org.NoteKeyGlossary:       "ACZ",
	org.NoteKeyCustoms:        "CUS",
	org.NoteKeyHandling:       "HAN",
	org.NoteKeyPackaging:      "PKG",
	org.NoteKeyLoading:        "LOI",
	org.NoteKeyPrice:          "AAK",
	org.NoteKeyPriority:       "PRI",
	org.NoteKeyRegulatory:     "REG",
	org.NoteKeySafety:         "SAF",
	org.NoteKeyShipLine:       "SLR",
	org.NoteKeySupplier:       "SUR",
	org.NoteKeyTransport:      "TRA",
	org.NoteKeyDelivery:       "DEL",
	org.NoteKeyQuarantine:     "QIN",
	org.NoteKeyTax:            "TXD",
	org.NoteKeyOther:          "ZZZ",
}

// Map of GOBL Identity keys to the corresponding ISO/IEC 6523 code.
var orgIdentitySchemeMap = map[cbc.Key]cbc.Code{
	org.IdentityKeyGLN:  "0088",
	org.IdentityKeyGTIN: "0160",
}

// Map of GOBL Identity types (codes) to the corresponding ISO/IEC 6523 code.
// This is used for regime-specific identity types that use Type instead of Key.
var orgIdentityTypeSchemeMap = map[cbc.Code]cbc.Code{
	fr.IdentityTypeSIREN: "0002", // French SIREN (legal identifier)
	fr.IdentityTypeSIRET: "0009", // French SIRET (private identifier)
	dk.IdentityTypeCVR:   "0184", // Danish CVR-nummer
}

var (
	orgInboxRegexpSchemeCode = regexp.MustCompile(`(\d{4}):.*`)
)

func normalizeOrgNote(n *org.Note) {
	if n == nil {
		return
	}
	if n.Key == cbc.KeyEmpty {
		return
	}
	if code, ok := orgNoteTextSubjectMap[n.Key]; ok {
		n.Ext = n.Ext.Merge(tax.Extensions{
			untdid.ExtKeyTextSubject: code,
		})
	}
}

func normalizeOrgItem(item *org.Item) {
	if item == nil {
		return
	}
	if item.Unit == org.UnitEmpty {
		item.Unit = org.UnitOne
	}
}

func normalizeOrgIdentity(i *org.Identity) {
	if i == nil {
		return
	}

	// Check for key-based identity mapping first
	if i.Key != cbc.KeyEmpty {
		if scheme, ok := orgIdentitySchemeMap[i.Key]; ok {
			i.Ext = i.Ext.Merge(tax.Extensions{
				iso.ExtKeySchemeID: scheme,
			})
			return
		}
	}

	// Check for type-based identity mapping (used by some regimes like France)
	if i.Type != cbc.CodeEmpty {
		if scheme, ok := orgIdentityTypeSchemeMap[i.Type]; ok {
			i.Ext = i.Ext.Merge(tax.Extensions{
				iso.ExtKeySchemeID: scheme,
			})
		}
	}
}

func normalizeOrgInbox(i *org.Inbox) {
	if i == nil || i.Code == cbc.CodeEmpty {
		return
	}
	if orgInboxRegexpSchemeCode.MatchString(i.Code.String()) {
		i.Scheme = cbc.Code(i.Code.String()[0:4])
		i.Code = cbc.Code(i.Code.String()[5:])
	}
}

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.Field("unit",
			// BR-23: unit of measure is required
			rules.Assert("01", "unit is required (BR-23)", is.Present),
		),
	)
}

func orgAttachmentRules() *rules.Set {
	return rules.For(new(org.Attachment),
		rules.Field("code",
			rules.Assert("01", "code is required", is.Present),
		),
	)
}

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		rules.Field("inboxes",
			rules.Assert("01", "cannot have more than one inbox (BT-34, BT-49)",
				is.Length(0, 1),
			),
		),
	)
}

func orgInboxRules() *rules.Set {
	return rules.For(new(org.Inbox),
		// BR-62, BR-63: scheme required when code is present
		rules.Assert("01", "scheme cannot be blank when code is set (BR-62, BR-63)",
			is.Func("scheme required with code", orgInboxSchemeRequiredWithCode),
		),
		// code required when scheme is present
		rules.Assert("02", "code cannot be blank when scheme is set",
			is.Func("code required with scheme", orgInboxCodeRequiredWithScheme),
		),
	)
}

func orgInboxSchemeRequiredWithCode(val any) bool {
	i, ok := val.(*org.Inbox)
	return !ok || i == nil || i.Code == cbc.CodeEmpty || i.Scheme != cbc.CodeEmpty
}

func orgInboxCodeRequiredWithScheme(val any) bool {
	i, ok := val.(*org.Inbox)
	return !ok || i == nil || i.Scheme == cbc.CodeEmpty || i.Code != cbc.CodeEmpty
}

func orgAddressRules() *rules.Set {
	return rules.For(new(org.Address),
		rules.Field("country",
			// Most addresses in EN16931 need a country: BR-9, BR-11, BR-20, BR-57
			rules.Assert("01", "country is required (BR-9, BR-11, BR-20, BR-57)", is.Present),
		),
	)
}
