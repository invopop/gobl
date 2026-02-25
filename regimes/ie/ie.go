// Package ie provides a regime definition for Ireland.
package ie

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Irish regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.IE.Tax(),
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Ireland",
			i18n.GA: "Éire",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Ireland's tax system is administered by the Revenue Commissioners (Na
				Coimisinéirí Ioncaim). As an EU member state, Ireland follows the EU VAT
				Directive with locally adapted rates.

				VAT rates include a 23% standard rate for most goods and services, a 13.5%
				reduced rate for tourism, hospitality, and certain construction services,
				a 9% second reduced rate for newspapers, e-publications, and sports
				facilities, a 4.8% rate for livestock, and a 0% zero rate for food,
				children's clothing, oral medicines, and exports.

				Businesses are identified by their VAT registration number in the format IE
				followed by 7 digits and 1-2 letters. The threshold for mandatory VAT
				registration is EUR 80,000 for goods and EUR 40,000 for services.

				Ireland supports credit notes for invoice corrections. E-invoicing via
				PEPPOL is supported for B2G transactions.
			`),
		},
		TimeZone:   "Europe/Dublin",
		Categories: taxCategories,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Validator:  Validate,
		Normalizer: Normalize,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
