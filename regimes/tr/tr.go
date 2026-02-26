// Package tr provides the tax regime definition for Türkiye.
//
// Additional context:
//   - e-Invoicing (e-Fatura for B2B, e-Arşiv for B2C) uses UBL-TR via GIB;
//     format-level concerns are deferred to a future addons/tr/efatura addon.
//   - No reverse charge: only applies inbound (Turkish buyer paying
//     non-resident); exports are KDV-exempt.
//   - Stopaj (withholding tax) and ÖTV (special consumption tax) are not
//     implemented (product-specific, change frequently).
//
// See also:
//   - https://www.fonoa.com/resources/blog/practical-guide-to-turkish-e-invoicing
//   - https://edicomgroup.com/electronic-invoicing/turkey
//   - https://taxsummaries.pwc.com/turkey/corporate/other-taxes

package tr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
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
		TimeZone: "Europe/Istanbul",
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
				// Credit and debit notes are not valid certifying documents under
				// Tax Procedure Law (TPL) Article 227. Adjustments are handled via
				// corrective invoices that reference the original document.
				// https://www.grc-legal.com/en/credit-debit-note-under-turkish-legislation/
					bill.InvoiceTypeCorrective,
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
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
