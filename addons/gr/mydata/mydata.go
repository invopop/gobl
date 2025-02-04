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
		Sources: []*cbc.Source{
			{
				Title:       i18n.NewString("myDATA API Documentation v1.0.7"),
				URL:         "https://www.aade.gr/sites/default/files/2023-10/myDATA%20API%20Documentation_v1.0.7_eng.pdf",
				ContentType: "application/pdf",
			},
			{
				Title:       i18n.NewString("Greek Peppol BIS Billing 3.0"),
				URL:         "https://www.gsis.gr/sites/default/files/eInvoice/Instructions%20to%20B2G%20Suppliers%20and%20certified%20PEPPOL%20Providers%20for%20the%20Greek%20PEPPOL%20BIS-EN-%20v1.0.pdf",
				ContentType: "application/pdf",
			},
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
