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
	tax.RegisterAddon(&addon{})
}

type addon struct {
	tax.BaseAddon
}

func (addon) Key() cbc.Key {
	return KeyV1
}

func (addon) Extensions() []*cbc.KeyDefinition {
	return extensions
}

func (addon) Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}

func (addon) Corrections() tax.CorrectionSet {
	return invoiceCorrectionDefinitions
}
