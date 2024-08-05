package hu

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
)

const (
	TagOutOfScope            cbc.Key = "out-of-scope"
	TagDomesticReverseCharge cbc.Key = "domestic-reverse-charge"
	TagTravelAgency          cbc.Key = "travel-agency"
	TagSecondHand            cbc.Key = "second-hand"
	TagArt                   cbc.Key = "art"
	TagAntiques              cbc.Key = "antiques"
)

var invoiceTags = common.InvoiceTagsWith([]*cbc.KeyDefinition{
	{
		Key: TagOutOfScope,
		Name: i18n.String{
			i18n.EN: "Out of Scope",
			i18n.HU: "Hatályaon kívül",
		},
	},
	{
		Key: TagDomesticReverseCharge,
		Name: i18n.String{
			i18n.EN: "Domestic Reverse Charge",
			i18n.HU: "Belföldi fordított adózás",
		},
	},
	{
		Key: TagTravelAgency,
		Name: i18n.String{
			i18n.EN: "Travel Agency",
			i18n.HU: "Utazási iroda",
		},
	},
	{
		Key: TagSecondHand,
		Name: i18n.String{
			i18n.EN: "Second Hand",
			i18n.HU: "Használt cikk",
		},
	},
	{
		Key: TagArt,
		Name: i18n.String{
			i18n.EN: "Art",
			i18n.HU: "Műalkotás",
		},
	},
	{
		Key: TagAntiques,
		Name: i18n.String{
			i18n.EN: "Antiques",
			i18n.HU: "Antikvitás",
		},
	},
})
