// Package zatca provides extensions and validations for the Saudi Arabia
// ZATCA (Zakat, Tax and Customs Authority) e-invoicing requirements.
package zatca

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 is the key for the ZATCA e-invoicing addon version 1.
	V1 cbc.Key = "sa-zatca-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		V1.String(),
		rules.GOBL.Add("SA-ZATCA-V1"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		taxComboRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Saudi Arabia ZATCA",
			i18n.AR: "هيئة الزكاة والضريبة والجمارك",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the Saudi Arabia ZATCA (Zakat, Tax and Customs Authority) e-invoicing
				requirements based on UBL 2.1 with EN 16931 as an intermediate layer and KSA-specific
				extensions (BR-KSA-* rules).

				ZATCA e-invoicing covers both standard tax invoices (B2B/B2G) sent for clearance
				and simplified tax invoices (B2C) sent for reporting through the FATOORA platform.

				This addon extends EN 16931 with Saudi-specific fields and validations including
				invoice type transactions, address requirements, and supply date handling.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ZATCA E-Invoicing Developer Portal"),
				URL:   "https://zatca.gov.sa/en/E-Invoicing/SystemsDevelopers/Pages/E-Invoice-specifications.aspx",
			},
		},
		Extensions: extensions,
		Scenarios:  scenarios,
		Normalizer: normalize,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	}
}
