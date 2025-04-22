// Package no provides tax scenarios specific to NO VAT regulations.
package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Document tag keys - define any Norway-specific document tags here
var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagSecondHand,
			Name: i18n.String{
				i18n.EN: "Second-hand Goods",
				i18n.NO: "Brukte varer",
			},
		},
		{
			Key: TagBooks,
			Name: i18n.String{
				i18n.EN: "Books and Periodicals",
				i18n.NO: "Bøker og tidsskrifter",
			},
		},
		{
			Key: TagECommerce,
			Name: i18n.String{
				i18n.EN: "VOEC E-commerce",
				i18n.NO: "VOEC E-handel",
			},
		},
		{
			Key: TagArtworks,
			Name: i18n.String{
				i18n.EN: "Works of Art",
				i18n.NO: "Kunstgjenstander",
			},
		},
	},
}

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charge
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Name: i18n.String{
				i18n.EN: "Reverse Charge",
				i18n.NO: "Omvendt avgiftsplikt",
			},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse Charge",
			},
		},
		// Simplified Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Name: i18n.String{
				i18n.EN: "Simplified Invoice",
				i18n.NO: "Forenklet faktura",
			},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Simplified Invoice (for transactions below NOK 1,000)",
			},
		},
		// Zero-Rated Export
		{
			Tags: []cbc.Key{tax.RateZero},
			Name: i18n.String{
				i18n.EN: "Zero-Rated Export",
				i18n.NO: "Nullsatseksport",
			},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.RateZero,
				Text: "Zero-Rated Export",
			},
		},
		// Second-hand goods margin scheme
		{
			Tags: []cbc.Key{TagSecondHand},
			Name: i18n.String{
				i18n.EN: "Second-hand Goods Margin Scheme",
				i18n.NO: "Avansesystem for brukte varer",
			},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagSecondHand,
				Text: "Margin Scheme - Second-hand Goods (Chapter Va MVAL)",
			},
		},
		// Artworks margin scheme
		{
			Tags: []cbc.Key{TagArtworks},
			Name: i18n.String{
				i18n.EN: "Works of Art Margin Scheme",
				i18n.NO: "Avansesystem for kunstgjenstander",
			},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagArtworks,
				Text: "Margin Scheme - Works of Art (Chapter Va MVAL)",
			},
		},
		// E-commerce VOEC scheme
		{
			Tags: []cbc.Key{TagECommerce},
			Name: i18n.String{
				i18n.EN: "VOEC E-commerce",
				i18n.NO: "VOEC E-handel",
			},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagECommerce,
				Text: "VOEC Scheme - B2C E-commerce (§ 3-30 MVAL)",
			},
		},
		// Books and periodicals zero-rating
		{
			Tags: []cbc.Key{TagBooks},
			Name: i18n.String{
				i18n.EN: "Books and Periodicals",
				i18n.NO: "Bøker og tidsskrifter",
			},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagBooks,
				Text: "Zero-Rated Books and Periodicals (§ 6-4 MVAL)",
			},
		},
	},
}
