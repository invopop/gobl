package ar

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// RateIncreased is the key for the increased rate in Argentina.
const RateIncreased = "increased"

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		//
		// IVA (Impuesto al Valor Agregado)
		//
		{
			Code:     tax.CategoryVAT,
			Retained: false,
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "VAT law - Article 28",
					},
					URL: "https://biblioteca.arca.gob.ar/cuadroslegislativos/getAdjunto.aspx?i=7833",
				},
			},
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.ES: "IVA",
			},
			Title: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.ES: "Impuesto al Valor Agregado",
			},
			Description: &i18n.String{
				i18n.EN: here.Doc(`
					Known in Spanish as "Impuesto al Valor Agregado" (IVA), Argentina's VAT is a
					consumption tax applied to the sale of goods and services. The standard rate is 21%,
					with reduced rates of 10.5% for certain essential goods and services, and 0% for
					exports and specific categories.
				`),
			},
			Keys: tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: RateIncreased,
					Name: i18n.String{
						i18n.EN: "Increased Rate",
						i18n.ES: "Alícuota Incrementada",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(1995, 4, 1),
							Percent: num.MakePercentage(270, 3),
						},
					},
					Description: i18n.String{
						i18n.EN: here.Doc(`
							Applies to the sale of gas, electricity and water regulated by meter and other services stated in the points 4, 5 and 6 of the inciso e) of the article 3.
						`),
						i18n.ES: here.Doc(`
							Se aplica para las ventas de gas, energía eléctrica y aguas reguladas por medidor y demás prestaciones comprendidas en los puntos 4, 5 y 6, del inciso e) del artículo 3°.
						`),
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "General Rate",
						i18n.ES: "Alícuota General",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(1995, 4, 1),
							Percent: num.MakePercentage(210, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced Rate",
						i18n.ES: "Alícuota Reducida",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(1995, 4, 1),
							Percent: num.MakePercentage(105, 3),
						},
					},
					Description: i18n.String{
						i18n.EN: here.Doc(`
							Applies to the sale of construction work, medicine, public transportation, livestock and agricultural products for food, etc.
						`),
						i18n.ES: here.Doc(`
							Se aplica para las ventas de trabajo de construcción, medicina, transporte público, ganado y productos agrícolas para alimentos, etc.
						`),
					},
				},
			},
		},
	}
}
