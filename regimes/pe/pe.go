// Package pe provides the tax regime definition for Peru.
package pe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
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
		TimeZone: "America/Lima",
	}
}
