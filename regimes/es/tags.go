package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Universal tax tags
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

// Tax tags which may be used in the Basque Country.
const (
	TagProvider  cbc.Key = "provider"
	TagServices  cbc.Key = "services"
	TagGoods     cbc.Key = "goods"
	TagExempt    cbc.Key = "exempt"
	TagArticle20 cbc.Key = "article-20"
	TagArticle21 cbc.Key = "article-21"
	TagArticle22 cbc.Key = "article-22"
	TagArticle23 cbc.Key = "article-23"
	TagArticle25 cbc.Key = "article-25"
	TagOther     cbc.Key = "other"
)

var invoiceTags = []*tax.Tag{
	// Simplified Invoice
	{
		Key: common.TagSimplified,
		Name: i18n.String{
			i18n.EN: "Simplified Invoice",
			i18n.ES: "Factura Simplificada",
		},
	},
	// Customer rates (mainly for digital goods inside EU)
	{
		Key: common.TagCustomerRates,
		Name: i18n.String{
			i18n.EN: "Customer rates",
			i18n.ES: "Tarifas aplicables al destinatario",
		},
	},
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
	// Reverse Charge Mechanism
	{
		Key: common.TagReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.ES: "Inversión del sujeto pasivo",
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
	// Customer issued invoices
	{
		Key: common.TagSelfBilled,
		Name: i18n.String{
			i18n.EN: "Customer issued invoice",
			i18n.ES: "Facturación por el destinatario",
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
}

var commonVATTags = []*tax.Tag{
	{
		Key: TagProvider,
		Name: i18n.String{
			i18n.ES: "Operacion en recargo de equivalencia o regimen simplificado",
			i18n.EN: "Operation in equivalence surcharge or simplified regime",
		},
	},
	{
		Key: TagServices,
		Name: i18n.String{
			i18n.ES: "Prestacion de servicios",
			i18n.EN: "Provision of services",
		},
	},
	{
		Key: TagGoods,
		Name: i18n.String{
			i18n.ES: "Entrega de bienes",
			i18n.EN: "Delivery of goods",
		},
	},
}

var exemptTaxTags = []*tax.Tag{
	{
		Key: TagExempt.With(TagArticle20),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to Article 20 of the Foral VAT Law",
			i18n.ES: "Exenta por el artículo 20 de la Norma Foral del IVA",
		},
		Meta: cbc.Meta{
			KeyTicketBAICausaExencion: "E1",
		},
	},
	{
		Key: TagExempt.With(TagArticle21),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to Article 21 of the Foral VAT Law",
			i18n.ES: "Exenta por el artículo 21 de la Norma Foral del IVA",
		},
		Meta: cbc.Meta{
			KeyTicketBAICausaExencion: "E2",
		},
	},
	{
		Key: TagExempt.With(TagArticle22),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to Article 22 of the Foral VAT Law",
			i18n.ES: "Exenta por el artículo 22 de la Norma Foral del IVA",
		},
		Meta: cbc.Meta{
			KeyTicketBAICausaExencion: "E3",
		},
	},
	{
		Key: TagExempt.With(TagArticle23),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to Articles 23 and 24 of the Foral VAT Law",
			i18n.ES: "Exenta por el artículos 23 y 24 de la Norma Foral del IVA",
		},
		Meta: cbc.Meta{
			KeyTicketBAICausaExencion: "E4",
		},
	},
	{
		Key: TagExempt.With(TagArticle25),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to Article 25 of the Foral VAT law",
			i18n.ES: "Exenta por el artículo 25 de la Norma Foral del IVA",
		},
		Meta: cbc.Meta{
			KeyTicketBAICausaExencion: "E5",
		},
	},
	{
		Key: TagExempt.With(TagOther),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to other reasons",
			i18n.ES: "Exenta por otra causa",
		},
		Meta: cbc.Meta{
			KeyTicketBAICausaExencion: "E6",
		},
	},
}
