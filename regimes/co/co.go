// Package co handles tax regime data for Colombia.
package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Local tax categories.
const (
	TaxCategoryIC        cbc.Code = "IC"  // Impuesto Consumo
	TaxCategoryICA       cbc.Code = "ICA" // Impuesto de Industria y Comercio
	TaxCategoryINC       cbc.Code = "INC"
	TaxCategoryReteIVA   cbc.Code = "RVAT" // ReteIVA
	TaxCategoryReteRenta cbc.Code = "RR"   // ReteRenta
	TaxCategoryReteICA   cbc.Code = "RICA" // ReteICA
)

// DIAN official codes to include in stamps.
const (
	StampProviderDIANCUDE cbc.Key = "dian-cude"
	StampProviderDIANQR   cbc.Key = "dian-qr"
)

// Special keys to use in meta data.
const (
	KeyDIAN cbc.Key = "dian"
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.CO,
		Currency: "COP",
		Name: i18n.String{
			i18n.EN: "Colombia",
			i18n.ES: "Colombia",
		},
		Validator:  Validate,
		Calculator: Calculate,
		Zones:      zones, // see zones.go
		Preceding: &tax.PrecedingDefinitions{ // see preceding.go
			Types: []cbc.Key{
				bill.InvoiceTypeCreditNote,
			},
			Stamps: []cbc.Key{
				StampProviderDIANCUDE,
				StampProviderDIANQR,
			},
			CorrectionMethods: correctionMethodList,
		},
		Categories: []*tax.Category{
			//
			// VAT
			//
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
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.ES: "Cero",
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
							i18n.ES: "Estándar",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2017, 1, 1),
								Percent: num.MakePercentage(190, 3),
							},
							{
								Since:   cal.NewDate(2006, 1, 1),
								Percent: num.MakePercentage(160, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.ES: "Reducido",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2006, 1, 1),
								Percent: num.MakePercentage(50, 3),
							},
						},
					},
				},
			},
			//
			// IC - national
			//
			{
				Code: TaxCategoryIC,
				Name: i18n.String{
					i18n.ES: "IC",
				},
				Desc: i18n.String{
					i18n.EN: "Consumption Tax",
					i18n.ES: "Impuesto sobre Consumo",
				},
				Retained: false,
				Rates:    []*tax.Rate{},
			},
			//
			// ICA - local taxes
			//
			{
				Code: TaxCategoryICA,
				Name: i18n.String{
					i18n.ES: "ICA",
				},
				Desc: i18n.String{
					i18n.EN: "Industry and Commerce Tax",
					i18n.ES: "Impuesto de Industria y Comercio",
				},
				Retained: false,
				Rates:    []*tax.Rate{},
			},
			//
			// ReteIVA
			//
			{
				Code: TaxCategoryReteIVA,
				Name: i18n.String{
					i18n.ES: "ReteIVA",
				},
				Desc: i18n.String{
					i18n.ES: "Retención en la fuente por el Impuesto al Valor Agregado",
				},
				Retained: true,
				Rates:    []*tax.Rate{},
			},
			//
			// ReteICA
			//
			{
				Code: TaxCategoryReteICA,
				Name: i18n.String{
					i18n.ES: "ReteICA",
				},
				Desc: i18n.String{
					i18n.ES: "Retención en la fuente por el Impuesto de Industria y Comercio",
				},
				Retained: true,
				Rates:    []*tax.Rate{},
			},
			//
			// ReteRenta
			//
			{
				Code: TaxCategoryReteRenta,
				Name: i18n.String{
					i18n.ES: "Retefuente",
				},
				Desc: i18n.String{
					i18n.ES: "Retención en la fuente por el Impuesto de la Renta",
				},
				Retained: true,
				Rates:    []*tax.Rate{},
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
	case *org.Party:
		return normalizePartyWithTaxIdentity(obj)
	}
	return nil
}
