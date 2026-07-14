package si

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
			i18n.SL: "DDV",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.SL: "Davek na dodano vrednost",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.SL: "Splošna stopnja",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 7, 1),
						Percent: num.MakePercentage(220, 3),
					},
					{
						Since:   cal.NewDate(2002, 1, 1),
						Percent: num.MakePercentage(200, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.SL: "Nižja stopnja",
				},
				// Food, medicines, passenger transport, accommodation, and
				// other goods and services listed in Annex I to the VAT Act.
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(95, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Special Reduced Rate for Publications",
					i18n.SL: "Posebna nižja stopnja",
				},
				// Books, newspapers and periodicals, in print or electronic
				// form, in force since 1 January 2020.
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Financial Administration of the Republic of Slovenia (FURS) - VAT",
				},
				URL: "https://www.fu.gov.si/en/taxes_and_other_duties/areas_of_work/value_added_tax_vat",
				At:  cal.NewDateTime(2026, 7, 3, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.EN: "European Commission - VAT rates applied in the Member States of the European Union",
				},
				URL: "https://taxation-customs.ec.europa.eu/system/files/2021-06/vat_rates_en.pdf",
				At:  cal.NewDateTime(2026, 7, 3, 0, 0, 0),
			},
		},
	},
}
