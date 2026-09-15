package lu

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
			i18n.FR: "TVA",
			i18n.DE: "MwSt.",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.FR: "Taxe sur la valeur ajoutée",
			i18n.DE: "Mehrwertsteuer",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Guichet.lu - VAT rates",
				},
				URL: "https://guichet.public.lu/en/entreprises/fiscalite/tva/assujettissement/taux-tva.html",
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
					i18n.FR: "Taux normal",
					i18n.DE: "Normalsatz",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(170, 3),
					},
					{
						Since:   cal.NewDate(2023, 1, 1),
						Percent: num.MakePercentage(160, 3),
					},
					{
						Since:   cal.NewDate(2015, 1, 1),
						Percent: num.MakePercentage(170, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.FR: "Taux intermédiaire",
					i18n.DE: "Mittlerer Satz",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(140, 3),
					},
					{
						Since:   cal.NewDate(2023, 1, 1),
						Percent: num.MakePercentage(130, 3),
					},
					{
						Since:   cal.NewDate(2015, 1, 1),
						Percent: num.MakePercentage(140, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.FR: "Taux réduit",
					i18n.DE: "Ermäßigter Satz",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(80, 3),
					},
					{
						Since:   cal.NewDate(2023, 1, 1),
						Percent: num.MakePercentage(70, 3),
					},
					{
						Since:   cal.NewDate(2015, 1, 1),
						Percent: num.MakePercentage(80, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.FR: "Taux super-réduit",
					i18n.DE: "Stark ermäßigter Satz",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2015, 1, 1),
						Percent: num.MakePercentage(30, 3),
					},
				},
			},
		},
	},
}
