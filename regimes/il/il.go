// Package il provides the tax region definition for Israel.
package il

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

// New provides the tax region definition for IL.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "IL",
		Currency:  currency.ILS,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Israel",
			i18n.HE: "ישראל",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Israel's value-added tax, known as Ma'am (מע"מ), is administered by the Israel Tax Authority (ITA) under the VAT Law 5736-1975.

        		The general VAT rate is 18%, raised from 17% on 1 January 2025. A zero rate applies to exports, services provided to non-residents, supplies to the Eilat Free Trade Zone, and fresh fruits and vegetables. Certain financial services, residential real estate leases, educational and healthcare services, and non-profit organisations are exempt from VAT.

        		Businesses are classified into three categories based on annual turnover: Osek Murshe (Authorized Dealer) for businesses above NIS 120,000 that must register, charge VAT, and file bimonthly returns; Osek Patur (Exempt Dealer) for businesses below the threshold that cannot charge or recover VAT; and Osek Zair (Small Dealer), introduced in 2025 as a simplified alternative for sole operators below the threshold.

        		VAT-registered businesses are identified by a 9-digit Mispar Osek Murshe. For sole proprietors, this is typically the same as the personal Mispar Zehut. For companies and other legal entities, it corresponds to the Corporations Authority registration number, where prefixes indicate entity type: 51 (companies), 58 (associations), among others.

        		Since May 2024, B2B invoices above certain thresholds must be pre-cleared with the ITA via the SHAAM platform, which assigns an Allocation Number required for input VAT deduction. The threshold is being reduced in phases from NIS 25,000 to NIS 5,000 by June 2026. SHAAM integration is planned as a future GOBL addon.

        		Simplified tax invoices are permitted when the recipient is not a registered dealer, per Section 46 of the VAT Law.
			`),
		},
		TimeZone: "Asia/Jerusalem",
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
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
