package gb

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Regime provides the tax region definition
func Regime() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.GB,
		Currency: "GBP",
		Name: i18n.String{
			i18n.EN: "United Kingdom",
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
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
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
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2011, 1, 4),
								Percent: num.MakePercentage(50, 3),
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
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
