package ad

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var categories = []*tax.CategoryDef{
	{
		Code:     tax.CategoryVAT,
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Name: i18n.String{
			i18n.EN: "IGI",
			i18n.CA: "IGI",
		},
		Title: i18n.String{
			i18n.EN: "General Indirect Tax",
			i18n.CA: "Impost General Indirecte",
		},
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{cbc.KeyEmpty, tax.KeyStandard},
				Rate: "increased",
				Name: i18n.String{
					i18n.EN: "Increased Rate",
					i18n.CA: "Tipus Incrementat",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(95, 3)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.CA: "Tipus General",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(45, 3)},
				},
			},
			{
				Keys: []cbc.Key{cbc.KeyEmpty, tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.CA: "Tipus Intermedi",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(25, 3)},
				},
			},
			{
				Keys: []cbc.Key{cbc.KeyEmpty, tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.CA: "Tipus Reduït",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(1, 2)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.CA: "Tipus Zero",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(0, 0)},
				},
			},
		},
	},
}
