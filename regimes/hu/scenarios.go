package hu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Document tag keys
const (
	TagDomesticReverseCharge cbc.Key = "domestic-reverse-charge"
	TagTravelAgency          cbc.Key = "travel-agency"
	TagSecondHand            cbc.Key = "second-hand"
	TagArt                   cbc.Key = "art"
	TagAntiques              cbc.Key = "antiques"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Name: i18n.String{
				i18n.EN: "Standard invoice",
				i18n.HU: "Számla",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Name: i18n.String{
				i18n.EN: "Credit note",
				i18n.HU: "Jóváírás",
			},
		},
	},
}

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.KeyDefinition{
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
	},
}
