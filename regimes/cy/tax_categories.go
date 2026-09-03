package cy

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
			i18n.EL: "ΦΠΑ",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.EL: "Φόρος Προστιθέμενης Αξίας",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Cyprus Tax Department - VAT Rates"),
				URL:   "https://www.gov.cy/mof-tax/documents/syntelestes-f-p-a/",
			},
			{
				Title: i18n.NewString("European Commission - VAT Rates"),
				URL:   "https://europa.eu/youreurope/business/finance-and-tax/vat/vat-rules-rates/index_en.htm",
			},
			{
				Title: i18n.NewString("Cyprus Treasury - Historical VAT Rates"),
				URL:   "https://www.gov.cy/media/sites/42/2026/01/9Α_Θέματα-Φόρου-Προστιθέμενης-Αξίας-Φ.Π.Α.-1-2.pdf",
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
					i18n.EL: "Κανονικός συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(19, 2), Since: cal.NewDate(2014, 1, 13)},
					{Percent: num.MakePercentage(18, 2), Since: cal.NewDate(2013, 1, 14)},
					{Percent: num.MakePercentage(17, 2), Since: cal.NewDate(2012, 3, 1)},
					{Percent: num.MakePercentage(15, 2), Since: cal.NewDate(2003, 1, 1)},
					{Percent: num.MakePercentage(13, 2), Since: cal.NewDate(2002, 7, 1)},
					{Percent: num.MakePercentage(10, 2), Since: cal.NewDate(2000, 7, 1)},
					{Percent: num.MakePercentage(8, 2), Since: cal.NewDate(1993, 10, 1)},
					{Percent: num.MakePercentage(5, 2), Since: cal.NewDate(1992, 7, 1)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Higher Reduced Rate",
					i18n.EL: "Υψηλότερος μειωμένος συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(9, 2), Since: cal.NewDate(2014, 1, 13)},
					{Percent: num.MakePercentage(8, 2), Since: cal.NewDate(2005, 8, 1)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Lower Reduced Rate",
					i18n.EL: "Χαμηλότερος μειωμένος συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(5, 2), Since: cal.NewDate(2000, 7, 1)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.EL: "Υπερμειωμένος συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(3, 2), Since: cal.NewDate(2023, 7, 21)},
				},
			},
		},
	},
}
