package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Point keys define when the tax point occurs, i.e. when the tax liability
// is triggered. These correspond to UNCL 2005 date/time/period function codes
// used in BT-8 of the EN 16931 standard.
const (
	// PointIssue indicates tax is due on the invoice issue date.
	// Corresponds to UNCL 2005 code 3.
	PointIssue cbc.Key = "issue"

	// PointDelivery indicates tax is due on the delivery date.
	// Corresponds to UNCL 2005 code 35.
	PointDelivery cbc.Key = "delivery"

	// PointPayment indicates tax is due when payment is received.
	// Corresponds to UNCL 2005 code 432.
	PointPayment cbc.Key = "payment"
)

// PointDefs lists the supported tax point keys and their descriptions.
var PointDefs = []*cbc.Definition{
	{
		Key: PointIssue,
		Name: i18n.String{
			i18n.EN: "Issue Date",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Tax point is the invoice issue date. Corresponds to UNCL 2005 code 3.
			`),
		},
	},
	{
		Key: PointDelivery,
		Name: i18n.String{
			i18n.EN: "Delivery Date",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Tax point is the actual delivery date. Corresponds to UNCL 2005 code 35.
			`),
		},
	},
	{
		Key: PointPayment,
		Name: i18n.String{
			i18n.EN: "Payment Date",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Tax point is the date on which payment is made. Corresponds to UNCL 2005 code 432.
			`),
		},
	},
}
