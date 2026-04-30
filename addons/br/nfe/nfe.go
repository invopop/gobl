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
	// Key identifies the NF-e addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "br-nfe"

	// V4 is the key for the NF-e 4.00 layout
	V4 cbc.Key = Key + "-v4"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("BR-NFE"),
		is.InContext(tax.AddonIn(V4)),
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
	case *pay.Record:
		normalizePayRecord(obj)
	}
}
