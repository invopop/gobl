package ticket

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Official stamps or codes validated by government agencies
const (
	// StampRef is the key to identify the reference provided by the AdE once the ticket is accepted
	// This code needs to be printed on the ticket
	StampRef cbc.Key = "ade-ref"
	// StampVoidRef is the key to identify the reference provided by the AdE once the ticket is voided
	// This code needs to be printed on the ticket representing the void action
	StampVoidRef cbc.Key = "ade-void-ref"
)

var stampUsage = &cbc.Usage{
	Keys: []cbc.Key{V1},
}

var stamps = []*cbc.Definition{
	{
		Key: StampRef,
		Name: i18n.String{
			i18n.EN: "AdE Reference",
			i18n.IT: "Riferimento AdE",
		},
		Desc: i18n.String{
			i18n.EN: "Reference provided by the AdE once the ticket is accepted, to be printed on the ticket.",
			i18n.IT: "Riferimento fornito dall'AdE all'accettazione del documento commerciale, da stampare sul documento.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampVoidRef,
		Name: i18n.String{
			i18n.EN: "AdE Void Reference",
			i18n.IT: "Riferimento Annullo AdE",
		},
		Desc: i18n.String{
			i18n.EN: "Reference provided by the AdE once the ticket is voided, to be printed on the void ticket.",
			i18n.IT: "Riferimento fornito dall'AdE all'annullamento del documento commerciale, da stampare sul documento di annullo.",
		},
		Usage: stampUsage,
	},
}
