// Package au provides the tax regime definition for Australia.
package au

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Australia.
const CountryCode = "AU"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("au", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax regime definition for Australia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.AUD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Australia",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Australia's Goods and Services Tax (GST) is a broad-based tax of 10%
				applied to most goods, services and other items sold or consumed in
				Australia. It is administered by the Australian Taxation Office (ATO).

				Businesses are identified by their Australian Business Number (ABN), an
				11-digit identifier whose validity is confirmed with a weighted
				modulus-89 check.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ATO - GST"),
				URL:   "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst",
			},
		},
		TimeZone:   "Australia/Sydney",
		Categories: taxCategories(),
	}
}
