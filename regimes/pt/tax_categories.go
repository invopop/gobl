package pt

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	// VAT
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.PT: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.PT: "Imposto sobre o Valor Acrescentado",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.PT: "Tipo Geral",
				},
				Values: []*tax.RateValueDef{
					{
						Ext: tax.Extensions{
							ExtKeyRegion: RegionAzores,
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(160, 3),
					},
					{
						Ext: tax.Extensions{
							ExtKeyRegion: RegionMadeira,
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(220, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(230, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.PT: "Taxa Intermédia", //nolint:misspell
				},
				Values: []*tax.RateValueDef{
					{
						Ext: tax.Extensions{
							ExtKeyRegion: RegionAzores,
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(90, 3),
					},
					{
						Ext: tax.Extensions{
							ExtKeyRegion: RegionMadeira,
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(130, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.PT: "Taxa Reduzida",
				},
				Values: []*tax.RateValueDef{
					{
						Ext: tax.Extensions{
							ExtKeyRegion: RegionAzores,
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(40, 3),
					},
					{
						Ext: tax.Extensions{
							ExtKeyRegion: RegionMadeira,
						},
						Since:   cal.NewDate(2024, 10, 1),
						Percent: num.MakePercentage(40, 3),
					},
					{
						Ext: tax.Extensions{
							ExtKeyRegion: RegionMadeira,
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(50, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(60, 3),
					},
				},
			},
			{
				// Other is a special case for rates that are not defined.
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateOther,
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.PT: "Outro",
				},
				Values: []*tax.RateValueDef{},
			},
		},
	},
}
