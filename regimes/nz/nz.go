// Package nz provides the tax regime definition for New Zealand.
package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for New Zealand.
const CountryCode = "NZ"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("nz", rules.GOBL.Add(CountryCode),
		billInvoiceRules(),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New instantiates a new New Zealand tax regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.NZD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "New Zealand",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Inland Revenue - GST"),
				URL:   "https://www.ird.govt.nz/gst",
			},
		},
		TimeZone: "Pacific/Auckland",
		Description: i18n.String{
			i18n.EN: "New Zealand applies a single-rate Goods and Services Tax (GST), administered by Inland Revenue (IRD). Most taxable supplies are charged at the standard rate, and businesses are identified by an IRD number validated with a modulus-11 check digit.",
		},
		Categories: taxCategories(),
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
