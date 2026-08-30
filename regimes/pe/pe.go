// Package pe provides the tax regime definition for Peru.
package pe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Peru.
const CountryCode = "PE"

// Identity type codes for non-tax identities commonly used in Peru
// (SUNAT Catalogue 06).
const (
	// IdentityTypeDNI represents the "Documento Nacional de Identidad",
	// the national identity document for Peruvian citizens.
	IdentityTypeDNI cbc.Code = "DNI"
	// IdentityTypeCE represents the "Carné de Extranjería", the identity
	// document issued to foreign residents in Peru.
	IdentityTypeCE cbc.Code = "CE"
	// IdentityTypePassport represents a passport, used to identify
	// non-resident individuals.
	IdentityTypePassport cbc.Code = "PAS"
)

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("pe", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
		orgIdentityRules(),
		billInvoiceRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
	norm.RegisterWithGuard(is.InContext(tax.RegimeIn(CountryCode)),
		norm.For(normalizeOrgIdentity),
	)
}

// New instantiates a new Peruvian tax regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.PEN,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Peru",
			i18n.ES: "Perú",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Peru applies the IGV (Impuesto General a las Ventas), a
				value-added tax administered by SUNAT, at a standard rate that
				combines the IGV itself with the Municipal Promotion Tax (IPM).
				A temporary reduced rate applies to qualifying micro and small
				enterprises in the restaurant, hotel and tourist accommodation
				sectors.

				Businesses and individuals are identified by their RUC
				(Registro Único de Contribuyentes), an eleven-digit number
				with a taxpayer type prefix and a mod-11 check digit.
			`),
		},
		TimeZone:   "America/Lima",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		// Only the "nota de crédito" (SUNAT Catalogue 09) and "nota de
		// débito" (Catalogue 10) exist as correction documents in Peru.
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
