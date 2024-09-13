// Package facturae provides the FacturaE addon for Spanish invoices.
package facturae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

const (
	// KeyV3 for FacturaE versions 3.x
	KeyV3 cbc.Key = "es-facturae-v3"
)

func init() {
	tax.RegisterAddon(newAddon())
}

func newAddon() *tax.Addon {
	return &tax.Addon{
		Key:         KeyV3,
		Extensions:  extensions,
		Normalizer:  normalize,
		Scenarios:   scenarios,
		Validator:   validate,
		Corrections: invoiceCorrectionDefinitions,
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
