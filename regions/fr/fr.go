package fr

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
)

// Region provides the tax region definition
func Region() *tax.Region {
	return &tax.Region{
		Country:  l10n.FR,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "France",
			i18n.FR: "La France",
		},
		ValidateDocument: Validate,
		Categories: []*tax.Category{
			//
			// VAT
			//
			{
				Code: common.TaxCategoryVAT,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.FR: "TVA",
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.FR: "Taxe sur la Valeur Ajoutée",
				},
				Retained: false,
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(0, 3),
							},
						},
					},
					{
						Key: common.TaxRateStandard,
						Name: i18n.String{
							i18n.EN: "Standard Rate",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2011, 1, 4),
								Percent: num.MakePercentage(200, 3),
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
	return nil
}
