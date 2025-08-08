package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// TaxRateIsland is used as a suffix to regular tax rates in order to denote
// the reduced rates that apply to islands.
const TaxRateIsland cbc.Key = "island"

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.EL: "ΦΠΑ",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.EL: "Φόρος προστιθέμενης αξίας",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("VAT Rates"),
				URL:   "https://www.gov.gr/en/sdg/taxes/vat/general/basic-vat-rates",
			},
			{
				Title: i18n.String{
					i18n.EN: "Value Added Tax/Goods and Services Tax (VAT/GST) (1976-2023)",
					i18n.EL: "Φόρος Προστιθέμενης Αξίας/Φόρος Αγαθών και Υπηρεσιών (ΦΠΑ/GST) (1976-2023)",
				},
				URL: "https://www.oecd.org/tax/tax-policy/tax-database/",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General rate",
					i18n.EL: "Κανονικός συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(24, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced rate",
					i18n.EL: "Μειωμένος συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(13, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-reduced rate",
					i18n.EL: "Υπερμειωμένος συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(6, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral.With(TaxRateIsland),
				Name: i18n.String{
					i18n.EN: "Standard rate (Island)",
					i18n.EL: "Κανονικός συντελεστής (Νησί)",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(17, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral.With(TaxRateIsland),
				Name: i18n.String{
					i18n.EN: "Reduced rate (Island)",
					i18n.EL: "Μειωμένος συντελεστής (Νησί)",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(9, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced.With(TaxRateIsland),
				Name: i18n.String{
					i18n.EN: "Super-reduced rate (Island)",
					i18n.EL: "Υπερμειωμένος συντελεστής (Νησί)",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(4, 2),
					},
				},
			},
		},
	},
}
