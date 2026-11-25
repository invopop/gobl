package uy

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// IVA rates and tax reform references:
// - Law 18.083 (Tax Reform, effective 2007-07-01): established the current 22% standard
//   and 10% reduced ("tasa mínima") rates by modifying Article 18 of Título 10
//   (Texto Ordenado 1996).
// - Law text: https://www.impo.com.uy/bases/leyes/18083-2006
// - Consolidated rates (Art. 18, Título 10): https://www.impo.com.uy/bases/todgi-2023/10-2024/10

var taxCategories = []*tax.CategoryDef{
	//
	// VAT (IVA)
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
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "IVA - Título 10, Texto Ordenado 2023 (Art. 18: rates)",
					i18n.ES: "IVA - Título 10, Texto Ordenado 2023 (Art. 18: tasas)",
				},
				URL: "https://www.impo.com.uy/bases/todgi-2023/10-2024/10",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Tasa Básica",
				},
				Description: i18n.String{
					i18n.EN: "Applies to the majority of goods and services.",
					i18n.ES: "Se aplica a la mayoría de bienes y servicios.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2007, 7, 1),
						Percent: num.MakePercentage(220, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Tasa Mínima",
				},
				Description: i18n.String{
					i18n.EN: "Applies to basic necessities including food, medicine, hotel services, passenger transport, and health services.",
					i18n.ES: "Se aplica a artículos de primera necesidad como alimentos, medicamentos, servicios de hotelería, transporte de pasajeros y servicios de salud.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2007, 7, 1),
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
		},
	},
}
