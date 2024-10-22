// Package gb provides the United Kingdom tax regime.
package gb

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Identification code types unique to the United Kingdom.
const (
	IdentityTypeCRN cbc.Code = "CRN" // Company Registration Number
)

var (
	altCountryCodes = []l10n.Code{
		l10n.XI, // Northern Ireland
		l10n.XU, // UK except Northern Ireland (Brexit)
	}
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:         "GB",
		AltCountryCodes: altCountryCodes,
		Currency:        currency.GBP,
		Name: i18n.String{
			i18n.EN: "United Kingdom",
		},
		TimeZone:   "Europe/London",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			common.InvoiceScenarios(),
		},
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
		Categories:    taxCategories,
		IdentityTypes: identityTypeDefinitions,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
	}
}

// Validate checks the document type and determines if it can be validated. Note that in
// the GB tax regime we don't need to validate the presence of the supplier's tax ID.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj, altCountryCodes...)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
