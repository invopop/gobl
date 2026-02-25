// Package fr provides the tax region definition for France.
package fr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Identification keys used for additional codes not
// covered by the standard fields.
const (
	IdentityTypeSIREN cbc.Code = "SIREN" // SIREN is the main local tax code used in france, we use the normalized VAT version for the tax ID.
	IdentityTypeSIRET cbc.Code = "SIRET" // SIRET number combines the SIREN with a branch number.
	IdentityTypeRCS   cbc.Code = "RCS"   // Trade and Companies Register.
	IdentityTypeRM    cbc.Code = "RM"    // Directory of Traders.
	IdentityTypeNAF   cbc.Code = "NAF"   // Identifies the main branch of activity of the company or self-employed person.
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "FR",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "France",
			i18n.FR: "La France",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				France's tax system is administered by the Direction Générale des Finances
				Publiques (DGFiP). As an EU member state, France follows the EU VAT Directive
				with locally adapted rates.

				TVA (Taxe sur la Valeur Ajoutée) applies at standard, intermediate, reduced,
				and super-reduced rates covering various categories of goods and services.

				Businesses are identified by three closely related numbers: the VAT code
				(numéro de TVA intracommunautaire), an 11-digit number starting with a
				2-digit checksum followed by the 9-digit SIREN; the SIREN itself, a 9-digit
				company identifier from the national register (Répertoire SIRENE); and the
				SIRET, which extends the SIREN with a 5-digit establishment number to form
				a 14-digit code.

				France supports both corrective invoices and credit notes for invoice
				corrections. E-invoicing via the Chorus Pro platform is mandatory for B2G
				transactions, with B2B e-invoicing being progressively mandated through the
				CTC (Continuous Transaction Controls) framework.
			`),
		},
		TimeZone: "Europe/Paris",
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				// France supports both corrective methods
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote, // Code 381
					bill.InvoiceTypeCorrective, // Code 384
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
	case *org.Identity:
		return validateIdentity(obj)
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
