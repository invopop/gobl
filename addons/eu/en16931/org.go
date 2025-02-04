package en16931

import (
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

func validateOrgItem(item *org.Item) error {
	return validation.ValidateStruct(item,
		validation.Field(&item.Unit,
			validation.Required.Error("cannot be blank (BR-23)"),
		),
	)
}
