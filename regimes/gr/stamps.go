package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

var stamps = []*cbc.Definition{
	{
		Key: StampIAPRMark,
		Name: i18n.String{
			i18n.EN: "MARK",
			i18n.EL: "ΜΑΡΚ",
		},
		Desc: i18n.String{
			i18n.EN: "Unique registration number assigned by myDATA to the document. Must be copied to preceding document references when correcting.",
		},
	},
	{
		Key: StampIAPRUID,
		Name: i18n.String{
			i18n.EN: "Unique ID",
		},
		Desc: i18n.String{
			i18n.EN: "Unique identifier (UID) of the document.",
		},
	},
	{
		Key: StampIAPRHash,
		Name: i18n.String{
			i18n.EN: "Authentication Code",
			i18n.EL: "Σήμανση",
		},
		Desc: i18n.String{
			i18n.EN: "Authentication code of the document, to be printed on it.",
		},
	},
	{
		Key: StampIAPRQR,
		Name: i18n.String{
			i18n.EN: "QR Code URL",
		},
		Desc: i18n.String{
			i18n.EN: "URL to include in the QR code printed on the document.",
		},
	},
	{
		Key: StampIAPRProvider,
		Name: i18n.String{
			i18n.EN: "E-Invoicing Provider",
			i18n.EL: "Υ.ΠΑ.Η.Ε.Σ",
		},
		Desc: i18n.String{
			i18n.EN: "Electronic invoicing service provider that issued the document.",
		},
	},
}
