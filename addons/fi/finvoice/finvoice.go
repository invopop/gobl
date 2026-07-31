// Package finvoice provides validations for the Finnish Finvoice 3.0 format.
package finvoice

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the Finvoice addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "fi-finvoice"

	// V3 is the key for the Finvoice version 3.0
	V3 cbc.Key = Key + "-v3"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("FI-FINVOICE"),
		is.InContext(tax.AddonIn(V3)),
		billInvoiceRules(),
	)
	norm.RegisterWithGuard(
		is.InContext(tax.AddonIn(V3)),
		norm.For(normalizePayInstructions),
		norm.For(normalizePayCreditTransfer),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V3,
		Name: i18n.String{
			i18n.EN: "Finland Finvoice 3.0",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the Finnish Finvoice 3.0 format for electronic invoicing.
				Finvoice conforms to the European Norm (EN) 16931, so this addon only adds
				the Finvoice-specific requirements on top of the EN 16931 rules: a named
				customer, and the payment data needed to build the Finvoice EpiDetails
				payment block (credit transfer instructions with an IBAN, a payment
				reference, and a due date), which is mandatory on every invoice.

				For more information on Finvoice, visit
				[www.finanssiala.fi](https://www.finanssiala.fi/en/topics/finvoice-implementation-guidelines/).
			`),
		},
	}
}
