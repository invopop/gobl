package jp

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "CT",
			i18n.JA: "消費税", // Shōhi zei
		},
		Title: i18n.String{
			i18n.EN: "Consumption Tax",
			i18n.JA: "消費税", // Shōhi zei
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "National Tax Agency - Consumption Tax",
					i18n.JA: "国税庁 - 消費税", // Kokuzei chō - Shōhi zei
				},
				URL: "https://www.nta.go.jp/english/taxes/consumption_tax/01.htm",
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
					i18n.JA: "標準税率", // Hyōjun zeiritsu
				},
				Description: i18n.String{
					i18n.EN: "Applies to most goods and services unless specified otherwise.",
					// Toku ni sadame ga nai kagiri, hotondo no shōhin sābisu ni tekiyō saremasu.
					i18n.JA: "特に定めがない限り、ほとんどの商品・サービスに適用されます。",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 10, 1),
						Percent: num.MakePercentage(10, 2),
					},
					{
						Since:   cal.NewDate(2014, 4, 1),
						Percent: num.MakePercentage(8, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.JA: "軽減税率", // Keigen zeiritsu
				},
				Description: i18n.String{
					i18n.EN: "Applies to food and non-alcoholic beverages (excluding dining out) and newspaper subscriptions (2+ per week).",
					// Inshoku ryōhin (gaishoku wo nozoku) oyobi teiki kōdoku no shinbun (shū ni kai ijō hakkō) ni tekiyō saremasu.
					i18n.JA: "飲食料品（外食を除く）および定期購読の新聞（週2回以上発行）に適用されます。",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 10, 1),
						Percent: num.MakePercentage(8, 2),
					},
				},
			},
		},
	},
}
