package ro

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.RO: "TVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.RO: "Taxa pe Valoarea Adăugată",
		},

		// ---- Trustworthy sources about tax systems in Romania (ANAF) ----
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ANAF - Cod fiscal 227/2015 Cotele de TVA (VAT rates)"),
				URL:   "https://static.anaf.ro/static/10/Anaf/legislatie/Cod_fiscal_norme_2023.htm",
			},
			{
				Title: i18n.NewString("MODIFICĂRI ADUSE LEGII NR. 227/2015 PRIVIND CODUL FISCAL PRIN LEGEA NR. 141/2025*. VAT 19% -> 21%"),
				URL:   "https://static.anaf.ro/static/10/Brasov/Brasov/tva_2025.pdf",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.RO: "Cota Standard",
				},

				// ---- General Rate History ----
				Values: []*tax.RateValueDef{
					{
						// 1 Aug 2025
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(210, 3), // 21%
					},
					{
						// 1 January 2017
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(190, 3), // 19%
					},
					{
						// 1 January 2016
						Since:   cal.NewDate(2016, 1, 1),
						Percent: num.MakePercentage(200, 3), // 20%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.RO: "Cota Redusă",
				},

				// ---- Reduced Rate History ----
				Values: []*tax.RateValueDef{
					{
						// 1 Aug 2025
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(110, 3), // 11%, unifies former 5%/9% rates
					},
				},
			},
		},
	},
}
