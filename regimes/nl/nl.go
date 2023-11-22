// Package nl provides the Dutch region definition
package nl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// New provides the Dutch region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.NL,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "The Netherlands",
			i18n.NL: "Nederland",
		},
		TimeZone:   "Europe/Amsterdam",
		Validator:  Validate,
		Calculator: Calculate,
		Scenarios: []*tax.ScenarioSet{
			common.InvoiceScenarios(),
		},
		Tags: common.InvoiceTags(),
		Categories: []*tax.Category{
			//
			// VAT
			//
			{
				Code: tax.CategoryVAT,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.NL: "BTW",
				},
				Title: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.NL: "Belasting Toegevoegde Waarde",
				},
				Retained: false,
				Rates: []*tax.Rate{
					{
						Key: tax.RateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.NL: `0%-tarief`,
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(0, 3),
							},
						},
					},
					{
						Key: tax.RateStandard,
						Name: i18n.String{
							i18n.EN: "Standard Rate",
							i18n.NL: "Standaardtarief",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(210, 3),
							},
						},
					},
					{
						Key: tax.RateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.NL: "Gereduceerd Tarief",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(90, 3),
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

// Calculate performs region specific calculations on the document.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return NormalizeTaxIdentity(obj)
	}
	return nil
}
