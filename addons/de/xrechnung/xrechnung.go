// Package xrechnung provides extensions and validations for the German XRechnung standard version 3.0.2 for electronic invoicing.
package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V3 is the key for the XRechnung version 3.x
	V3 cbc.Key = "de-xrechnung-v3"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V3,
		Name: i18n.String{
			i18n.EN: "German XRechnung 3.0.2",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
			Extensions to support the German XRechnung standard version 3.0.2 for electronic invoicing.
			XRechnung is based on the European Norm (EN) 16931 and is mandatory for business-to-government
			(B2G) invoices in Germany. This addon provides the necessary structures and validations to
			ensure compliance with the XRechnung format.

			For more information on XRechnung, visit:
			https://www.xrechnung.de/
			`),
		},
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Scenarios:  scenarios,
		Extensions: extensions,
		Normalizer: normalize,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Combo:
		NormalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return ValidateInvoice(obj)
	case *tax.Combo:
		return ValidateTaxCombo(obj)
	}
	return nil
}
