// Package saft provides the SAF-T addon for Portuguese invoices.
package saft

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for SAF-T (PT) versions 1.x
	V1 cbc.Key = "pt-saft-v1"
)

func init() {
	tax.RegisterAddon(newAddon())
}

func newAddon() *tax.Addon {
	return &tax.Addon{
		Key:        V1,
		Extensions: extensions,
		Normalizer: normalize,
		Scenarios:  scenarios,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}
