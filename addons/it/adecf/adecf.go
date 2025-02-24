// Package adecf handles the validation rules in order to use
// GOBL with the Italian Agenzia delle Entrate format.
package adecf

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Key to identify the AdE adecf addon
const (
	// V1 for AdE format
	V1 cbc.Key = "it-adecf-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Italy AdE adecf v1.x",
		},
		Tags:       []*tax.TagSet{},
		Extensions: extensions,
		Validator:  validate,
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}
