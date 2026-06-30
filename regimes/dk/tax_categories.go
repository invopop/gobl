package dk

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// TaxCategoryExcise is the Danish excise duty (punktafgift), a tax distinct from
// VAT levied on specific goods (alcohol, tobacco, energy, …). Its rates vary per
// product, so they are supplied per line rather than predefined.
const TaxCategoryExcise cbc.Code = "EXCISE"

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.DA: "moms",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.DA: "Merværdiafgift",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.DA: "Standardsats",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(1992, 1, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Danish Tax Agency - VAT rates",
				},
				URL: "https://skat.dk/erhverv/moms",
			},
		},
	},
	{
		Code:     TaxCategoryExcise,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "Excise Duty",
			i18n.DA: "Punktafgift",
		},
		Title: i18n.String{
			i18n.EN: "Danish Excise Duty",
			i18n.DA: "Punktafgift",
		},
		// Excise rates vary per product (alcohol, tobacco, energy, …); users supply
		// the applicable percentage on the invoice rather than referencing a fixed
		// rate, as with the Spanish IPSI category.
		Rates: []*tax.RateDef{},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Danish Tax Agency - Excise duties",
				},
				URL: "https://skat.dk/erhverv/punktafgifter",
			},
		},
	},
}
