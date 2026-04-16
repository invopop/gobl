// Package fi provides tax regime support for Finland.
package fi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-1 alpha-2 code for Finland.
const CountryCode = "FI"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("fi", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
}

// New instantiates a new Finland regime
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Finland",
			i18n.FI: "Suomi",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Finland's tax system is administered by the Finnish Tax Administration
				(Verohallinto). As an EU member state, Finland follows the EU VAT
				Directive.

				VAT (Arvonlisävero, ALV) applies at a standard rate of 25.5%, a reduced
				rate of 13.5% on food, restaurants, books, transport, accommodation, and
				cultural events, and a further reduced rate of 10% on newspapers and
				magazines. The 13.5% rate took effect in 2026, replacing the previous 14%
				tier and absorbing several categories formerly at 10%. Exports and
				intra-EU sales to VAT-liable buyers are zero-rated.

				Businesses are identified by their Business ID (Y-tunnus), a 7-digit
				number plus a check digit (format 1234567-8). The Finnish VAT number
				is formed by prefixing FI and removing the hyphen (e.g. FI12345678).

				Invoice corrections are not restricted to specific document types; any
				corrective document referencing the original invoice is accepted.

				E-invoicing is mandatory for B2G transactions since April 2021 under
				Act 241/2019, which implements EU Directive 2014/55/EU. Finvoice and
				TEAPPSXML are the primary domestic formats; PEPPOL BIS is also supported.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Verohallinto - VAT Invoice Requirements"),
				URL:   "https://www.vero.fi/en/detailed-guidance/guidance/48090/vat-invoice-requirements/",
			},
			{
				Title: i18n.NewString("Verohallinto - Business ID Validation"),
				URL:   "https://www.vero.fi/globalassets/tietoa-verohallinnosta/ohjelmistokehittajille/yritys--ja-yhteisötunnuksen-ja-henkilötunnuksen-tarkistusmerkin-tarkistuslaskenta.pdf",
			},
			{
				Title: i18n.NewString("Valtiokonttori - Invoicing the State"),
				URL:   "https://www.valtiokonttori.fi/en/services/government-e-invoices/invoicing-the-state/",
			},
			{
				Title: i18n.NewString("Finlex - Act on Electronic Invoicing 241/2019"),
				URL:   "https://www.finlex.fi/fi/laki/alkup/2019/20190241",
			},
		},
		TimeZone:   "Europe/Helsinki",
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},
		Normalizer: Normalize,
	}
}

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
