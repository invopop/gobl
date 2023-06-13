// Package mx provides the Mexican tax regime.
package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Tax rates specific to Mexico.
const (
	TaxRateExempt cbc.Key = "exempt"
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.MX,
		Currency: currency.MXN,
		Name: i18n.String{
			i18n.EN: "Mexico",
			i18n.ES: "México",
		},
		Validator:  Validate,
		Calculator: Calculate,
		Categories: []*tax.Category{
			{
				Code: common.TaxCategoryVAT,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.ES: "IVA",
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.ES: "Impuesto al Valor Agregado",
				},
				Retained: false,
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateStandard,
						Name: i18n.String{
							i18n.EN: "Standard Rate",
							i18n.ES: "Tasa General",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(160, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced (Border) Rate",
							i18n.ES: "Tasa Reducida (Fronteriza)",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(80, 3),
							},
						},
					},
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.ES: "Tasa Cero",
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
							i18n.ES: "Exenta",
						},
						Exempt: true,
					},
				},
			},
		},
	}
}

// Validate validates a document against the tax regime.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate performs regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return common.NormalizeTaxIdentity(obj)
	}
	return nil
}
