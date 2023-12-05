package pl

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Tax rates specific to Poland.
const (
	TaxRateNotPursuant       cbc.Key = "np"
	TaxRateNotPursuantArt100 cbc.Key = "np-art100sec1point4"
)

var taxCategories = []*tax.Category{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.PL: "VAT",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.PL: "Podatek od Wartości Dodanej",
		},
		Retained: false,
		Rates: []*tax.Rate{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.PL: "Stawka Podstawowa",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(230, 3),
						Since:   cal.NewDate(2011, 1, 1),
					},
					{
						Percent: num.MakePercentage(220, 3),
						Since:   cal.NewDate(1993, 7, 8),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "First Reduced Rate",
					i18n.PL: "Stawka Obniżona Pierwsza",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(80, 3),
						Since:   cal.NewDate(2011, 1, 1),
					},
					{
						Percent: num.MakePercentage(70, 3),
						Since:   cal.NewDate(2000, 9, 4),
					},
				},
			},
			{
				Key: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Second Reduced Rate",
					i18n.PL: "Stawka Obniżona Druga",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(50, 3),
						Since:   cal.NewDate(2011, 1, 1),
					},
					{
						Percent: num.MakePercentage(30, 3),
						Since:   cal.NewDate(2000, 9, 4),
					},
				},
			},
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.PL: "Stawka Zerowa",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
				Extensions: []cbc.Key{
					ExtKeyKSeFVATZero,
				},
			},
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PL: "Zwolnione",
				},
				Exempt:     true,
				Extensions: []cbc.Key{},
			},
			{
				Key: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Special Rate",
					i18n.PL: "Stawka Specjalna",
				},
				Extensions: []cbc.Key{
					ExtKeyKSeFVATSpecial,
				},
			},

			/*
				 * Still working on refactoring these...
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
			*/
		},
	},
}
