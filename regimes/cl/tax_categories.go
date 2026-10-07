package cl

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		{
			Code: tax.CategoryVAT,
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.ES: "IVA",
			},
			Title: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.ES: "Impuesto al Valor Agregado",
			},
			Retained: false,
			Keys:     tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard Rate",
						i18n.ES: "Tasa General",
					},
					Description: i18n.String{
						i18n.EN: "Standard rate applied to most goods and services unless they are exempt.",
						i18n.ES: "Tasa estándar aplicada a la mayoría de bienes y servicios a menos que estén exentos.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2003, 10, 1),
							Percent: num.MakePercentage(190, 3),
						},
					},
				},
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "Decree Law 825 - Law on Sales and Services Tax (Art. 14, standard rate)",
						i18n.ES: "Decreto Ley 825 - Ley sobre Impuesto a las Ventas y Servicios (Art. 14, tasa general)",
					},
					URL: "https://www.sii.cl/normativa_legislacion/dl825.pdf",
				},
			},
		},
	}
}
