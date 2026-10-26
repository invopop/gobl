package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Official IAPR codes to include in stamps.
const (
	StampIAPRQR       cbc.Key = "iapr-qr"
	StampIAPRMark     cbc.Key = "iapr-mark"
	StampIAPRHash     cbc.Key = "iapr-hash"
	StampIAPRUID      cbc.Key = "iapr-uid"
	StampIAPRProvider cbc.Key = "iapr-provider"
)

var stampUsage = &cbc.Usage{
	Codes: []cbc.Code{CountryCode},
}

var stamps = []*cbc.Definition{
	{
		Key: StampIAPRMark,
		Name: i18n.String{
			i18n.EN: "IAPR: MARK",
			i18n.EL: "ΑΑΔΕ: ΜΑΡΚ",
		},
		Desc: i18n.String{
			i18n.EN: "Unique registration number assigned by myDATA to the document.",
			i18n.EL: "Μοναδικός αριθμός καταχώρησης που αποδίδεται από το myDATA στο παραστατικό.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampIAPRUID,
		Name: i18n.String{
			i18n.EN: "IAPR: Unique ID",
			i18n.EL: "ΑΑΔΕ: Μοναδικός Κωδικός",
		},
		Desc: i18n.String{
			i18n.EN: "Unique identifier of the document registered in myDATA.",
			i18n.EL: "Μοναδικός κωδικός του παραστατικού που καταχωρήθηκε στο myDATA.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampIAPRHash,
		Name: i18n.String{
			i18n.EN: "IAPR: Authentication Code",
			i18n.EL: "ΑΑΔΕ: Σήμανση",
		},
		Desc: i18n.String{
			i18n.EN: "Authentication code of the document issued through an e-invoicing provider.",
			i18n.EL: "Κωδικός αυθεντικοποίησης του παραστατικού που εκδόθηκε μέσω παρόχου ηλεκτρονικής τιμολόγησης.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampIAPRProvider,
		Name: i18n.String{
			i18n.EN: "IAPR: E-Invoicing Provider",
			i18n.EL: "ΑΑΔΕ: Πάροχος Ηλεκτρονικής Τιμολόγησης",
		},
		Desc: i18n.String{
			i18n.EN: "E-invoicing provider that issued the document.",
			i18n.EL: "Πάροχος ηλεκτρονικής τιμολόγησης που εξέδωσε το παραστατικό.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampIAPRQR,
		Name: i18n.String{
			i18n.EN: "IAPR: QR Code",
			i18n.EL: "ΑΑΔΕ: Κωδικός QR",
		},
		Desc: i18n.String{
			i18n.EN: "URL included in the QR code to verify the document.",
			i18n.EL: "URL που περιλαμβάνεται στον κωδικό QR για την επαλήθευση του παραστατικού.",
		},
		Usage: stampUsage,
	},
}
