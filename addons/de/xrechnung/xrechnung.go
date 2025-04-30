// Package xrechnung provides extensions and validations for the German XRechnung standard version 3.0.2 for electronic invoicing.
package xrechnung

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// BR-DE-14 - set percent in convertor as this rule requires it, even if it is 0. BT-119
// BR-DE-16 - covered by tax ID being required. BT-31
// BR-DE-18 - references format of payment.terms.details. BT-20
// BR-DE-20 - partialy covered by gobl validation of IBAN format. BT-84
// BR-DE-21 - look at BT-24 mapping of gobl. BT-24
// BR-DE-22 - refers to attachments. BG-24
// BR-DE-27 - handled by gobl validation of phone number. BT-42
// BR-DE-28 - handled by gobl validation of email address. BT-43

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
	}
	return nil
}
