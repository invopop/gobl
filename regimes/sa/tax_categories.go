package sa

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
				i18n.AR: "ضريبة القيمة المضافة",
			},
			Title: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.AR: "ضريبة القيمة المضافة",
			},
			Retained: false,
			Keys:     tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard Rate",
						i18n.AR: "النسبة الأساسية",
					},
					Description: i18n.String{
						i18n.EN: "Applies to most goods and services unless specified otherwise.",
						i18n.AR: "تنطبق على معظم السلع والخدمات ما لم ينص على خلاف ذلك.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2020, 7, 1),
							Percent: num.MakePercentage(150, 3),
						},
						{
							Since:   cal.NewDate(2018, 1, 1),
							Percent: num.MakePercentage(50, 3),
						},
					},
				},
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "ZATCA - VAT Regulations",
						i18n.AR: "هيئة الزكاة والضريبة والجمارك",
					},
					URL: "https://zatca.gov.sa/en/RulesRegulations/VAT/Pages/default.aspx",
				},
			},
		},
	}
}
