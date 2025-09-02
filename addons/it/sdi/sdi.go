// Package sdi handles the extensions and validation rules in order to use
// GOBL with the Italian SDI and FatturaPA format.
package sdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for SDI's FatturaPA verions 1.x
	V1 cbc.Key = "it-sdi-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Italy SDI FatturaPA v1.x",
		},
		Extensions: extensions,
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Inboxes:    inboxes,
		Normalizer: normalize,
		Scenarios:  scenarios,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Advance:
		normalizePayAdvance(obj)
	case *org.Address:
		normalizeAddress(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *pay.Instructions:
		return validatePayInstructions(obj)
	case *pay.Advance:
		return validatePayAdvance(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}
