// Package nfe handles extensions and validation rules to issue NF-e in
// Brazil.
package nfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// V4 is the key for the NF-e 4.00 layout
	V4 cbc.Key = "br-nfe-v4"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		V4.String(),
		rules.GOBL.Add("BR-NFE-V4"),
		is.HasContext(tax.AddonIn(V4)),
		billInvoiceRules(),
		billLineRules(),
		payInstructionsRules(),
		payAdvanceRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V4,
		Name: i18n.String{
			i18n.EN: "Brazil NF-e 4.00",
		},
		Normalizer: normalize,
		Extensions: extensions,
		Scenarios:  scenarios,
		Identities: identities,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Advance:
		normalizePayAdvance(obj)
	}
}
