// Package pt provides models for dealing with the Portuguese tax regime.
package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
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

// Zone code definitions for Portugal based on districts and
// autonomous regions based on ISO 3166-2:PT.
const (
	ZoneAveiro         l10n.Code = "01"
	ZoneBeja           l10n.Code = "02"
	ZoneBraga          l10n.Code = "03"
	ZoneBraganca       l10n.Code = "04"
	ZoneCasteloBranco  l10n.Code = "05"
	ZoneCoimbra        l10n.Code = "06"
	ZoneEvora          l10n.Code = "07"
	ZoneFaro           l10n.Code = "08"
	ZoneGuarda         l10n.Code = "09"
	ZoneLeiria         l10n.Code = "10"
	ZoneLisboa         l10n.Code = "11"
	ZonePortalegre     l10n.Code = "12"
	ZonePorto          l10n.Code = "13"
	ZoneSantarem       l10n.Code = "14"
	ZoneSetubal        l10n.Code = "15"
	ZoneVianaDoCastelo l10n.Code = "16"
	ZoneVilaReal       l10n.Code = "17"
	ZoneViseu          l10n.Code = "18"
	ZoneAzores         l10n.Code = "20" // Autonomous Region
	ZoneMadeira        l10n.Code = "30" // Autonomous Region
)

// New instantiates a new Portugal regime for the given zone.
func New() *tax.Regime {
	return &tax.Regime{
		Country:   l10n.PT,
		Currency:  currency.EUR,
		Validator: Validate,
		Name: i18n.String{
			i18n.EN: "Portugal",
			i18n.PT: "Portugal",
		},
		Zones: zones,
		Categories: []*tax.Category{
			// VAT
			{
				Code: common.TaxCategoryVAT,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.PT: "IVA",
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.PT: "Imposto sobre o Valor Acrescentado",
				},
				Retained: false,
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateStandard,
						Name: i18n.String{
							i18n.EN: "Standard Rate",
							i18n.PT: "Tipo Geral",
						},
						Values: []*tax.RateValue{
							{
								Zones:   []l10n.Code{ZoneAzores},
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(160, 3),
							},
							{
								Zones:   []l10n.Code{ZoneAzores},
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(180, 3),
							},
							{
								Zones:   []l10n.Code{ZoneMadeira},
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(220, 3),
							},
							{
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(230, 3),
							},
						},
					},
					{
						Key: common.TaxRateIntermediate,
						Name: i18n.String{
							i18n.EN: "Intermediate Rate",
							i18n.PT: "Taxa Intermédia", //nolint:misspell
						},
						Values: []*tax.RateValue{
							{
								Zones:   []l10n.Code{ZoneAzores},
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(90, 3),
							},
							{
								Zones:   []l10n.Code{ZoneMadeira},
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(120, 3),
							},
							{
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(130, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.PT: "Taxa Reduzida",
						},
						Values: []*tax.RateValue{
							{
								Zones:   []l10n.Code{ZoneAzores},
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(40, 3),
							},
							{
								Zones:   []l10n.Code{ZoneMadeira},
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(50, 3),
							},
							{
								Since:   cal.NewDate(2011, 1, 1),
								Percent: num.MakePercentage(60, 3),
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
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will attempt to clean the object passed to it.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
