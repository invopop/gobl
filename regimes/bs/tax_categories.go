package bs

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
			Code: tax.CategoryVAT,
			Name: i18n.String{
				i18n.EN: "VAT",
			},
			Title: i18n.String{
				i18n.EN: "Value Added Tax",
			},
			Retained: false,
			Keys:     tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard Rate",
					},
					Description: i18n.String{
						i18n.EN: "Applies to most taxable goods and services.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2022, 1, 1),
							Percent: num.MakePercentage(10, 2), // 10%
						},
						{
							Since:   cal.NewDate(2018, 7, 1),
							Percent: num.MakePercentage(12, 2), // 12%
						},
						{
							Since:   cal.NewDate(2015, 1, 1),
							Percent: num.MakePercentage(75, 3), // 7.5%
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced Rate",
					},
					Description: i18n.String{
						i18n.EN: "Applies to some qualifying essential goods.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2025, 4, 1),
							Percent: num.MakePercentage(5, 2), // 5%
						},
					},
				},
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.NewString("PRIME MINISTER DAVIS’S CONTRIBUTION ON THE COMMONWEALTH OF THE BAHAMAS 2025/2026 MID-YEAR BUDGET STATEMENT"),
					URL:   "https://opm.gov.bs/prime-minister-daivs-contribution-2025-2026-mid-year-budget-statement/",
				},
				{
					Title: i18n.NewString("ABOUT - VAT, Department of Inland Revenue"),
					URL:   "https://inlandrevenue.finance.gov.bs/value-added-tax/about-vat/",
				},
				{
					Title: i18n.NewString("VAT Guidance on Transitional Arrangements"),
					URL:   "https://inlandrevenue.finance.gov.bs/wp-content/uploads/2015/02/2015-11-01-Bahamas-VAT-Guidance_Transitional-Arrangements-Final-.pdf",
				},
				{
					Title: i18n.NewString("OFFICIAL GAZETTE - THE BAHAMAS (published on 1st July 2018)"),
					URL:   "https://inlandrevenue.finance.gov.bs/wp-content/uploads/2022/02/ValueAddedTaxAmendmentAct2018.pdf",
				},
				{
					Title: i18n.NewString("VAT GUIDANCE ON 2021 AMENDMENTS"),
					URL:   "https://inlandrevenue.finance.gov.bs/wp-content/uploads/2021/12/VAT-Guidance-2021-Amendments.pdf",
				},
				{
					Title: i18n.NewString("VALUE ADDED TAX (AMENDMENT) BILL, 2025"),
					URL:   "https://laws.bahamas.gov.bs/cms/images/LEGISLATION/BILLS/2025/2025-0009/2025-0009.pdf",
				},
			},
		},
	}
}
