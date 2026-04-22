// Package facturae provides the FacturaE addon for Spanish invoices.
package facturae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the FacturaE addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "es-facturae"

	// V3 for FacturaE versions 3.x
	V3 cbc.Key = Key + "-v3"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("ES-FACTURAE"),
		is.InContext(tax.AddonIn(V3)),
		billInvoiceRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V3,
		Name: i18n.String{
			i18n.EN: "Spain FacturaE",
		},
		Extensions:  extensions,
		Normalizer:  normalize,
		Scenarios:   scenarios,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	}
}
