package adecf

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Italian extension keys required by the AdE adecf format.
const (
	ExtKeyExempt cbc.Key = "it-adecf-exempt"
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
}
