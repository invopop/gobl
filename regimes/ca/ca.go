// Package ca provides models for dealing with Canada.
package ca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Tax categories specific for Canada.
const (
	TaxCategoryHST cbc.Code = "HST"
	TaxCategoryPST cbc.Code = "PST"
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "CA",
		Currency:  currency.CAD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Canada",
		},
		TimeZone:   "America/Toronto", // Toronto
		Validator:  Validate,
		Normalizer: Normalize,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Categories: []*tax.CategoryDef{
			//
			// General Sales Tax (GST)
			//
			{
				Code: tax.CategoryGST,
				Name: i18n.String{
					i18n.EN: "GST",
				},
				Title: i18n.String{
					i18n.EN: "General Sales Tax",
				},
				Sources: []*cbc.Source{
					{
						Title: i18n.String{
							i18n.EN: "GST/HST provincial rates table",
						},
						URL: "https://www.canada.ca/en/revenue-agency/services/tax/businesses/topics/gst-hst-businesses/charge-collect-which-rate/calculator.html",
					},
				},
				Retained: false,
				Rates: []*tax.RateDef{
					{
						Rate: tax.RateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
						},
						Description: i18n.String{
							i18n.EN: "Some supplies are zero-rated under the GST, mainly: basic groceries, agricultural products, farm livestock, most fishery products such, prescription drugs and drug-dispensing services, certain medical devices, feminine hygiene products, exports, many transportation services where the origin or destination is outside Canada",
						},
						Values: []*tax.RateValueDef{
							{
								Percent: num.MakePercentage(0, 3),
							},
						},
					},
					{
						Rate: tax.RateGeneral,
						Name: i18n.String{
							i18n.EN: "General rate",
						},
						Description: i18n.String{
							i18n.EN: "For the majority of sales of goods and services: it applies to all products or services for which no other rate is expressly provided.",
						},

						Values: []*tax.RateValueDef{
							{
								Since:   cal.NewDate(2022, 1, 1),
								Percent: num.MakePercentage(5, 2),
							},
						},
					},
				},
			},
			//
			// Harmonized Sales Tax (HST)
			//
			{
				Code: TaxCategoryHST,
				Name: i18n.String{
					i18n.EN: "HST",
				},
				Title: i18n.String{
					i18n.EN: "Harmonized Sales Tax",
				},
				// TODO: determine local rates
				Rates: []*tax.RateDef{},
			},

			//
			// Provincial Sales Tax (PST)
			//
			{
				Code: TaxCategoryPST,
				Name: i18n.String{
					i18n.EN: "PST",
				},
				Title: i18n.String{
					i18n.EN: "Provincial Sales Tax",
				},
				// TODO: determine local rates
				Rates: []*tax.RateDef{},
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
