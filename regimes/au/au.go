// Package au provides the tax region definition for Australia.
package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for Australia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "AU",
		Currency:  currency.AUD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Australia",
		},
		TimeZone: "Australia/Sydney",
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Australia operates a Goods and Services Tax (GST) system administered by the
				Australian Taxation Office (ATO). GST is a broad-based tax of 10% on most goods,
				services and other items sold or consumed in Australia.

				Key features:
				- Standard GST rate: 10%
				- GST-free (zero-rated) supplies: Basic food, most health services, medical aids,
				  educational courses, childcare, exports
				- Input-taxed (exempt) supplies: Financial supplies, residential rent, residential
				  premises sales
				- ABN (Australian Business Number): 11-digit identifier for businesses

				Tax invoices with a taxable value of A$1,000 or more must include either the
				buyer's name or ABN.
			`),
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories(),
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Party:
		normalizeParty(obj)
	}
}

// Stub functions that will be implemented in subsequent commits
func validateInvoice(inv *bill.Invoice) error {
	return nil
}

func validateTaxIdentity(tID *tax.Identity) error {
	return nil
}

func normalizeParty(party *org.Party) {
	// Will be implemented with tax identity normalization
}
