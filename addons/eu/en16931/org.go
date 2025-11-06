package en16931

import (
	"regexp"

	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Map of GOBL keys to the corresponding UNTDID 4451 code.
var orgNoteTextSubjectMap = map[cbc.Key]cbc.Code{
	org.NoteKeyGoods:          "AAA",
	org.NoteKeyPayment:        "PMT",
	org.NoteKeyLegal:          "ABY",
	org.NoteKeyDangerousGoods: "AAC",
	org.NoteKeyAck:            "AAE",
	org.NoteKeyRate:           "AAF",
	org.NoteKeyReason:         "ACD",
	org.NoteKeyDispute:        "ACE",
	org.NoteKeyCustomer:       "CUR",
	org.NoteKeyGlossary:       "ACZ",
	org.NoteKeyCustoms:        "CUS",
	org.NoteKeyGeneral:        "AAI",
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

func normalizeOrgInbox(i *org.Inbox) {
	if i == nil || i.Code == cbc.CodeEmpty {
		return
	}
	if orgInboxRegexpSchemeCode.MatchString(i.Code.String()) {
		i.Scheme = cbc.Code(i.Code.String()[0:4])
		i.Code = cbc.Code(i.Code.String()[5:])
	}
}

func validateOrgItem(item *org.Item) error {
	return validation.ValidateStruct(item,
		validation.Field(&item.Unit,
			validation.Required.Error("cannot be blank (BR-23)"),
		),
	)
}

func validateOrgAttachment(a *org.Attachment) error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Code,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateOrgParty(p *org.Party) error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Inboxes,
			validation.Length(0, 1).Error("cannot have more than one inbox (BT-34, BT-49)"),
			validation.Skip,
		),
	)
}

func validateOrgInbox(i *org.Inbox) error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Scheme,
			validation.When(i.Code != cbc.CodeEmpty,
				validation.Required.Error("cannot be blank with code (BR-62, BR-63)"),
			),
			validation.Skip,
		),
		validation.Field(&i.Code,
			validation.When(i.Scheme != cbc.CodeEmpty,
				validation.Required.Error("cannot be blank with scheme"),
			),
			validation.Skip,
		),
	)
}

func validateOrgAddress(a *org.Address) error {
	return validation.ValidateStruct(a,
		// Most addresses in EN16931 need a country: BR-9, BR-11, BR-20, BR-57
		validation.Field(&a.Country,
			validation.Required,
			validation.Skip,
		),
	)
}
