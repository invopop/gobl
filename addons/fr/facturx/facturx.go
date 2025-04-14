// Package facturx provides extensions and validations for the French Factur-X standard version 1.07 for electronic invoicing.
package facturx

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 is the key for the Factur-X version 1.x
	V1 cbc.Key = "fr-facturx-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "French Factur-X v1",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French Factur-X 1.x standard for electronic invoicing.
				Factur-X is based on the European Norm (EN) 16931 and is mandatory for business-to-government
				(B2G) invoices in France. This addon provides the necessary structures and validations to
				ensure compliance with the Factur-X specifications.

				For more information on Factur-X, visit [fnfe-mpe.org](https://fnfe-mpe.org/factur-x/factur-x_en/#).
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
