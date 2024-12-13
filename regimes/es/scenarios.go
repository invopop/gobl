package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Tax tags that can be applied in Spain.
const (
	TagCopy             cbc.Key = "copy"
	TagSummary          cbc.Key = "summary"
	TagSimplifiedScheme cbc.Key = "simplified-scheme"
	TagCustomerIssued   cbc.Key = "customer-issued"
	TagTravelAgency     cbc.Key = "travel-agency"
	TagSecondHandGoods  cbc.Key = "second-hand-goods"
	TagArt              cbc.Key = "art"
	TagAntiques         cbc.Key = "antiques"
	TagCashBasis        cbc.Key = "cash-basis"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		// Copy of the original document
		{
			Key: TagCopy,
			Name: i18n.String{
				i18n.EN: "Copy",
				i18n.ES: "Copia",
			},
		},
		// Summary document
		{
			Key: TagSummary,
			Name: i18n.String{
				i18n.EN: "Summary",
				i18n.ES: "Recapitulativa",
			},
		},
		// Simplified Scheme (Modules)
		{
			Key: TagSimplifiedScheme,
			Name: i18n.String{
				i18n.EN: "Simplified tax scheme",
				i18n.ES: "Contribuyente en régimen simplificado",
			},
		},

		// Travel agency
		{
			Key: TagTravelAgency,
			Name: i18n.String{
				i18n.EN: "Special scheme for travel agencies",
				i18n.ES: "Régimen especial de las agencias de viajes",
			},
		},
		// Secondhand stuff
		{
			Key: TagSecondHandGoods,
			Name: i18n.String{
				i18n.EN: "Special scheme for second-hand goods",
				i18n.ES: "Régimen especial de los bienes usados",
			},
		},
		// Art
		{
			Key: TagArt,
			Name: i18n.String{
				i18n.EN: "Special scheme of works of art",
				i18n.ES: "Régimen especial de los objetos de arte",
			},
		},
		// Antiques
		{
			Key: TagAntiques,
			Name: i18n.String{
				i18n.EN: "Special scheme of antiques and collectables",
				i18n.ES: "Régimen especial de las antigüedades y objetos de colección",
			},
		},
		// Special Regime of "Cash Criteria"
		{
			Key: TagCashBasis,
			Name: i18n.String{
				i18n.EN: "Special scheme on cash basis",
				i18n.ES: "Régimen especial del criterio de caja",
			},
		},
	},
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// ** Special Messages **
		// Reverse Charges
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse Charge / Inversión del sujeto pasivo.",
			},
		},
		// Simplified Scheme (Modules)
		{
			Tags: []cbc.Key{TagSimplifiedScheme},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagSimplifiedScheme,
				Text: "Factura expedida por contibuyente en régimen simplificado.",
			},
		},
		// Customer issued invoices
		{
			Tags: []cbc.Key{tax.TagSelfBilled},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSelfBilled,
				Text: "Facturación por el destinatario.",
			},
		},
		// Travel agency
		{
			Tags: []cbc.Key{TagTravelAgency},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagTravelAgency,
				Text: "Régimen especial de las agencias de viajes.",
			},
		},
		// Secondhand stuff
		{
			Tags: []cbc.Key{TagSecondHandGoods},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagSecondHandGoods,
				Text: "Régimen especial de los bienes usados.",
			},
		},
		// Art
		{
			Tags: []cbc.Key{TagArt},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagArt,
				Text: "Régimen especial de los objetos de arte.",
			},
		},
		// Antiques
		{
			Tags: []cbc.Key{TagAntiques},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagAntiques,
				Text: "Régimen especial de las antigüedades y objetos de colección.",
			},
		},
		// Special Regime of "Cash Criteria"
		{
			Tags: []cbc.Key{TagCashBasis},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagCashBasis,
				Text: "Régimen especial del criterio de caja.",
			},
		},
	},
}
