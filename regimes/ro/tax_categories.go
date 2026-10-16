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
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ANAF - Codul fiscal (Legea 227/2015), cotele de TVA"),
				URL:   "https://static.anaf.ro/static/10/Anaf/legislatie/Cod_fiscal_norme_2023.htm",
			},
			{
				Title: i18n.NewString("Legea nr. 141/2025 privind unele măsuri fiscal-bugetare (TVA 21% / 11%)"),
				URL:   "https://static.anaf.ro/static/10/Anaf/legislatie/L_141_2025.pdf",
			},
			{
				Title: i18n.NewString("ANAF - Modificări aduse Codului fiscal prin Legea nr. 141/2025"),
				URL:   "https://static.anaf.ro/static/10/Brasov/Brasov/tva_2025.pdf",
			},
			{
				Title: i18n.NewString("ANAF - Cote de TVA pentru livrarea de locuințe, începând cu data de 01 august 2025"),
				URL:   "https://static.anaf.ro/static/3/Cluj/20250811163039_cj_tva_locuinte_05aug2025.pdf",
			},
			{
				Title: i18n.NewString("Legea nr. 227/2015 privind Codul fiscal, art. 291 (TVA 20% / 19%, 9%, 5%)"),
				URL:   "https://legislatie.just.ro/Public/DetaliiDocument/171282",
			},
			{
				Title: i18n.NewString("OUG nr. 58/2010, art. 44 (TVA 24% din 1 iulie 2010)"),
				URL:   "https://legislatie.just.ro/Public/DetaliiDocumentAfis/119870",
			},
			{
				Title: i18n.NewString("OUG nr. 200/2008 (TVA 5% pentru locuințe sociale din 15 decembrie 2008)"),
				URL:   "https://legislatie.just.ro/Public/DetaliiDocumentAfis/99787",
			},
			{
				Title: i18n.NewString("Legea nr. 571/2003 privind Codul fiscal, art. 140 (TVA 19%, 9% din 1 ianuarie 2004)"),
				URL:   "https://legislatie.just.ro/Public/DetaliiDocument/48725",
			},
			{
				Title: i18n.NewString("OUG nr. 215/1999 (TVA 19% din 1 ianuarie 2000)"),
				URL:   "https://legislatie.just.ro/Public/DetaliiDocumentAfis/20397",
			},
			{
				Title: i18n.NewString("OG nr. 3/1992 privind taxa pe valoarea adăugată (TVA 18% din 1 iulie 1993, 9% din 1 ianuarie 1995)"),
				URL:   "https://legislatie.just.ro/Public/DetaliiDocument/2092",
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
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(210, 3), // 21%
					},
					{
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(190, 3), // 19%
					},
					{
						Since:   cal.NewDate(2016, 1, 1),
						Percent: num.MakePercentage(200, 3), // 20%
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(240, 3), // 24%
					},
					{
						Since:   cal.NewDate(2000, 1, 1),
						Percent: num.MakePercentage(190, 3), // 19%
					},
					{
						Since:   cal.NewDate(1998, 2, 1),
						Percent: num.MakePercentage(220, 3), // 22%
					},
					{
						Since:   cal.NewDate(1993, 7, 1),
						Percent: num.MakePercentage(180, 3), // 18%
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
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(110, 3), // 11%
					},
					{
						Since:   cal.NewDate(2004, 1, 1),
						Percent: num.MakePercentage(90, 3), // 9%
					},
					{
						Since:   cal.NewDate(1998, 2, 1),
						Percent: num.MakePercentage(110, 3), // 11%
					},
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(90, 3), // 9%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.RO: "Cota Redusă de 5%",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(110, 3), // 11%
					},
					{
						Since:   cal.NewDate(2008, 12, 15),
						Percent: num.MakePercentage(50, 3), // 5%
					},
				},
			},
		},
	},
}
