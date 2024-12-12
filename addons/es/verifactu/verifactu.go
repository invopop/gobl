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

	// QRPrefix is the compulsory text to appear above the QR code.
	QRPrefix string = "QR tributario:"
	// QRSuffix is the compulsory text to appear below the QR code.
	QRSuffix string = "VERI*FACTU"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Spain Verifactu V1",
		},
		Extensions:  extensions,
		Validator:   validate,
		Scenarios:   scenarios,
		Normalizer:  normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}
