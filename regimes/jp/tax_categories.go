package jp

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		// Consumption Tax (modelled using the shared VAT category)
		{
			Code: tax.CategoryVAT,
			Name: i18n.String{
				i18n.EN: "Consumption Tax",
				i18n.JA: "消費税",
			},
			Title: i18n.String{
				i18n.EN: "Japanese Consumption Tax",
				i18n.JA: "消費税",
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.NewString("National Tax Agency - Information about Consumption Tax"),
					URL:   "https://www.nta.go.jp/english/taxes/consumption_tax/index.htm",
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
					Description: i18n.String{
						i18n.EN: "The standard rate of Consumption Tax applies to all supplies of goods and services unless a reduced rate, zero rate, or exemption applies. The 10% headline rate combines a 7.8% national Consumption Tax and a 2.2% Local Consumption Tax.",
					},
					// Values must be ordered newest-first (descending by date).
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
						i18n.EN: here.Doc(`
							The 8% reduced rate applies to food and drink for human
							consumption excluding alcoholic drinks and "dining out" (eating
							at the establishment), and to newspapers issued twice a week or
							more under a subscription. Note that catering, and food and drink
							provided at fee-charging retirement homes, are standard-rated, and
							that mixed "linked goods" bundles qualify only under specific
							value conditions - these boundaries cannot be determined from an
							invoice line alone and are not enforced by GOBL.
						`),
					},
					// The reduced rate was introduced together with the 10% standard rate
					// on 1 October 2019; it did not exist before that date.
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
}
