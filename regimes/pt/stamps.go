package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

var stamps = []*cbc.Definition{
	{
		Key: StampProviderATATCUD,
		Name: i18n.String{
			i18n.EN: "ATCUD",
			i18n.PT: "ATCUD",
		},
		Desc: i18n.String{
			i18n.EN: "Unique document code, made of the series validation code and the document's sequential number, to be printed on the document.",
		},
	},
	{
		Key: StampProviderATQR,
		Name: i18n.String{
			i18n.EN: "AT QR Code",
			i18n.PT: "Código QR",
		},
		Desc: i18n.String{
			i18n.EN: "Contents of the QR code to be printed on the document.",
		},
	},
	{
		Key: StampProviderATHash,
		Name: i18n.String{
			i18n.EN: "Signature Hash",
		},
		Desc: i18n.String{
			i18n.EN: "Characters extracted from the document's signature hash, to be printed on the document.",
		},
	},
	{
		Key: StampProviderATHashFull,
		Name: i18n.String{
			i18n.EN: "Full Signature Hash",
		},
		Desc: i18n.String{
			i18n.EN: "Complete signature hash of the document.",
		},
	},
	{
		Key: StampProviderATAppID,
		Name: i18n.String{
			i18n.EN: "Application ID",
		},
		Desc: i18n.String{
			i18n.EN: "Identifier of the certified application used to register the document.",
		},
	},
	{
		Key: StampProviderATTimestamp,
		Name: i18n.String{
			i18n.EN: "Timestamp",
		},
		Desc: i18n.String{
			i18n.EN: "Date and time at which the document was registered.",
		},
	},
	{
		Key: StampProviderATDocCode,
		Name: i18n.String{
			i18n.EN: "AT Document Code",
		},
		Desc: i18n.String{
			i18n.EN: "Identification code assigned by the AT when the document is communicated.",
		},
	},
}
