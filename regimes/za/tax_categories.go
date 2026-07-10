package za

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// taxCategories defines the VAT category and rates for South Africa.
//
// Exempt, reverse-charge, export, intra-community, and outside-scope
// supplies are already covered by the global VAT keys (tax.GlobalVATKeys)
// and require no local rate: they carry no percentage at all, as opposed to
// zero-rated supplies which do carry an explicit (if 0%) rate. This
// distinction is legally significant under the VAT Act — zero-rated
// vendors may still recover input VAT, exempt vendors may not — and is
// preserved here rather than collapsed.
var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
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
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("SARS - Value-Added Tax"),
				URL:   "https://www.sars.gov.za/types-of-tax/value-added-tax/",
			},
			{
				Title: i18n.NewString("Zero-rated and exempt supplies (Parliamentary Monitoring Group)"),
				URL:   "https://static.pmg.org.za/docs/Zero-rated%20and%20exempt%20supplies.pdf",
			},
		},
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
				},
				Values: []*tax.RateValueDef{
					{
						// Increased 1 April 2018 under the Rates and Monetary
						// Amounts and Amendment of Revenue Laws Act 21 of 2018.
						// A further increase to 15.5%/16% was announced for
						// 2025/2026 but withdrawn before taking effect.
						Since:   cal.NewDate(2018, 4, 1),
						Percent: num.MakePercentage(15, 2),
					},
					{
						// Increased 7 April 1993.
						Since:   cal.NewDate(1993, 4, 7),
						Percent: num.MakePercentage(14, 2),
					},
					{
						// VAT introduced 30 September 1991 at 10%.
						Since:   cal.NewDate(1991, 9, 30),
						Percent: num.MakePercentage(10, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to exports, international transport, basic foodstuffs (e.g. brown bread, maize meal, rice, vegetables, milk, eggs), and fuel (petrol, diesel, illuminating paraffin). Input VAT remains recoverable, unlike exempt supplies.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1991, 9, 30),
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
		},
	},
}
