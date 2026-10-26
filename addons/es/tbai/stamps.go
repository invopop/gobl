package tbai

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Official stamps or codes validated by government agencies
const (
	// StampCode contains the code required to be presented alongside
	// the QR code.
	StampCode cbc.Key = "tbai-code"
	// StampQR contains the URL included in the QR code.
	StampQR cbc.Key = "tbai-qr"
)

var stampUsage = &cbc.Usage{
	Keys: []cbc.Key{V1},
}

var stamps = []*cbc.Definition{
	{
		Key: StampCode,
		Name: i18n.String{
			i18n.EN: "TicketBAI Code",
			i18n.ES: "Código TicketBAI",
			i18n.EU: "TicketBAI Kodea",
		},
		Desc: i18n.String{
			i18n.EN: "Identifier of the TicketBAI invoice that must be presented alongside the QR code.",
			i18n.ES: "Identificador de la factura TicketBAI que debe presentarse junto al código QR.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampQR,
		Name: i18n.String{
			i18n.EN: "QR Code",
			i18n.ES: "Código QR",
			i18n.EU: "QR Kodea",
		},
		Desc: i18n.String{
			i18n.EN: "URL included in the QR code to verify the invoice with the tax agency.",
			i18n.ES: "URL incluida en el código QR para verificar la factura ante la hacienda foral.",
		},
		Usage: stampUsage,
	},
}
