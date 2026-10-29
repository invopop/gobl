package sk

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Slovakia's VAT rates.
// Only the standard rate carries a pre-2025 value (20%, in force since 2011). The
// reduced and super-reduced rates are deliberately not backfilled: before 2025 there
// was a single 10% reduced rate on príloha č. 7 goods (medicines, books), whose scope
// sits closer to today's 5% super-reduced than to the new 19% tier, so attaching that
// history to either rate would misrepresent which goods were taxed at which rate.
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.SK: "DPH",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.SK: "Daň z pridanej hodnoty",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Slov-Lex - Zákon č. 222/2004 Z. z. § 27 (Sadzba dane)"),
				URL:   "https://www.slov-lex.sk/ezbierky/pravne-predpisy/SK/ZZ/2004/222/",
				At:    cal.NewDateTime(2026, 8, 6, 0, 0, 0),
			},
			{
				// Source of the 20% standard rate, introduced from 2011-01-01 by § 85j and
				// carried in § 27 itself from 2015-01-01.
				Title: i18n.NewString("Slov-Lex - Zákon č. 490/2010 Z. z. (novela zákona č. 222/2004 Z. z.)"),
				URL:   "https://www.slov-lex.sk/ezbierky/pravne-predpisy/SK/ZZ/2010/490/",
				At:    cal.NewDateTime(2026, 8, 6, 0, 0, 0),
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
					i18n.SK: "Základná sadzba",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 1, 1),
						Percent: num.MakePercentage(230, 3), // 23.0%
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(200, 3), // 20.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.SK: "Znížená sadzba",
				},
				Description: i18n.String{
					i18n.EN: "Non-basic foodstuffs, electricity and certain catering services.",
					i18n.SK: "Potraviny okrem základných, elektrina a niektoré reštauračné služby.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 1, 1),
						Percent: num.MakePercentage(190, 3), // 19.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-reduced Rate",
					i18n.SK: "Druhá znížená sadzba",
				},
				Description: i18n.String{
					i18n.EN: "Basic foodstuffs, medicines, medical devices, books, and accommodation.",
					i18n.SK: "Základné potraviny, lieky, zdravotnícke pomôcky, knihy a ubytovanie.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 1, 1),
						Percent: num.MakePercentage(50, 3), // 5.0%
					},
				},
			},
		},
	},
}
