package pl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Tax rates specific to Poland.
const (
	TaxRateExempt            cbc.Key = "exempt"
	TaxRateNotPursuant       cbc.Key = "np"
	TaxRateNotPursuantArt100 cbc.Key = "np-art100sec1point4"
	TaxRateReverseCharge     cbc.Key = "reverse-charge"
	TaxRateZeroWDT           cbc.Key = "zero-wdt"
	TaxRateZeroDomestic      cbc.Key = "zero-domestic"
	TaxRateZeroExport        cbc.Key = "zero-export"
)

var taxCategories = []*tax.Category{
	{
		Code: common.TaxCategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.PL: "VAT",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.PL: "Podatek od wartości dodanej",
		},
		Retained: false,
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.PL: "Stawka Podstawowa",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(230, 3),
					},
					{
						Percent: num.MakePercentage(220, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "First Reduced Rate",
					i18n.PL: "Stawka Obniżona Pierwsza",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(80, 3),
					},
					{
						Percent: num.MakePercentage(70, 3),
					},
				},
			},
			{
				Key: common.TaxRateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Second Reduced Rate",
					i18n.PL: "Stawka Obniżona Druga",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			{
				Key: common.TaxRateSpecial,
				Name: i18n.String{
					i18n.EN: "Lump sum taxi rate",
					i18n.PL: "Ryczałt dla taksówek",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(40, 3),
					},
					{
						Percent: num.MakePercentage(30, 3),
					},
				},
			},
			{
				Key: TaxRateZeroWDT,
				Name: i18n.String{
					i18n.EN: "Zero Rate - WDT",
					i18n.PL: "Stawka Zerowa - WDT",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: TaxRateZeroDomestic,
				Name: i18n.String{
					i18n.EN: "Zero Rate - domestic",
					i18n.PL: "Stawka Zerowa - krajowe",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: TaxRateZeroExport,
				Name: i18n.String{
					i18n.EN: "Zero Rate - export",
					i18n.PL: "Stawka Zerowa - export",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: TaxRateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PL: "Zwolnione z opodatkowania",
				},
				Exempt: true,
			},
			{
				Key: TaxRateNotPursuant,
				Name: i18n.String{
					i18n.EN: "Not pursuant, pursuant to art100 section 1 point4",
					i18n.PL: "Niepodlegające opodatkowaniu na postawie wyłączeniem art100 sekcja 1 punkt 4",
				},
				Exempt: true,
			},
			{
				Key: TaxRateNotPursuantArt100,
				Name: i18n.String{
					i18n.EN: "Not pursuant excluding art100 section 1 point4",
					i18n.PL: "Niepodlegające opodatkowaniu z wyłączeniem art100 sekcja 1 punkt 4",
				},
				Exempt: true,
			},
			{
				Key: TaxRateReverseCharge,
				Name: i18n.String{
					i18n.EN: "Reverse Charge",
					i18n.PL: "Odwrotne obciążenie",
				},
				Exempt: true,
			},
		},
	},
}
