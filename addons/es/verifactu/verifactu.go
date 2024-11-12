// Package verifactu provides the Verifactu addon
package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for Verifactu versions 1.x
	V1 cbc.Key = "es-verifactu-v1"
)

// Official stamps or codes validated by government agencies
const (
	// StampQR contains the URL included in the QR code.
	StampQR cbc.Key = "verifactu-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Spain Verifactu",
		},
		Extensions:  extensions,
		Validator:   validate,
		Scenarios:   scenarios,
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
