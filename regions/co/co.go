package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
)

// Local tax categories.
const (
	TaxCategoryIC        org.Code = "IC"  // Impuesto Consumo
	TaxCategoryICA       org.Code = "ICA" // Impuesto de Industria y Comercio
	TaxCategoryINC       org.Code = "INC"
	TaxCategoryReteIVA   org.Code = "RVAT" // ReteIVA
	TaxCategoryReteRenta org.Code = "RR"   // ReteRenta
	TaxCategoryReteICA   org.Code = "RICA" // ReteICA
)

// Keys used in meta data
const (
	KeyPost org.Key = "post"
)

// DIAN official codes to include in stamps.
const (
	StampProviderDIANCUFE org.Key = "dian-cufe"
	StampProviderDIANQR   org.Key = "dian-qr"
)

// Region provides the tax region definition
func Region() *tax.Region {
	return &tax.Region{
		Country:  l10n.CO,
		Currency: "COP",
		Name: i18n.String{
			i18n.EN: "Colombia",
			i18n.ES: "Colombia",
		},
		ValidateDocument:     Validate,
		ValidateTaxIdentity:  ValidateTaxIdentity,
		NormalizeTaxIdentity: NormalizeTaxIdentity,
		Localities: tax.Localities{
			{Code: "AMA", Name: i18n.String{i18n.ES: "Amazonas"}, Meta: org.Meta{KeyPost: "91"}},
			{Code: "ANT", Name: i18n.String{i18n.ES: "Antioquia"}, Meta: org.Meta{KeyPost: "05"}},
			{Code: "ARA", Name: i18n.String{i18n.ES: "Arauca"}, Meta: org.Meta{KeyPost: "81"}},
			{Code: "ATL", Name: i18n.String{i18n.ES: "Atlántico"}, Meta: org.Meta{KeyPost: "08"}},
			{Code: "DC", Name: i18n.String{i18n.ES: "Bogotá"}, Meta: org.Meta{KeyPost: "11"}},
			{Code: "BOL", Name: i18n.String{i18n.ES: "Bolívar"}, Meta: org.Meta{KeyPost: "13"}},
			{Code: "BOY", Name: i18n.String{i18n.ES: "Boyacá"}, Meta: org.Meta{KeyPost: "15"}},
			{Code: "CAL", Name: i18n.String{i18n.ES: "Caldas"}, Meta: org.Meta{KeyPost: "17"}},
			{Code: "CAQ", Name: i18n.String{i18n.ES: "Caquetá"}, Meta: org.Meta{KeyPost: "18"}},
			{Code: "CAS", Name: i18n.String{i18n.ES: "Casanare"}, Meta: org.Meta{KeyPost: "85"}},
			{Code: "CAU", Name: i18n.String{i18n.ES: "Cauca"}, Meta: org.Meta{KeyPost: "19"}},
			{Code: "CES", Name: i18n.String{i18n.ES: "Cesar"}, Meta: org.Meta{KeyPost: "20"}},
			{Code: "CHO", Name: i18n.String{i18n.ES: "Chocó"}, Meta: org.Meta{KeyPost: "27"}},
			{Code: "COR", Name: i18n.String{i18n.ES: "Córdoba"}, Meta: org.Meta{KeyPost: "23"}},
			{Code: "CUN", Name: i18n.String{i18n.ES: "Cundinamarca"}, Meta: org.Meta{KeyPost: "25"}},
			{Code: "GUA", Name: i18n.String{i18n.ES: "Guainía"}, Meta: org.Meta{KeyPost: "94"}},
			{Code: "GUV", Name: i18n.String{i18n.ES: "Guaviare"}, Meta: org.Meta{KeyPost: "95"}},
			{Code: "HUI", Name: i18n.String{i18n.ES: "Huila"}, Meta: org.Meta{KeyPost: "41"}},
			{Code: "LAG", Name: i18n.String{i18n.ES: "La Guajira"}, Meta: org.Meta{KeyPost: "44"}},
			{Code: "MAG", Name: i18n.String{i18n.ES: "Magdalena"}, Meta: org.Meta{KeyPost: "47"}},
			{Code: "MET", Name: i18n.String{i18n.ES: "Meta"}, Meta: org.Meta{KeyPost: "50"}},
			{Code: "NAR", Name: i18n.String{i18n.ES: "Nariño"}, Meta: org.Meta{KeyPost: "52"}},
			{Code: "NSA", Name: i18n.String{i18n.ES: "Norte de Santander"}, Meta: org.Meta{KeyPost: "54"}},
			{Code: "PUT", Name: i18n.String{i18n.ES: "Putumayo"}, Meta: org.Meta{KeyPost: "86"}},
			{Code: "QUI", Name: i18n.String{i18n.ES: "Quindío"}, Meta: org.Meta{KeyPost: "63"}},
			{Code: "RIS", Name: i18n.String{i18n.ES: "Risaralda"}, Meta: org.Meta{KeyPost: "66"}},
			{Code: "SAP", Name: i18n.String{i18n.ES: "San Andrés y Providencia"}, Meta: org.Meta{KeyPost: "88"}},
			{Code: "SAN", Name: i18n.String{i18n.ES: "Santander"}, Meta: org.Meta{KeyPost: "68"}},
			{Code: "SUC", Name: i18n.String{i18n.ES: "Sucre"}, Meta: org.Meta{KeyPost: "70"}},
			{Code: "TOL", Name: i18n.String{i18n.ES: "Tolima"}, Meta: org.Meta{KeyPost: "73"}},
			{Code: "VAC", Name: i18n.String{i18n.ES: "Valle de Cauca"}, Meta: org.Meta{KeyPost: "76"}},
			{Code: "VAU", Name: i18n.String{i18n.ES: "Vaupés"}, Meta: org.Meta{KeyPost: "97"}},
			{Code: "VID", Name: i18n.String{i18n.ES: "Vichada"}, Meta: org.Meta{KeyPost: "99"}},
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
							i18n.ES: "Zero",
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
	}
	return nil
}
