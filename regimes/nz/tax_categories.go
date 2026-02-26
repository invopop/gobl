package nz

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		{
			Code: tax.CategoryGST,
			Name: i18n.String{
				i18n.EN: "GST",
			},
			Title: i18n.String{
				i18n.EN: "Goods and Services Tax",
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "Charging GST",
					},
					URL: "https://www.ird.govt.nz/gst/charging-gst",
				},
			},
			Retained: false,
			Keys:     tax.GlobalGSTKeys(),
			Rates: []*tax.RateDef{
				{
					Rate: tax.RateZero,
					Name: i18n.String{
						i18n.EN: "Zero rate",
					},
					Description: i18n.String{
						i18n.EN: "Some supplies can be zero-rated or exempt under GST rules.",
					},
					Values: []*tax.RateValueDef{
						{
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "General rate",
					},
					Description: i18n.String{
						i18n.EN: "Most taxable supplies are charged at 15% GST.",
					},
					Values: []*tax.RateValueDef{
						{
							Percent: num.MakePercentage(15, 2),
						},
					},
				},
			},
		},
	}
}
