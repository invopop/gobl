// Package gb provides the United Kingdom tax regime.
package gb

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Identification code types unique to the United Kingdom.
const (
	IdentityTypeCRN cbc.Code = "CRN" // Company Registration Number
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.GB,
		Currency: "GBP",
		Name: i18n.String{
			i18n.EN: "United Kingdom",
		},
		Validator: Validate,
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
