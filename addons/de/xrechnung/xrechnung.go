// Package xrechnung provides extensions and validations for the German XRechnung standard version 3.0.2 for electronic invoicing.
package xrechnung

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
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
			i18n.EN: "German XRechnung 3.X",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the German XRechnung version 3.X standard for electronic invoicing.
				XRechnung is based on the European Norm (EN) 16931 and is mandatory for business-to-government
				(B2G) invoices in Germany. This addon provides the necessary structures and validations to
				ensure compliance with the XRechnung specifications.

				For more information on XRechnung, visit [www.xrechnung.de](https://www.xrechnung.de/).
			`),
		},
		Validator: validate,
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *pay.Instructions:
		return validatePaymentInstructions(obj)
	}
	return nil
}
