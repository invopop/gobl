// Package us provides models for dealing with the United States of America.
package us

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Identification codes unique to the United States.
const (
	IdentityTypeEIN cbc.Code = "EIN" // Employer Identification Number
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.US,
		Currency: currency.USD,
		Name: i18n.String{
			i18n.EN: "United States of America",
		},
		TimeZone:  "America/Chicago", // Around the middle
		Validator: Validate,
		Tags:      common.InvoiceTags(),
		Categories: []*tax.Category{
			//
			// Sales Tax
			//
			{
				Code: common.TaxCategoryST,
				Name: i18n.String{
					i18n.EN: "ST",
				},
				Title: i18n.String{
					i18n.EN: "Sales Tax",
				},
				Retained: false,
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
