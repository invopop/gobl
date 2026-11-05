package hr

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Croatia's VAT rates, dated from 1 July 2013, when the current VAT Act
// (NN 73/2013) came into force. Rates under the previous VAT Act (NN 47/1995)
// are not backfilled. Zero-rated supplies use the global zero key.
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.HR: "PDV",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.HR: "Porez na dodanu vrijednost",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Narodne novine 73/2013 - Zakon o porezu na dodanu vrijednost, članak 38"),
				URL:   "https://narodne-novine.nn.hr/clanci/sluzbeni/2013_06_73_1451.html",
				At:    cal.NewDateTime(2026, 10, 19, 0, 0, 0),
			},
			{
				// Source of the 13% reduced rate, replacing 10% from 2014-01-01 (articles 2 and 11).
				Title: i18n.NewString("Narodne novine 148/2013 - Zakon o izmjenama i dopunama Zakona o porezu na dodanu vrijednost"),
				URL:   "https://narodne-novine.nn.hr/clanci/sluzbeni/2013_12_148_3138.html",
				At:    cal.NewDateTime(2026, 10, 19, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.EN: "Croatian Tax Administration - VAT",
					i18n.HR: "Porezna uprava - Porez na dodanu vrijednost",
				},
				URL: "https://porezna-uprava.gov.hr/en/vat/7362",
				At:  cal.NewDateTime(2026, 10, 19, 0, 0, 0),
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
					i18n.HR: "Opća stopa",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 7, 1),
						Percent: num.MakePercentage(250, 3), // 25.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.HR: "Snižena stopa",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2014, 1, 1),
						Percent: num.MakePercentage(130, 3), // 13.0%
					},
					{
						Since:   cal.NewDate(2013, 7, 1),
						Percent: num.MakePercentage(100, 3), // 10.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super Reduced Rate",
					i18n.HR: "Super snižena stopa",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 7, 1),
						Percent: num.MakePercentage(50, 3), // 5.0%
					},
				},
			},
		},
	},
}
