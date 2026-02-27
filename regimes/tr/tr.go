// Package tr provides the tax regime definition for Türkiye.
package tr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax regime definition for TR.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "TR",
		Currency:  currency.TRY,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Türkiye",
			i18n.TR: "Türkiye",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
			GIB (Revenue Administration) is the tax authority. Businesses use a
			10-digit tax number (VKN) as their tax identity. Individuals and sole
			traders use their 11-digit national identity number (TCKN) which serves
			as a recognized tax identifier domestically. For international trade,
			a VKN is required. Either a VKN or TCKN must be provided on invoices.

			E-invoicing is mandatory above a revenue threshold. B2B between
			registered parties uses e-Fatura; B2C and unregistered buyers use e-Arşiv.
			Both use the UBL-TR XML format.

			Credit and debit notes are not valid; adjustments require a corrective
			invoice that references the original.

			Reverse charge applies only when a Turkish buyer pays a foreign supplier,
			not to invoices issued by Turkish sellers. Exports are VAT-exempt.

			Withholding tax (Stopaj) and special consumption tax (ÖTV) are out of scope;
			their rates are product-specific and change frequently.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Revenue Administration",
					i18n.TR: "Gelir İdaresi Başkanlığı",
				},
				URL: "https://www.gib.gov.tr",
			},
			{
				Title: i18n.String{
					i18n.EN: "Tax Guide - Invest in Türkiye",
					i18n.TR: "Vergi Rehberi - Türkiye'ye Yatırım",
				},
				URL: "https://www.invest.gov.tr/en/investmentguide/pages/tax-guide.aspx",
			},
		},
		TimeZone: "Europe/Istanbul",
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					// Credit and debit notes are not valid certifying documents under
					// Tax Procedure Law (TPL) Article 227. Adjustments are handled via
					// corrective invoices that reference the original document.
					bill.InvoiceTypeCorrective,
				},
			},
		},
		Identities: identityDefinitions,
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
