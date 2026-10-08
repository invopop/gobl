// Package bs provides the tax regime definition for Bahamas.
package bs

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Bahamas.
const CountryCode = "BS"

func init() {
	tax.RegisterRegimeDef(New()) // regimes.go
	rules.Register("bs", rules.GOBL.Add(CountryCode), // rules.go
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax regime definition for New Zealand.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.NZD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Commonwealth of The Bahamas",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				New Zealand's Goods and Services Tax (GST) is a broad-based tax
				applied to most goods and services supplied in New Zealand. It is
				administered by Inland Revenue (IRD).

				Businesses are identified by their IRD number, an 8 or 9 digit
				identifier whose validity is confirmed with a weighted modulus-11 check.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Inland Revenue - GST",
					i18n.MI: "Te Tari Taake - GST",
				},
				URL: "https://www.ird.govt.nz/gst",
			},
		},
		TimeZone:   "Pacific/Auckland",
		Categories: taxCategories(),
	}
}
