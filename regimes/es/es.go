// Package es provides tax regime support for Spain.
package es

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

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryIRPF cbc.Code = "IRPF"
	TaxCategoryIGIC cbc.Code = "IGIC"
	TaxCategoryIPSI cbc.Code = "IPSI"
)

// Specific tax rate codes.
const (
	// IRPF non-standard Rates (usually for self-employed)
	TaxRatePro                cbc.Key = "pro"                 // Professional Services
	TaxRateProStart           cbc.Key = "pro-start"           // Professionals, first 2 years
	TaxRateModules            cbc.Key = "modules"             // Module system
	TaxRateAgriculture        cbc.Key = "agriculture"         // Agricultural
	TaxRateAgricultureSpecial cbc.Key = "agriculture-special" // Agricultural special
	TaxRateCapital            cbc.Key = "capital"             // Rental or Interest

	// Special tax rate surcharge extension
	TaxRateEquivalence cbc.Key = "eqs"
)

// Official stamps or codes validated by government agencies
const (
	// TicketBAI (Basque Country) codes used for stamps.
	StampProviderTBAICode cbc.Key = "tbai-code"
	StampProviderTBAIQR   cbc.Key = "tbai-qr"
)

// Inbox key and role definitions
const (
	InboxKeyFACE cbc.Key = "face"

	// Main roles defined in FACE
	InboxRoleFiscal    cbc.Key = "fiscal"    // Fiscal / 01
	InboxRoleRecipient cbc.Key = "recipient" // Receptor / 02
	InboxRolePayer     cbc.Key = "payer"     // Pagador / 03
	InboxRoleCustomer  cbc.Key = "customer"  // Comprador / 04

)

// Custom keys used typically in meta information.
const (
	KeyAddressCode                 cbc.Key = "post"
	KeyFacturaE                    cbc.Key = "facturae"
	KeyFacturaETaxTypeCode         cbc.Key = "facturae-tax-type-code"
	KeyFacturaEInvoiceDocumentType cbc.Key = "facturae-invoice-document-type"
	KeyFacturaEInvoiceClass        cbc.Key = "facturae-invoice-class"
	KeyTicketBAICausaExencion      cbc.Key = "ticketbai-causa-exencion"
	KeyTicketBAIIDType             cbc.Key = "ticketbai-id-type"
)

// New provides the Spanish tax regime definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.ES,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "Spain",
			i18n.ES: "España",
		},
		Validator:     Validate,
		Calculator:    Calculate,
		Zones:         zones,
		IdentityTypes: taxIdentityTypes,
		Tags:          invoiceTags,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Preceding: &tax.PrecedingDefinitions{
			Corrections:       correctionList,
			CorrectionMethods: correctionMethodList,
		},
		Categories: []*tax.Category{
			//
			// VAT
			//
			{
				Code:     common.TaxCategoryVAT,
				Retained: false,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.ES: "IVA",
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.ES: "Impuesto sobre el Valor Añadido",
				},
				Tags: vatTaxTags,
				Meta: cbc.Meta{
					KeyFacturaETaxTypeCode: "01",
				},
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.ES: "Tipo Cero",
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
							i18n.ES: "Tipo General",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2012, 9, 1),
								Percent: num.MakePercentage(210, 3),
							},
							{
								Since:   cal.NewDate(2010, 7, 1),
								Percent: num.MakePercentage(180, 3),
							},
							{
								Since:   cal.NewDate(1995, 1, 1),
								Percent: num.MakePercentage(160, 3),
							},
							{
								Since:   cal.NewDate(1993, 1, 1),
								Percent: num.MakePercentage(150, 3),
							},
						},
					},
					{
						Key: common.TaxRateStandard.With(TaxRateEquivalence),
						Name: i18n.String{
							i18n.EN: "Standard Rate + Equivalence Surcharge",
							i18n.ES: "Tipo General + Recargo de Equivalencia",
						},
						Values: []*tax.RateValue{
							{
								Since:     cal.NewDate(2012, 9, 1),
								Percent:   num.MakePercentage(210, 3),
								Surcharge: num.NewPercentage(52, 3),
							},
							{
								Since:     cal.NewDate(2010, 7, 1),
								Percent:   num.MakePercentage(180, 3),
								Surcharge: num.NewPercentage(40, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.ES: "Tipo Reducido",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2012, 9, 1),
								Percent: num.MakePercentage(100, 3),
							},
							{
								Since:   cal.NewDate(2010, 7, 1),
								Percent: num.MakePercentage(80, 3),
							},
							{
								Since:   cal.NewDate(1995, 1, 1),
								Percent: num.MakePercentage(70, 3),
							},
							{
								Since:   cal.NewDate(1993, 1, 1),
								Percent: num.MakePercentage(60, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced.With(TaxRateEquivalence),
						Name: i18n.String{
							i18n.EN: "Reduced Rate + Equivalence Surcharge",
							i18n.ES: "Tipo Reducido + Recargo de Equivalencia",
						},
						Values: []*tax.RateValue{
							{
								Since:     cal.NewDate(2012, 9, 1),
								Percent:   num.MakePercentage(100, 3),
								Surcharge: num.NewPercentage(14, 3),
							},
							{
								Since:     cal.NewDate(2010, 7, 1),
								Percent:   num.MakePercentage(80, 3),
								Surcharge: num.NewPercentage(10, 3),
							},
						},
					},
					{
						Key: common.TaxRateSuperReduced,
						Name: i18n.String{
							i18n.EN: "Super-Reduced Rate",
							i18n.ES: "Tipo Superreducido",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(1995, 1, 1),
								Percent: num.MakePercentage(40, 3),
							},
							{
								Since:   cal.NewDate(1993, 1, 1),
								Percent: num.MakePercentage(30, 3),
							},
						},
					},
					{
						Key: common.TaxRateSuperReduced.With(TaxRateEquivalence),
						Name: i18n.String{
							i18n.EN: "Super-Reduced Rate + Equivalence Surcharge",
							i18n.ES: "Tipo Superreducido + Recargo de Equivalencia",
						},
						Values: []*tax.RateValue{
							{
								Since:     cal.NewDate(1995, 1, 1),
								Percent:   num.MakePercentage(40, 3),
								Surcharge: num.NewPercentage(5, 3),
							},
						},
					},
				},
			},

			//
			// IGIC
			//
			{
				Code:     TaxCategoryIGIC,
				Retained: false,
				Name: i18n.String{
					i18n.EN: "IGIC",
					i18n.ES: "IGIC",
				},
				Meta: cbc.Meta{
					KeyFacturaETaxTypeCode: "03",
				},
				Desc: i18n.String{
					i18n.EN: "Canary Island General Indirect Tax",
					i18n.ES: "Impuesto General Indirecto Canario",
				},
				// This is a subset of the possible rates.
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.ES: "Tipo Cero",
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
							i18n.ES: "Tipo General",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(70, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.ES: "Tipo Reducido",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(30, 3),
							},
						},
					},
				},
			},

			//
			// IPSI
			//
			{
				Code:     TaxCategoryIPSI,
				Retained: false,
				Name: i18n.String{
					i18n.EN: "IPSI",
					i18n.ES: "IPSI",
				},
				Meta: cbc.Meta{
					KeyFacturaETaxTypeCode: "02",
				},
				Desc: i18n.String{
					i18n.EN: "Production, Services, and Import Tax",
					i18n.ES: "Impuesto sobre la Producción, los Servicios y la Importación",
				},
				// IPSI rates are complex and don't align well regular rates. Users are
				// recommended to include whatever percentage applies to their situation
				// directly in the invoice.
				Rates: []*tax.Rate{},
			},

			//
			// IRPF
			//
			{
				Code:     TaxCategoryIRPF,
				Retained: true,
				Name: i18n.String{
					i18n.EN: "IRPF",
					i18n.ES: "IRPF",
				},
				Meta: cbc.Meta{
					KeyFacturaETaxTypeCode: "04",
				},
				Desc: i18n.String{
					i18n.EN: "Personal income tax.",
					i18n.ES: "Impuesto sobre la renta de las personas físicas.",
				},
				Rates: []*tax.Rate{
					{
						Key: TaxRatePro,
						Name: i18n.String{
							i18n.EN: "Professional Rate",
							i18n.ES: "Profesionales",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2015, 7, 12),
								Percent: num.MakePercentage(150, 3),
							},
							{
								Since:   cal.NewDate(2015, 1, 1),
								Percent: num.MakePercentage(190, 3),
							},
							{
								Since:   cal.NewDate(2012, 9, 1),
								Percent: num.MakePercentage(210, 3),
							},
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(150, 3),
							},
						},
					},
					{
						Key: TaxRateProStart,
						Name: i18n.String{
							i18n.EN: "Professional Starting Rate",
							i18n.ES: "Profesionales Inicio",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(70, 3),
							},
						},
					},
					{
						Key: TaxRateCapital,
						Name: i18n.String{
							i18n.EN: "Rental or Interest Capital",
							i18n.ES: "Alquileres o Intereses de Capital",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(190, 3),
							},
						},
					},
					{
						Key: TaxRateModules,
						Name: i18n.String{
							i18n.EN: "Modules Rate",
							i18n.ES: "Tipo Modulos",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(10, 3),
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

// Calculate will perform any regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
