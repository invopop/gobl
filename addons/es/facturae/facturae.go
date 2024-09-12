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
	tax.RegisterAddon(&addon{})
}

type addon struct {
	tax.BaseAddon
}

func (addon) Key() cbc.Key {
	return KeyV3
}

func (addon) Extensions() []*cbc.KeyDefinition {
	return extensions
}

func (addon) Normalize(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	}
	return nil
}

func (addon) Scenarios() []*tax.ScenarioSet {
	return scenarios
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
