package tbai

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

var stamps = []*cbc.Definition{
	{
		Key: StampCode,
		Name: i18n.String{
			i18n.EN: "TicketBAI Code",
			i18n.ES: "Código TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: "TicketBAI identifier code to be presented alongside the QR code.",
		},
	},
	{
		Key: StampQR,
		Name: i18n.String{
			i18n.EN: "TicketBAI QR Code",
			i18n.ES: "Código QR TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: "URL included in the QR code.",
		},
	},
}
