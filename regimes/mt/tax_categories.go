package mt

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
			i18n.MT: "VAT",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.MT: "Taxxa fuq il-Valur Miżjud",
		},
		Retained: false,
		// Standard VAT keys cover Malta's needs: "zero" models exempt-with-credit
		// (Fifth Schedule Part One), "exempt" models exempt-without-credit (Part Two),
		// plus export, intra-community and reverse-charge.
		Keys: tax.GlobalVATKeys(),
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("MTCA - VAT Rates"),
				URL:   "https://mtca.gov.mt/business-tax/vat1/vat-compliance/vat-rates/vat-rates",
			},
			{
				Title: i18n.NewString("MTCA - VAT Rates and Exemptions FAQs (Eighth and Fifth Schedules)"),
				URL:   "https://mtca.gov.mt/docs/default-source/documents/business-tax/vat/faqs/vat-rates-and-exemptions---faqs.pdf",
			},
		},
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2004, 1, 1),
						Percent: num.MakePercentage(18, 2),
					},
					{
						// VAT was reintroduced by Cap. 406 on 1 January 1999 at 15%.
						Since:   cal.NewDate(1999, 1, 1),
						Percent: num.MakePercentage(15, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate (12%)",
				},
				Description: i18n.String{
					i18n.EN: "Custody and management of securities, credit management, hiring of pleasure boats, and care of the human body (Eighth Schedule items 12-15).",
				},
				Values: []*tax.RateValueDef{
					{
						// Introduced 1 January 2024 by Legal Notice 231 of 2023.
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(12, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate (7%)",
				},
				Description: i18n.String{
					i18n.EN: "Licensed tourist accommodation and use of sporting facilities (Eighth Schedule items 1, 11).",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(7, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Reduced Rate (5%)",
				},
				Description: i18n.String{
					i18n.EN: "Electricity, confectionery, medical accessories, printed matter, minor repairs, domestic care and cultural admissions (Eighth Schedule items 2-10).",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1999, 1, 1),
						Percent: num.MakePercentage(5, 2),
					},
				},
			},
		},
	},
}
