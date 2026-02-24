package jp

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// Consumption Tax (消費税)
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "CT",
			i18n.JA: "消費税",
		},
		Title: i18n.String{
			i18n.EN: "Consumption Tax",
			i18n.JA: "消費税",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Consumption Tax Trends - Japan",
				},
				URL: "https://www.nta.go.jp/english/taxes/consumption_tax/01.htm#c01",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard rate",
					i18n.JA: "標準税率",
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
					{
						Since:   cal.NewDate(1997, 4, 1),
						Percent: num.MakePercentage(5, 2),
					},
					{
						Since:   cal.NewDate(1989, 4, 1),
						Percent: num.MakePercentage(3, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced rate",
					i18n.JA: "軽減税率",
				},
				Description: i18n.String{
					i18n.EN: "Applies to food and beverages (excluding dining out and alcohol) and newspaper subscriptions.",
					i18n.JA: "飲食料品（外食・酒類を除く）及び新聞の定期購読に適用。",
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

	//
	// Withholding Income Tax (源泉徴収)
	//
	{
		Code:     TaxCategoryWHT,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "WHT",
			i18n.JA: "源泉徴収",
		},
		Title: i18n.String{
			i18n.EN: "Withholding Income Tax",
			i18n.JA: "源泉徴収所得税",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Japan Withholding Taxes - PwC",
				},
				URL: "https://taxsummaries.pwc.com/japan/corporate/withholding-taxes",
			},
		},
		Rates: []*tax.RateDef{
			{
				Rate: TaxRatePro,
				Name: i18n.String{
					i18n.EN: "Professional Rate",
					i18n.JA: "専門家報酬",
				},
				Description: i18n.String{
					i18n.EN: "Applies to professional service fees of ¥1,000,000 or less per payment. Includes 2.1% reconstruction surtax (2013-2037).",
					i18n.JA: "1回の支払金額が100万円以下の専門家報酬に適用。復興特別所得税2.1%を含む（2013年〜2037年）。",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(1021, 4),
					},
					{
						Since:   cal.NewDate(1989, 4, 1),
						Percent: num.MakePercentage(10, 2),
					},
				},
			},
			{
				Rate: TaxRateProOver,
				Name: i18n.String{
					i18n.EN: "Professional Rate (over ¥1M)",
					i18n.JA: "専門家報酬（100万円超）",
				},
				Description: i18n.String{
					i18n.EN: "Applies to professional service fees exceeding ¥1,000,000 per payment. Includes 2.1% reconstruction surtax (2013-2037).",
					i18n.JA: "1回の支払金額が100万円を超える専門家報酬に適用。復興特別所得税2.1%を含む（2013年〜2037年）。",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(2042, 4),
					},
					{
						Since:   cal.NewDate(1989, 4, 1),
						Percent: num.MakePercentage(20, 2),
					},
				},
			},
		},
	},
}
