package cz

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Czech VAT rates, as set in § 47 of the VAT act and changed by each
// amendment cited below. The history stops at 2010.
//
// A second reduced rate of 10% existed from 2015 to 2023 and is deliberately
// left out: rate values only have a start date, so a super-reduced rate would
// still resolve to 10% on documents issued after it was abolished.
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.CS: "DPH",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.CS: "Daň z přidané hodnoty",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Zákon č. 235/2004 Sb., o dani z přidané hodnoty, § 47 (Sazby daně)"),
				URL:   "https://www.zakonyprolidi.cz/cs/2004-235#p47",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Zákon č. 349/2023 Sb., čl. LX (snížená sazba 12 % od 1. 1. 2024)"),
				URL:   "https://www.zakonyprolidi.cz/cs/2023-349",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Zákon č. 500/2012 Sb., čl. III (sazby 21 % a 15 % od 1. 1. 2013)"),
				URL:   "https://www.zakonyprolidi.cz/cs/2012-500",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Zákon č. 370/2011 Sb., čl. I (snížená sazba 14 % od 1. 1. 2012)"),
				URL:   "https://www.zakonyprolidi.cz/cs/2011-370",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Zákon č. 362/2009 Sb., čl. VII (sazby 20 % a 10 % od 1. 1. 2010)"),
				URL:   "https://www.zakonyprolidi.cz/cs/2009-362",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.CS: "Základní sazba",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(210, 3), // 21.0%
					},
					{
						Since:   cal.NewDate(2010, 1, 1),
						Percent: num.MakePercentage(200, 3), // 20.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.CS: "Snížená sazba",
				},
				Description: i18n.String{
					i18n.EN: "Foodstuffs, medicines, accommodation, water and heat.",
					i18n.CS: "Potraviny, léky, ubytování, voda a teplo.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(120, 3), // 12.0%
					},
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(150, 3), // 15.0%
					},
					{
						Since:   cal.NewDate(2012, 1, 1),
						Percent: num.MakePercentage(140, 3), // 14.0%
					},
					{
						Since:   cal.NewDate(2010, 1, 1),
						Percent: num.MakePercentage(100, 3), // 10.0%
					},
				},
			},
		},
	},
}
