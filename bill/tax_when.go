package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Tax When keys define when the tax point occurs, i.e. when the tax liability
// is triggered. These correspond to UNCL 2005 date/time/period function codes
// used in BT-8 of the EN 16931 standard.
const (
	// TaxWhenIssue indicates tax is due on the invoice issue date.
	// Corresponds to UNCL 2005 code 3.
	TaxWhenIssue cbc.Key = "issue"

	// TaxWhenDelivery indicates tax is due on the delivery date.
	// Corresponds to UNCL 2005 code 35.
	TaxWhenDelivery cbc.Key = "delivery"

	// TaxWhenPaid indicates tax is due when payment is received.
	// Corresponds to UNCL 2005 code 432.
	TaxWhenPaid cbc.Key = "paid"
)

// TaxWhenDefs lists the supported tax when keys and their descriptions.
var TaxWhenDefs = []*cbc.Definition{
	{
		Key: TaxWhenIssue,
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
		Key: TaxWhenDelivery,
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
		Key: TaxWhenPaid,
		Name: i18n.String{
			i18n.EN: "Paid Date",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Tax point is the date on which payment is made. Corresponds to UNCL 2005 code 432.
			`),
		},
	},
}
