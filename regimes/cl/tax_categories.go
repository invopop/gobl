package cl

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// IVA (Impuesto al Valor Agregado)
	//
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
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "VAT law - Article 14",
				},
				URL: "https://www.bcn.cl/leychile/navegar?idNorma=6369",
			},
		},
		Keys: tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.ES: "Tasa General",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2003, 10, 1),
						Percent: num.MakePercentage(19, 2),
					},
				},
			},
		},
	},
}
