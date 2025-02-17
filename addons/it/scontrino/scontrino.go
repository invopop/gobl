// Package scontrino handles the validation rules in order to use
// GOBL with the Italian Agenzia delle Entrate format.
package scontrino

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Key to identify the AdE scontrino addon
const (
	// V1 for AdE format
	V1 cbc.Key = "it-scontrino-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Italy AdE scontrino v1.x",
		},
		Tags:       []*tax.TagSet{},
		Normalizer: normalize,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
