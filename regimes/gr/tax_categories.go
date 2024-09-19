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
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Value Added Tax/Goods and Services Tax (VAT/GST) (1976-2023)",
					i18n.EL: "Φόρος Προστιθέμενης Αξίας/Φόρος Αγαθών και Υπηρεσιών (ΦΠΑ/GST) (1976-2023)",
				},
				URL: "https://www.oecd.org/tax/tax-policy/tax-database/",
			},
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard rate",
					i18n.EL: "Κανονικός συντελεστής",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(24, 2),
					},
				},
			},
			{
				Key: tax.RateReduced,
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
				Key: tax.RateSuperReduced,
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
				Key: tax.RateStandard.With(TaxRateIsland),
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
				Key: tax.RateReduced.With(TaxRateIsland),
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
				Key: tax.RateSuperReduced.With(TaxRateIsland),
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
			{
				Key:    tax.RateExempt,
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.EL: "Απαλλαγή",
				},
			},
		},
	},
}
