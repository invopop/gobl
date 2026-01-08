// Package arca provides the ARCA addon for Argentina.
package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	// V4 for ARCA version 4
	V4 cbc.Key = "ar-arca-v4"
)

// ARCA Official Codes to include in stamps
const (
	// CAE is the code assigned by ARCA to certify that the invoice has been reported.
	StampCAE cbc.Key = "arca-cae"
	// CAEExpiry is the expiry date of the CAE (normally 3 days)
	StampCAEExpiry cbc.Key = "arca-cae-expiry"
	// QR is the QR code URL that contains information about the invoice including the CAE and the CAE expiry date.
	StampQR cbc.Key = "arca-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V4,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA V4",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Argentina ARCA V4",
				},
				URL: "https://www.afip.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG-v4-0.pdf",
			},
		},
		Description: i18n.String{
			i18n.EN: "Support for the Argentina ARCA v4 standard for electronic invoicing.",
		},
		Extensions: extensions,
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Corrections: invoiceCorrectionDefinitions,
		Normalizer:  normalize,
		Validator:   validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeBillInvoice(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateBillInvoice(obj)
	case *bill.Charge:
		return validateBillCharge(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}
