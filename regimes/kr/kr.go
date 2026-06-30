// Package kr provides the tax regime definition for South Korea.
package kr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for South Korea.
const CountryCode = "KR"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("kr", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New instantiates a new South Korean tax regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.KRW,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "South Korea",
			i18n.KO: "대한민국",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				South Korea applies VAT (부가가치세), administered by the National Tax
				Service (국세청), on most goods and services, with zero-rating for exports
				and certain internationally supplied services.

				Businesses are identified by their Business Registration Number
				(사업자등록번호), a 10-digit number with a check digit. It is the identifier
				required on Korean tax invoices (세금계산서).
			`),
		},
		TimeZone:   "Asia/Seoul",
		Categories: taxCategories,
		// Korean VAT law provides for revised tax invoices (수정세금계산서) that adjust a
		// previously issued invoice, either downwards (e.g. returns, discounts) or
		// upwards (e.g. additional supply), mapped here to credit and debit notes.
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
