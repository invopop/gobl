package ticket

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

var stamps = []*cbc.Definition{
	{
		Key: StampRef,
		Name: i18n.String{
			i18n.EN: "AdE Reference",
		},
		Desc: i18n.String{
			i18n.EN: "Reference provided by the AdE once the ticket is accepted. Must be printed on the ticket and copied to preceding document references when correcting.",
		},
	},
	{
		Key: StampVoidRef,
		Name: i18n.String{
			i18n.EN: "AdE Void Reference",
		},
		Desc: i18n.String{
			i18n.EN: "Reference provided by the AdE once the ticket is voided. Must be printed on the ticket representing the void action.",
		},
	},
}
