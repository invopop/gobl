package is

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
				i18n.IS: "VSK",
			},
			Title: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.IS: "Virðisaukaskattur",
			},
			Retained: false,
			Keys:     tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard Rate",
						i18n.IS: "Almennt skatthlutfall",
					},
					Description: i18n.String{
						i18n.EN: "Applies to most goods and services unless a reduced rate or zero-rating applies.",
						i18n.IS: "Gildir um flestar vörur og þjónustu nema að lægra þrep eða núll-hlutfall eigi við.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2015, 1, 1),
							Percent: num.MakePercentage(240, 3),
						},
						{
							Since:   cal.NewDate(2010, 1, 1),
							Percent: num.MakePercentage(255, 3),
						},
						{
							Since:   cal.NewDate(1990, 1, 1),
							Percent: num.MakePercentage(220, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced Rate",
						i18n.IS: "Lægra skatthlutfall",
					},
					Description: i18n.String{
						i18n.EN: "Accommodation, books and periodicals, food, radio/TV licences and certain cultural services.",
						i18n.IS: "Gisting, bækur og tímarit, matvæli, afnotagjöld útvarps og tiltekin menningarþjónusta.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2015, 1, 1),
							Percent: num.MakePercentage(110, 3),
						},
						{
							Since:   cal.NewDate(2007, 3, 1),
							Percent: num.MakePercentage(70, 3),
						},
					},
				},
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "Skatturinn - Value Added Tax",
						i18n.IS: "Skatturinn - Virðisaukaskattur",
					},
					URL: "https://www.skatturinn.is/english/companies/value-added-tax/",
				},
				{
					Title: i18n.String{
						i18n.IS: "Skatturinn - Skattskylda og skattprósentur",
					},
					URL: "https://www.skatturinn.is/atvinnurekstur/virdisaukaskattur/skattskylda-og-skattprosentur/",
				},
				{
					Title: i18n.NewString("OECD - Consumption Tax Trends: Iceland"),
					URL:   "https://www.oecd.org/content/dam/oecd/en/topics/policy-sub-issues/consumption-tax-trends/consumption-tax-trends-iceland.pdf",
				},
			},
		},
	}
}
