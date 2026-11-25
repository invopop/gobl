package lt

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Lithuania's VAT rates.
// The former 9% reduced rate (in force until 2025-12-31) is deliberately not backfilled:
// on 2026-01-01 its scope split three ways, with accommodation, passenger transport and
// culture moving to the 12% reduced rate, books to the 5% super-reduced rate, and
// residential heating, hot water and firewood to the standard rate, so attaching that
// history to any tier would misrepresent which supplies were taxed at which rate.
// Any document applying the former 9% (supplies before 2026, or the transitional cases
// of Art. 2 of Law XV-287) must omit the rate key and set the percent directly.
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.LT: "PVM",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.LT: "Pridėtinės vertės mokestis",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("e-TAR - Pridėtinės vertės mokesčio įstatymas Nr. IX-751, 19 straipsnis"),
				URL:   "https://www.e-tar.lt/portal/lt/legalAct/TAR.ED68997709F5/asr",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
			{
				// Source of the 21% standard rate, in force from 2009-09-01.
				Title: i18n.NewString("e-Seimas - Pridėtinės vertės mokesčio įstatymo 2, 58 ir 91 straipsnių pakeitimo įstatymas Nr. XI-386"),
				URL:   "https://e-seimas.lrs.lt/portal/legalAct/lt/TAD/TAIS.350400",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
			{
				// Source of the 12% reduced rate and the current scope of Art. 19, in force
				// from 2026-01-01.
				Title: i18n.NewString("e-TAR - Pridėtinės vertės mokesčio įstatymo Nr. IX-751 19 straipsnio pakeitimo įstatymas Nr. XV-287"),
				URL:   "https://www.e-tar.lt/portal/lt/legalAct/14bbc1a04dd411f0b070ee7f1ceefc75",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
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
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2009, 9, 1),
						Percent: num.MakePercentage(210, 3), // 21.0%
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
					i18n.EN: "Accommodation under tourism legislation, passenger transport on regular routes and the passengers' baggage, and admission to art and culture institutions and events.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2026, 1, 1),
						Percent: num.MakePercentage(120, 3), // 12.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-reduced Rate",
				},
				Description: i18n.String{
					i18n.EN: "Reimbursed medicines, medical aids and foods for special medical purposes, prescription medicines, technical aids for persons with disabilities and their repair, newspapers, magazines and periodicals, and books and non-periodical information publications.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2004, 1, 1),
						Percent: num.MakePercentage(50, 3), // 5.0%
					},
				},
			},
		},
	},
}
