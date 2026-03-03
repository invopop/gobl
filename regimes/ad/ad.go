// Package ad provides the Andorran tax regime.
package ad

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

func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "AD",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Andorra",
			i18n.CA: "Andorra",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Andorra's tax system is administered by the Departament de Tributs i Fronteres.
				The primary indirect tax is the Impost General Indirecte (IGI), which functions
				as a value-added tax on goods and services.

				Entities and individuals are identified by a Número de Registre Tributari (NRT),
				which consists of 8 characters: a letter prefix indicating entity type, six digits,
				and a control letter.

				Electronic invoicing follows specific form structures: Form 980-A for traveler/export
				refunds, Form 980-B for non-resident B2B transactions (requiring a local fiscal
				representative), and Form 980-C for diplomatic exemptions under Article 15 of
				Law 11/2012.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Departament de Tributs i Fronteres"),
				URL:   "https://www.govern.ad/tributs-i-fronteres",
			},
		},
		TimeZone:   "Europe/Andorra",
		Scenarios:  scenarios,
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: categories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types:  []cbc.Key{bill.InvoiceTypeCreditNote},
			},
		},
	}
}

func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Party:
		return validateParty(obj)
	}
	return nil
}

func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
