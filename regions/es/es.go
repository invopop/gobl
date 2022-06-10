package es

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

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryIRPF tax.Code = "IRPF"
	TaxCategoryIGIC tax.Code = "IGIC"
	TaxCategoryIPSI tax.Code = "IPSI"
)

// Specific tax rate codes.
const (
	// IRPF non-standard Rates (usually for self-employed)
	TaxRatePro                tax.Key = "pro"                 // Professional Services
	TaxRateProStart           tax.Key = "pro-start"           // Professionals, first 2 years
	TaxRateModules            tax.Key = "modules"             // Module system
	TaxRateAgriculture        tax.Key = "agriculture"         // Agricultural
	TaxRateAgricultureSpecial tax.Key = "agriculture-special" // Agricultural special
	TaxRateCapital            tax.Key = "capital"             // Rental or Interest

	// Special tax rate surcharge extension
	TaxRateEquivalence tax.Key = "eqs"
)

// Scheme key definitions
const (
	SchemeSimplified      tax.Key = "simplified"
	SchemeCustomerIssued  tax.Key = "customer-issued"
	SchemeTravelAgency    tax.Key = "travel-agency"
	SchemeSecondHandGoods tax.Key = "second-hand-goods"
	SchemeArt             tax.Key = "art"
	SchemeAntiques        tax.Key = "antiques"
	SchemeCashBasis       tax.Key = "cash-basis"
)

// New provides the Spanish region definition
func New() *tax.Region {
	return &tax.Region{
		Country:  l10n.ES,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "Spain",
			i18n.ES: "España",
		},
		ValidateDocument: Validate,
		Schemes: []*tax.Scheme{
			// Reverse Charge Scheme
			{
				Key: common.SchemeReverseCharge,
				Name: i18n.String{
					i18n.EN: "Reverse Charge",
					i18n.ES: "Inversión del sujeto pasivo",
				},
				Categories: []tax.Code{
					common.TaxCategoryVAT,
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(common.SchemeReverseCharge),
					Text: "Reverse Charge / Inversión del sujeto pasivo.",
				},
			},
			// Customer Rates Scheme (digital goods)
			{
				Key: common.SchemeCustomerRates,
				Name: i18n.String{
					i18n.EN: "Customer Country Rates",
					i18n.ES: "Tasas del País del Cliente",
				},
				Description: i18n.String{
					i18n.EN: "Use the customers country to determine tax rates.",
				},
			},
			// Simplified Regime
			{
				Key: SchemeSimplified,
				Name: i18n.String{
					i18n.EN: "Simplified tax scheme",
					i18n.ES: "Contribuyente en régimen simplificado",
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(SchemeSimplified),
					Text: "Factura expedida por contibuyente en régimen simplificado.",
				},
			},
			// Customer issued invoices
			{
				Key: SchemeCustomerIssued,
				Name: i18n.String{
					i18n.EN: "Customer issued invoice",
					i18n.ES: "Facturación por el destinatario",
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(SchemeCustomerIssued),
					Text: "Facturación por el destinatario.",
				},
			},
			// Travel agency
			{
				Key: SchemeTravelAgency,
				Name: i18n.String{
					i18n.EN: "Special scheme for travel agencies",
					i18n.ES: "Régimen especial de las agencias de viajes",
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(SchemeTravelAgency),
					Text: "Régimen especial de las agencias de viajes.",
				},
			},
			// Secondhand stuff
			{
				Key: SchemeSecondHandGoods,
				Name: i18n.String{
					i18n.EN: "Special scheme for second-hand goods",
					i18n.ES: "Régimen especial de los bienes usados",
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(SchemeSecondHandGoods),
					Text: "Régimen especial de los bienes usados.",
				},
			},
			// Art
			{
				Key: SchemeArt,
				Name: i18n.String{
					i18n.EN: "Special scheme of works of art",
					i18n.ES: "Régimen especial de los objetos de arte",
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(SchemeArt),
					Text: "Régimen especial de los objetos de arte.",
				},
			},
			// Antiques
			{
				Key: SchemeAntiques,
				Name: i18n.String{
					i18n.EN: "Special scheme of antiques and collectables",
					i18n.ES: "Régimen especial de las antigüedades y objetos de colección",
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(SchemeAntiques),
					Text: "Régimen especial de las antigüedades y objetos de colección.",
				},
			},
			// Special Regime of "Cash Criteria"
			{
				Key: SchemeCashBasis,
				Name: i18n.String{
					i18n.EN: "Special scheme on cash basis",
					i18n.ES: "Régimen especial del criterio de caja",
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
					Src:  string(SchemeCashBasis),
					Text: "Régimen especial del criterio de caja.",
				},
			},
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
					i18n.ES: "Impuesto sobre el Valor Añadido",
				},
				Retained: false,
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.ES: "Tipo Zero",
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
							i18n.EN: "Standard Rate + Equivalence",
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
							i18n.EN: "Reduced Rate + Surcharge",
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
							i18n.EN: "Super-Reduced Rate + Equivalence",
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
			// IRPF
			//
			{
				Code:     TaxCategoryIRPF,
				Retained: true,
				Name: i18n.String{
					i18n.EN: "IRPF",
					i18n.ES: "IRPF",
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
							i18n.ES: "Professionales",
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
							i18n.ES: "Professionales Inicio",
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
	}
	return nil
}
