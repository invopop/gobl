// Package mydata handles the extensions and validation rules in order to use
// GOBL with the Greek MyData format.
package mydata

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for Greece MyData XML v1.x
	V1 cbc.Key = "gr-mydata-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Greece MyData v1.x",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Greece uses the myDATA and Peppol BIS Billing 3.0 formats for their e-invoicing/tax-reporting system.
				This addon will ensure that the GOBL documents have all the required fields to be able to correctly
				generate the myDATA XML reporting files.
			`),
		},
		Extensions: extensions,
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Normalizer: normalize,
		Scenarios:  scenarios,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Advance:
		normalizePayAdvance(obj)
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
