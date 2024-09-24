package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Tax tags that can be applied in Australia.
// Spurce: https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes
const (
	TagWine      cbc.Key = "wine"
	TagLuxuryCar cbc.Key = "luxury-car"
	TagExcisable cbc.Key = "excisable-goods"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.KeyDefinition{
		{
			Key: tax.TagReverseCharge,
			Name: i18n.String{
				i18n.EN: "Reverse Charge",
			},
		},
		{
			Key: TagWine,
			Name: i18n.String{
				i18n.EN: "Wine Equalisation Tax",
			},
		},
		{
			Key: TagLuxuryCar,
			Name: i18n.String{
				i18n.EN: "Luxury Car Tax",
			},
		},
		{
			Key: TagExcisable,
			Name: i18n.String{
				i18n.EN: "Excisable Goods",
			},
		},
	},
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "A NEW TAX SYSTEM (GOODS AND SERVICES TAX) ACT 1999 - SECT 84.10 Reverse charge on offshore supplies",
			},
		},
		{
			Tags: []cbc.Key{TagWine},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagWine,
				Text: "Wine Equalisation Tax.",
			},
		},
		{
			Tags: []cbc.Key{TagLuxuryCar},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagLuxuryCar,
				Text: "Luxury Car Tax.",
			},
		},
		{
			Tags: []cbc.Key{TagExcisable},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagExcisable,
				Text: "Excise duties.",
			},
		},
	},
}
