// Package dian provides the DIAN UBL 2.1 extensions used in Colombia.
package dian

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V2 for DIAN UBL 2.1 in Colombia
	V2 cbc.Key = "co-dian-v2"
)

// DIAN official codes to include in stamps.
const (
	StampCUDE cbc.Key = "dian-cude"
	StampQR   cbc.Key = "dian-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V2,
		Name: i18n.String{
			i18n.EN: "Colombia DIAN UBL 2.X",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Extensions to support the Colombian DIAN (Dirección de Impuestos y Aduanas Nacionales)
				specifications for electronic invoicing based on UBL 2.1.
			`),
		},
		Extensions:  extensions,
		Identities:  identities,
		Normalizer:  normalize,
		Validator:   validate,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(_ any) {
	// no normalizations yet
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
