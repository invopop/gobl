package fi

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// sources of truth:
// https://vatapp.net/vat-rates#vat-rates-fi
// https://www.vero.fi/en/businesses-and-corporations/taxes-and-charges/vat/rates-of-vat/
// https://www.vero.fi/en/businesses-and-corporations/taxes-and-charges/vat/rates-of-vat/new-vat-rate-from-1-september-2024--instructions-for-vat-reporting/
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.FI: "ALV",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.FI: "Arvonlisävero",
		},
		Description: &i18n.String{
			i18n.EN: here.Doc(`
                Known in Finnish as "Arvonlisävero" (ALV), this is a consumption tax applied to
                the purchase of goods and services. As a member of the European Union, Finland
                follows the EU VAT Directive, with a standard rate and reduced rates tailored to
                specific goods and services.
            `),
		},
		Rates: []*tax.RateDef{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.FI: "Yleinen verokanta",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 9, 1),
						Percent: num.MakePercentage(255, 3),
					},
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(240, 3),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(230, 3),
					},
					{
						Since:   cal.NewDate(1994, 6, 1),
						Percent: num.MakePercentage(220, 3),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate 1",
					i18n.FI: "Alennettu verokanta 1",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2023, 1, 1),
						Percent: num.MakePercentage(140, 3),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(130, 3),
					},
					{
						Since:   cal.NewDate(2009, 10, 1),
						Percent: num.MakePercentage(120, 3),
					},
					{
						Since:   cal.NewDate(1998, 1, 1),
						Percent: num.MakePercentage(170, 3),
					},
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
					{
						Since:   cal.NewDate(1994, 6, 1),
						Percent: num.MakePercentage(60, 3),
					},
				},
			},
			{
				Key: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate 2",
					i18n.FI: "Alennettu verokanta 2",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2023, 1, 1),
						Percent: num.MakePercentage(100, 3),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(90, 3),
					},
					{
						Since:   cal.NewDate(1998, 1, 1),
						Percent: num.MakePercentage(80, 3),
					},
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(60, 3),
					},
					{
						Since:   cal.NewDate(1994, 6, 1),
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.FI: "Nollaverokanta",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
		},
	},
}
