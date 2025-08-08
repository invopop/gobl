package be

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
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "VAT rates | FPS Finance",
				},
				URL: "https://finance.belgium.be/en/enterprises/vat/vat-obligation/rates-and-calculation/vat-rates#q1",
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
					i18n.NL: "Algemeen tarief",
					i18n.FR: "Taux normal",
					i18n.DE: "Allgemeiner Satz",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(210, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.NL: "Tussentarief",
					i18n.FR: "Taux intermédiaire",
					i18n.DE: "Mittlerer Satz",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.NL: "Verlaagd Tarief",
					i18n.FR: "Taux réduit",
					i18n.DE: "Ermäßigter Satz",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(60, 3),
					},
				},
			},
		},
	},
}
