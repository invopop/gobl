// Package is provides the tax regime definition for Iceland.
package is

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	rulesis "github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Iceland.
const CountryCode = "IS"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("is", rules.GOBL.Add(CountryCode), taxIdentityRules(), orgIdentityRules())
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
	norm.RegisterWithGuard(rulesis.InContext(tax.RegimeIn(CountryCode)),
		norm.For(normalizeOrgIdentity),
	)
}

// New instantiates a new Icelandic tax regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.IS.Tax(),
		Currency:  currency.ISK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Iceland",
			i18n.IS: "Ísland",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Iceland's tax system is administered by the Directorate of Internal
				Revenue (Skatturinn). Although Iceland is not an EU member state, it is
				part of the European Economic Area (EEA) and operates a value added tax
				broadly aligned with European practice.

				VSK (Virðisaukaskattur) applies at a standard and a reduced rate. 
				Exports are zero-rated. Health, education, public transport, real 
				estate leasing, sports and financial services are exempt under Article 
				12 of the VAT Act.

				Businesses and individuals alike are identified by their kennitala, a
				10-digit national identification number issued by Registers Iceland
				(Þjóðskrá Íslands). The same number serves as the VAT identifier; for
				international use it is prefixed with "IS". The kennitala encodes a date
				(birth date for individuals, registration date for companies) followed by 
				two random digits, a MOD-11 checksum digit, and a final century indicator.

				E-invoicing is mandatory for all B2G transactions since 2020 under
				Regulation 44/2019, using the Peppol BIS Billing 3.0 format together with
				the national CIUS TS-236. The Peppol Authority is Fjársýsla ríkisins
				(FJS). B2B and B2C e-invoicing is voluntary but widely adopted.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Registers Iceland - ID numbers (kennitala)"),
				URL:   "https://www.skra.is/english/people/my-registration/id-numbers/",
			},
			{
				Title: i18n.NewString("Skatturinn - Value Added Tax (VSK)"),
				URL:   "https://www.skatturinn.is/english/companies/value-added-tax/",
			},
			{
				Title: i18n.NewString("European Commission - eInvoicing in Iceland"),
				URL:   "https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/eInvoicing+in+Iceland",
			},
			{
				Title: i18n.NewString("Ecosio - E-invoicing in Iceland"),
				URL:   "https://ecosio.com/en/compliance/iceland/e-invoicing/",
			},
		},
		TimeZone:   "Atlantic/Reykjavik",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
	}
}
