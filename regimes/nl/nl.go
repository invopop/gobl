// Package nl provides the Dutch region definition
package nl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the Dutch region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "NL",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "The Netherlands",
			i18n.NL: "Nederland",
		},
		TimeZone:   "Europe/Amsterdam",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Categories: []*tax.CategoryDef{
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
				Keys:     tax.GlobalVATKeys(),
				Rates: []*tax.RateDef{
					{
						Keys: []cbc.Key{tax.KeyStandard},
						Rate: tax.RateGeneral,
						Name: i18n.String{
							i18n.EN: "General Rate",
							i18n.NL: "Algemeen Tarief",
						},
						Values: []*tax.RateValueDef{
							{
								Percent: num.MakePercentage(210, 3),
							},
						},
					},
					{
						Keys: []cbc.Key{tax.KeyStandard},
						Rate: tax.RateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.NL: "Gereduceerd Tarief",
						},
						Values: []*tax.RateValueDef{
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
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}

// Normalize performs region specific calculations on the document.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
