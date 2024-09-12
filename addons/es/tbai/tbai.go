// Package tbai provides the TicketBAI addon
package tbai

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

const (
	// KeyV1 for TicketBAI versions 1.x
	KeyV1 cbc.Key = "es-tbai-v1"
)

// Official stamps or codes validated by government agencies
const (
	StampCode cbc.Key = "tbai-code"
	StampQR   cbc.Key = "tbai-qr"
)

func init() {
	tax.RegisterAddon(newAddon())
}

func newAddon() *tax.Addon {
	return &tax.Addon{
		Key:         KeyV1,
		Extensions:  extensions,
		Validate:    validate,
		Normalize:   normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(_ any) error {
	return nil
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
