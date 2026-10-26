// Package tbai provides the TicketBAI addon
package tbai

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the TicketBAI addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "es-tbai"

	// V1 for TicketBAI versions 1.x
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("ES-TBAI"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		taxComboRules(),
	)
	norm.RegisterWithGuard(
		is.InContext(tax.AddonIn(V1)),
		norm.For(normalizeInvoice),
		norm.For(normalizeBillLine),
		norm.For(normalizeTaxCombo),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Spain TicketBAI",
		},
		Extensions:  extensions,
		Stamps:      stamps,
		Corrections: invoiceCorrectionDefinitions,
	}
}
