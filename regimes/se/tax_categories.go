package se

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
			i18n.EN: "VAT",
			i18n.SE: "Moms",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.SE: "Moms",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.SE: "Normalskattesats",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(1990, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
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
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
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
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Swedish Tax Agency - VAT rates and VAT exemption",
				},
				URL: "https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers/startingandrunningaswedishbusiness/declaringtaxesbusinesses/vat/vatratesandvatexemption.4.676f4884175c97df419255d.html",
				At:  cal.NewDateTime(2025, 4, 25, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.SE: "Skatteverket - Momssatser och undantag från moms",
				},
				URL: "https://www4.skatteverket.se/rattsligvagledning/394966.html",
				At:  cal.NewDateTime(2025, 4, 25, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.SE: "Normalskattesats",
					i18n.EN: "Standard rate",
				},
				URL: "https://www4.skatteverket.se/rattsligvagledning/edition/2025.4/429414.html",
				At:  cal.NewDateTime(2025, 4, 25, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.SE: "Skattesats 12 procent",
					i18n.EN: "Tax rate 12 percent",
				},
				URL: "https://www4.skatteverket.se/rattsligvagledning/edition/2025.4/394978.html",
				At:  cal.NewDateTime(2025, 4, 25, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.SE: "Skattesats 6 procent",
					i18n.EN: "Tax rate 6 percent",
				},
				URL: "https://www4.skatteverket.se/rattsligvagledning/edition/2025.4/394984.html",
				At:  cal.NewDateTime(2025, 4, 25, 0, 0, 0),
			},
		},
	},
}
