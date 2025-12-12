package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	V4 cbc.Key = "ar-arca-v4"
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
		Extensions:  extensions,
		Scenarios:   scenarios,
		Corrections: invoiceCorrectionDefinitions,
		Normalizer:  normalize,
		Validator:   validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *bill.Charge:
		normalizeCharge(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *bill.Charge:
		return validateCharge(obj)
	}
	return nil
}
