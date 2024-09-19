// Package tbai provides the TicketBAI addon
package tbai

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for TicketBAI versions 1.x
	V1 cbc.Key = "es-tbai-v1"
)

// Official stamps or codes validated by government agencies
const (
	StampCode cbc.Key = "tbai-code"
	StampQR   cbc.Key = "tbai-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "TicketBAI",
		},
		Extensions:  extensions,
		Validator:   validate,
		Normalizer:  normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(_ any) {
	// nothing to normalize yet
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
