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
	// VAT
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
					i18n.EN: "SII - What is the VAT rate?",
					i18n.ES: "SII - ¿Cuál es la tasa del Impuesto al Valor Agregado (IVA)?",
				},
				URL: "https://www.sii.cl/preguntas_frecuentes/impuestos_mensuales/001_130_0572.htm",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
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
						// Ley N° 19.888, published 2003-08-13, effective from
						// 2003-10-01. Chile has a single flat VAT rate with no
						// reduced or intermediate bands.
						Since:   cal.NewDate(2003, 10, 1),
						Percent: num.MakePercentage(190, 3),
					},
				},
			},
		},
	},
}
