package se

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	{
		// Source (English): https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers/startingandrunningaswedishbusiness/declaringtaxesbusinesses/vat/vatratesandvatexemption.4.676f4884175c97df419255d.html
		// Source (Swedish): https://www4.skatteverket.se/rattsligvagledning/394966.html
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.SE: "Moms",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.SE: "Moms",
		},
		Retained: false,
		Rates: []*tax.RateDef{
			// Source: https://www4.skatteverket.se/rattsligvagledning/edition/2025.2/429414.html
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.SE: "Normalskattesats",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(1990, 1, 1),
					},
				},
			},
			// Source: https://www4.skatteverket.se/rattsligvagledning/edition/2025.2/394978.html
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "First Reduced Rate",
					i18n.SE: "Första reducerade skattesatsen",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(120, 3),
						Since:   cal.NewDate(1996, 1, 1),
					},
				},
			},
			// Source: https://www4.skatteverket.se/rattsligvagledning/edition/2025.2/394984.html
			{
				Key: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Second Reduced Rate",
					i18n.SE: "Andra reducerade skattesatsen",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(60, 3),
						Since:   cal.NewDate(1996, 1, 1),
					},
				},
			},
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.SE: "Momsfri",
				},
				Exempt: true,
			},
		},
	},
}
