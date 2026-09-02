// Package jp provides the tax region definition for Japan.
package jp

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Japan.
const CountryCode = "JP"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("jp", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax region definition for Japan.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.JPY,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Japan",
			i18n.JA: "日本",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("National Tax Agency - Information about Consumption Tax"),
				URL:   "https://www.nta.go.jp/english/taxes/consumption_tax/index.htm",
			},
			{
				Title: i18n.NewString("National Tax Agency - Outline of the invoice system"),
				URL:   "https://www.nta.go.jp/english/taxes/consumption_tax/pdf/2022/general_17.pdf",
			},
		},
		TimeZone: "Asia/Tokyo",
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Japan levies a Consumption Tax (消費税, shōhizei) administered by the
				National Tax Agency (NTA). It is a value-added tax charged at a standard
				rate of 10% with a reduced rate of 8% for food and drink (excluding alcohol
				and dining out) and certain newspaper subscriptions. Exports and export-like
				transactions are zero-rated.

				Although the rate legally decomposes into a national Consumption Tax portion
				and a Local Consumption Tax portion (a fixed fraction of the national tax),
				the two are invoiced as a single combined rate. GOBL therefore models a single
				Consumption Tax category using the shared VAT definition.

				Since October 2023, Japan operates the Qualified Invoice System
				(適格請求書等保存方式, "Invoice Seido"). Only a registered "qualified invoice
				issuer" may issue invoices that allow a registered buyer to claim an input tax
				credit. Such issuers are identified by a registration number consisting of the
				letter "T" followed by 13 digits.
			`),
		},
		Categories: taxCategories(),
	}
}
