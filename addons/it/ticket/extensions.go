package ticket

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Italian extension keys required by the AdE ticket format.
const (
	ExtKeyExempt  cbc.Key = "it-ticket-exempt"
	ExtKeyProduct cbc.Key = "it-ticket-product"
)

var extensions = []*cbc.Definition{
	{
		// Used to clarify the reason for the exemption from VAT.
		Key: ExtKeyExempt,
		Name: i18n.String{
			i18n.EN: "Exemption Code",
			i18n.IT: "Natura Esenzione",
		},
		Values: []*cbc.Definition{
			{
				Code: "N1",
				Name: i18n.String{
					i18n.EN: "Excluded pursuant to Art. 15, DPR 633/72",
					i18n.IT: "Escluse ex. art. 15 del D.P.R. 633/1972",
				},
			},
			{
				Code: "N2",
				Name: i18n.String{
					i18n.EN: "Not subject",
					i18n.IT: "Non soggette",
				},
			},
			{
				Code: "N3",
				Name: i18n.String{
					i18n.EN: "Not taxable",
					i18n.IT: "Non imponibili",
				},
			},
			{
				Code: "N4",
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.IT: "Esenti",
				},
			},
			{
				Code: "N5",
				Name: i18n.String{
					i18n.EN: "Margin regime / VAT not exposed",
					i18n.IT: "Regime del margine/IVA non esposta in fattura",
				},
			},
			{
				Code: "N6",
				Name: i18n.String{
					i18n.EN: "Reverse charge",
					i18n.IT: "Inversione contabile",
				},
			},
		},
	},
	{
		Key: ExtKeyProduct,
		Name: i18n.String{
			i18n.EN: "AdE CF Product Key",
			i18n.IT: "Chiave Prodotto AdE CF",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Product keys are used by AdE CF to differentiate between goods
				and services.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "goods",
				Name: i18n.String{
					i18n.EN: "Delivery of goods",
					i18n.IT: "Consegna di beni",
				},
			},
			{
				Code: "services",
				Name: i18n.String{
					i18n.EN: "Provision of services",
					i18n.IT: "Prestazione di servizi",
				},
			},
		},
	},
}
