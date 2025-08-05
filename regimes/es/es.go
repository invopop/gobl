// Package es provides tax regime support for Spain.
package es

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

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryIRPF cbc.Code = "IRPF"
	TaxCategoryIGIC cbc.Code = "IGIC"
	TaxCategoryIPSI cbc.Code = "IPSI"
)

// Specific tax rate codes.
const (
	// IRPF non-standard Rates (usually for self-employed)
	TaxRatePro                cbc.Key = "pro"                 // Professional Services
	TaxRateProStart           cbc.Key = "pro-start"           // Professionals, first 2 years
	TaxRateModules            cbc.Key = "modules"             // Module system
	TaxRateAgriculture        cbc.Key = "agriculture"         // Agricultural
	TaxRateAgricultureSpecial cbc.Key = "agriculture-special" // Agricultural special
	TaxRateCapital            cbc.Key = "capital"             // Rental or Interest

	// Special tax rate surcharge extension
	TaxRateEquivalence cbc.Key = "eqs"
)

// Inbox key and role definitions.
// TODO: move to their own addon.
const (
	InboxKeyFACE cbc.Key = "face"

	// Main roles defined in FACE
	InboxRoleFiscal    cbc.Key = "fiscal"    // Fiscal / 01
	InboxRoleRecipient cbc.Key = "recipient" // Receptor / 02
	InboxRolePayer     cbc.Key = "payer"     // Pagador / 03
	InboxRoleCustomer  cbc.Key = "customer"  // Comprador / 04

)

// New provides the Spanish tax regime definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "ES",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Spain",
			i18n.ES: "España",
		},
		TimeZone: "Europe/Madrid",
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Identities: identityDefinitions,
		Categories: taxCategories,
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: correctionDefinitions,
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

// Normalize will perform any regime specific normalizations on the data.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
