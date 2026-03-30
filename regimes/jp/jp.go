// Package jp provides the tax regime definition for Japan.
package jp

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

// New provides the tax regime definition for Japan.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "JP",
		Currency:  currency.JPY,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Japan",
			i18n.JA: "日本", // Nihon
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Japan's tax system is administered by the National Tax Agency (NTA).
				The Consumption Tax (消費税, Shōhizei) functions as a VAT and applies
				at standard (10%) and reduced (8%) rates since October 2019.

				Businesses are identified by a 13-digit Corporate Number (法人番号,
				Hōjin Bangō). The Qualified Invoice System (インボイス制度), effective
				since October 2023, requires registered businesses to include a
				Registration Number (T + Corporate Number) on invoices for buyers
				to claim input tax credits.

				Credit notes are supported as Qualified Return Invoices (適格返還請求書)
				for returns, discounts, and rebates.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("National Tax Agency - Consumption Tax"),
				URL:   "https://www.nta.go.jp/english/taxes/consumption_tax/01.htm",
			},
		},
		TimeZone: "Asia/Tokyo",
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
