package nz

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
			Code: tax.CategoryGST,
			Name: i18n.String{
				i18n.EN: "GST",
				i18n.MI: "GST",
			},
			Title: i18n.String{
				i18n.EN: "Goods and Services Tax",
				i18n.MI: "Tāke mō ngā Rawa me ngā Ratonga",
			},
			Retained: false,
			Keys:     tax.GlobalGSTKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard Rate",
						i18n.MI: "Pāpātanga Paerewa",
					},
					Description: i18n.String{
						i18n.EN: "Applies to most goods and services unless they are zero-rated or exempt.",
						i18n.MI: "Ka pā ki te nuinga o ngā rawa me ngā ratonga ki te kore rātou e tāke-kore, e tāke-whakawātea rānei.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2010, 10, 1),
							Percent: num.MakePercentage(150, 3),
						},
					},
				},
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.NewString("Inland Revenue - GST rate"),
					URL:   "https://www.ird.govt.nz/gst",
				},
			},
		},
	}
}
