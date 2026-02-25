// Package us provides models for dealing with the United States of America.
package us

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Identification codes unique to the United States.
const (
	IdentityTypeEIN cbc.Code = "EIN" // Employer Identification Number
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "US",
		Currency: currency.USD,
		Name: i18n.String{
			i18n.EN: "United States of America",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The United States does not have a federal value-added tax (VAT) or goods and
				services tax (GST). Instead, sales taxes are levied at the state and local
				level, with rates and rules varying significantly across jurisdictions.

				Sales tax rates range from 0% in states with no sales tax (e.g. Oregon,
				Montana, Delaware, New Hampshire) to combined state and local rates exceeding
				10% in some areas. Sales tax is generally collected by the seller at the point
				of sale and remitted to the relevant state tax authority.

				Businesses are identified by their EIN (Employer Identification Number), a
				9-digit number assigned by the IRS (Internal Revenue Service) in the format
				XX-XXXXXXX. State-level tax registration is separate and varies by
				jurisdiction.

				There is no federal e-invoicing mandate. Both credit notes and debit notes
				are supported for invoice corrections.
			`),
		},
		TimeZone: "America/Chicago", // Around the middle
		Categories: []*tax.CategoryDef{
			//
			// Sales Tax
			//
			{
				Code: tax.CategoryST,
				Name: i18n.String{
					i18n.EN: "ST",
				},
				Title: i18n.String{
					i18n.EN: "Sales Tax",
				},
				Retained: false,
				Rates:    []*tax.RateDef{},
			},
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
	}
}
