package it

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// New instantiates a new Italian regime.
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.IT,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "Italy",
			i18n.IT: "Italia",
		},
		Validator:  Validate,
		Calculator: Calculate,
		Zones:      zones, // see zones.go
		Categories: []*tax.Category{
			{
				Code:     common.TaxCategoryVAT,
				Retained: false,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.IT: "IVA",
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.IT: "Imposta sul Valore Aggiunto",
				},
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.IT: "Aliquota Zero",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(0, 3),
							},
						},
					},
					{
						Key: common.TaxRateSuperReduced,
						Name: i18n.String{
							i18n.EN: "Minimum Rate",
							i18n.IT: "Aliquota Minima",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(4, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.IT: "Aliquota Ridotta",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(100, 3),
							},
						},
					},
					{
						Key: common.TaxRateStandard,
						Name: i18n.String{
							i18n.EN: "Ordinary Rate",
							i18n.IT: "Aliquota Ordinaria",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(220, 3),
							},
						},
					},
				},
			},
		},
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {

	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will perform any regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
