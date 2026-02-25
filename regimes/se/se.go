// Package se provides a regime definition for Sweden.
package se

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Swedish regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.SE.Tax(),
		Currency:  currency.SEK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Sweden",
			i18n.SE: "Sverige",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Sweden's tax system is administered by the Swedish Tax Agency
				(Skatteverket). As an EU member state, Sweden follows the EU VAT Directive
				with locally adapted rates.

				Moms (Mervärdesskatt) rates include a 25% standard rate for most goods and
				services, a 12% reduced rate for food, hotel accommodation, restaurant
				services, and certain repair services, and a 6% heavily reduced rate for
				passenger transport, books, newspapers, cultural events, and intellectual
				property. Exports and certain financial and healthcare services are exempt.

				Businesses are identified by their Organisationsnummer (organization number),
				a 10-digit number validated with the Luhn algorithm. The Swedish VAT number
				uses the format SE followed by the 10-digit organization number plus "01" as
				check digits. Individuals may be identified by their Personnummer (personal
				identity number, format YYMMDD-XXXX) or Samordningsnummer (coordination
				number for non-residents, where the day component is offset by 60). Sole
				proprietorships use the owner's personal number as their organization number.

				E-invoicing via PEPPOL BIS Billing 3.0 is mandatory for all B2G transactions
				since April 2019. F-tax (F-skatt) registration indicates that a business
				handles its own tax payments, exempting customers from withholding obligations.
				Reverse charge applies in specific sectors (construction, metals, waste) and
				for cross-border transactions within the EU. Simplified invoices may be used
				for transactions up to 4000 SEK.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("PEPPOL BIS Billing 3.0"),
				URL:   "https://docs.peppol.eu/poacc/billing/3.0/",
			},
			{
				Title: i18n.NewString("Skatteverket - VAT Rates"),
				URL:   "https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momssatspavarorochtjanster.4.58d555751259e4d66168000409.html",
			},
			{
				Title: i18n.NewString("Skatteverket - Invoice Requirements"),
				URL:   "https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momslagensregleromfakturering.4.58d555751259e4d66168000403.html",
			},
		},
		TimeZone:   "Europe/Stockholm",
		Identities: identityTypeDefinitions,
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
	case *org.Identity:
		return validateOrgIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
