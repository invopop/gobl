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
	// StampCode contains the code required to be presented alongside
	// the QR code.
	StampCode cbc.Key = "tbai-code"
	// StampQR contains the URL included in the QR code.
	StampQR cbc.Key = "tbai-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Spain TicketBAI",
		},
		Extensions:  extensions,
		Validator:   validate,
		Normalizer:  normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	// nothing to normalize yet
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
