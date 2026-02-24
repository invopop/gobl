package ad

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// RateIncreased is the key for the increased rate in Andorra.
const RateIncreased cbc.Key = "increased"

func taxCategories() []*tax.CategoryDef {
	//
	// IGI (VAT)
	//
	return []*tax.CategoryDef{
		{
			Code: tax.CategoryVAT,
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.CA: "IGI",
				i18n.ES: "IGI",
			},
			Title: i18n.String{
				i18n.EN: "General Indirect Tax",
				i18n.CA: "Impost General Indirecte",
				i18n.ES: "Impuesto General Indirecto",
			},
			Description: &i18n.String{
				i18n.EN: "The General Indirect Tax (IGI) is the main indirect tax levied on consumption in Andorra.",
				i18n.CA: "L'Impost General Indirecte (IGI) és el principal impost indirecte que grava el consum a Andorra.",
				i18n.ES: "El Impuesto General Indirecto (IGI) es el principal impuesto indirecto que grava el consumo en Andorra.",
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "Departament de Tributs i de Fronteres - Andorra",
						i18n.CA: "Departament de Tributs i de Fronteres - Andorra",
						i18n.ES: "Departamento de Tributos y Fronteras - Andorra",
					},
					URL: "https://www.e-tramits.ad/tramits/ca/impostos/igi",
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
						i18n.CA: "Tipus general",
						i18n.ES: "Tipo general",
					},
					Description: i18n.String{
						i18n.EN: "General IGI rate applied to most goods and services (4.5%).",
						i18n.CA: "Tipus general de l'IGI aplicat a la majoria de béns i serveis (4,5%).",
						i18n.ES: "Tipo general del IGI aplicado a la mayoría de bienes y servicios (4,5%).",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(45, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced Rate",
						i18n.CA: "Tipus reduït",
						i18n.ES: "Tipo reducido",
					},
					Description: i18n.String{
						i18n.EN: "Reduced IGI rate (1%).",
						i18n.CA: "Tipus reduït de l'IGI (1%).",
						i18n.ES: "Tipo reducido del IGI (1%).",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(10, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyZero},
					Rate: tax.RateSuperReduced,
					Name: i18n.String{
						i18n.EN: "Super-Reduced Rate",
						i18n.CA: "Tipus superreduït",
						i18n.ES: "Tipo superreducido",
					},
					Description: i18n.String{
						i18n.EN: "Super-reduced IGI rate (0%).",
						i18n.CA: "Tipus superreduït de l'IGI (0%).",
						i18n.ES: "Tipo superreducido del IGI (0%).",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateSpecial,
					Name: i18n.String{
						i18n.EN: "Special Rate",
						i18n.CA: "Tipus especial",
						i18n.ES: "Tipo especial",
					},
					Description: i18n.String{
						i18n.EN: "Special IGI rate (2.5%).",
						i18n.CA: "Tipus especial de l'IGI (2,5%).",
						i18n.ES: "Tipo especial del IGI (2,5%).",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(25, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: RateIncreased,
					Name: i18n.String{
						i18n.EN: "Increased Rate",
						i18n.CA: "Tipus incrementat",
						i18n.ES: "Tipo incrementado",
					},
					Description: i18n.String{
						i18n.EN: "Increased IGI rate (9.5%).",
						i18n.CA: "Tipus incrementat de l'IGI (9,5%).",
						i18n.ES: "Tipo incrementado del IGI (9,5%).",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(95, 3),
						},
					},
				},
			},
		},
	}
}
