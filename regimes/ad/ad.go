// Package ad provides the tax regime data for Andorra.
package ad

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.AD.Tax(),
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Andorra",
			i18n.CA: "Andorra",
			i18n.ES: "Andorra",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The main indirect tax in Andorra is the 'Impost General Indirecte (IGI)', and it is enforced since 1st of January 2013.
				The NRT (Número de Registre Tributari) is the tax identification number for companies in Andorra. It has the following format: 'X-999999-X'
				- A leading letter (identifying the type of person/entity):
				  - F: Individual Residents
				  - E: Non-resident Individuals
				  - L: Limited Liability Companies (S.L.)
				  - A: Joint-stock Corporations (S.A.)
				- Six digits.
				- A trailing control letter.

				Invoices allow corrections through credit notes (Nota d'Abonament) and debit notes (Nota de Càrrec).
				
				Invoice presentation requirements are:
				- Quarterly: Companies with a turnover of more than €250,000 (April, July, October, January).
				- Semestral: Companies with a turnover of less than €250,000 (July, January).
				- Start of activity: Generally declared semestrally (July and January), unless the special regime applies.

				Sources:
				- [Departament de Tributs i de Fronteres - Andorra](https://www.impostos.ad)
				- [Andorra NRT number guide](https://lookuptax.com/docs/tax-identification-number/andorra-tax-id-guide)
			`),
		},
		TimeZone:   "Europe/Andorra",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote, // CAT: Nota d'Abonament
					bill.InvoiceTypeDebitNote,  // CAT: Nota de Càrrec
				},
			},
		},
		Categories: taxCategories(),
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
