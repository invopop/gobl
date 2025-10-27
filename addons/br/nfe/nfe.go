// Package nfe handles extensions and validation rules to issue NF-e in
// Brazil.
package nfe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

const (
	// V4 is the key for the NF-e 4.00 layout
	V4 cbc.Key = "br-nfe-v4"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V4,
		Name: i18n.String{
			i18n.EN: "Brazil NF-e 4.00",
		},
		Validator:  validate,
		Normalizer: normalize,
		Extensions: extensions,
		Scenarios:  scenarios,
		Identities: identities,
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *bill.Line:
		return validateLine(obj)
	case *pay.Instructions:
		return validatePayInstructions(obj)
	case *pay.Advance:
		return validatePayAdvance(obj)
	}
	return nil
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Advance:
		normalizePayAdvance(obj)
	}
}
