package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Scheme key definitions
const (
	SchemeSimplified      cbc.Key = "simplified"
	SchemeCustomerIssued  cbc.Key = "customer-issued"
	SchemeTravelAgency    cbc.Key = "travel-agency"
	SchemeSecondHandGoods cbc.Key = "second-hand-goods"
	SchemeArt             cbc.Key = "art"
	SchemeAntiques        cbc.Key = "antiques"
	SchemeCashBasis       cbc.Key = "cash-basis"
)

var schemes = []*tax.Scheme{
	// Reverse Charge Scheme
	{
		Key: common.SchemeReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.ES: "Inversión del sujeto pasivo",
		},
		Categories: []cbc.Code{
			common.TaxCategoryVAT,
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(common.SchemeReverseCharge),
			Text: "Reverse Charge / Inversión del sujeto pasivo.",
		},
	},
	// Customer Rates Scheme (digital goods)
	{
		Key: common.SchemeCustomerRates,
		Name: i18n.String{
			i18n.EN: "Customer Country Rates",
			i18n.ES: "Tasas del País del Cliente",
		},
		Description: i18n.String{
			i18n.EN: "Use the customers country to determine tax rates.",
		},
	},
	// Simplified Regime
	{
		Key: SchemeSimplified,
		Name: i18n.String{
			i18n.EN: "Simplified tax scheme",
			i18n.ES: "Contribuyente en régimen simplificado",
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(SchemeSimplified),
			Text: "Factura expedida por contibuyente en régimen simplificado.",
		},
	},
	// Customer issued invoices
	{
		Key: SchemeCustomerIssued,
		Name: i18n.String{
			i18n.EN: "Customer issued invoice",
			i18n.ES: "Facturación por el destinatario",
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(SchemeCustomerIssued),
			Text: "Facturación por el destinatario.",
		},
	},
	// Travel agency
	{
		Key: SchemeTravelAgency,
		Name: i18n.String{
			i18n.EN: "Special scheme for travel agencies",
			i18n.ES: "Régimen especial de las agencias de viajes",
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(SchemeTravelAgency),
			Text: "Régimen especial de las agencias de viajes.",
		},
	},
	// Secondhand stuff
	{
		Key: SchemeSecondHandGoods,
		Name: i18n.String{
			i18n.EN: "Special scheme for second-hand goods",
			i18n.ES: "Régimen especial de los bienes usados",
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(SchemeSecondHandGoods),
			Text: "Régimen especial de los bienes usados.",
		},
	},
	// Art
	{
		Key: SchemeArt,
		Name: i18n.String{
			i18n.EN: "Special scheme of works of art",
			i18n.ES: "Régimen especial de los objetos de arte",
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(SchemeArt),
			Text: "Régimen especial de los objetos de arte.",
		},
	},
	// Antiques
	{
		Key: SchemeAntiques,
		Name: i18n.String{
			i18n.EN: "Special scheme of antiques and collectables",
			i18n.ES: "Régimen especial de las antigüedades y objetos de colección",
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(SchemeAntiques),
			Text: "Régimen especial de las antigüedades y objetos de colección.",
		},
	},
	// Special Regime of "Cash Criteria"
	{
		Key: SchemeCashBasis,
		Name: i18n.String{
			i18n.EN: "Special scheme on cash basis",
			i18n.ES: "Régimen especial del criterio de caja",
		},
		Note: &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  string(SchemeCashBasis),
			Text: "Régimen especial del criterio de caja.",
		},
	},
}
