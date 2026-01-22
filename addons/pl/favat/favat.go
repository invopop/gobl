// Package favat provides support for the Polish KSeF FA_VAT format.
package favat

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// Polish FA_VAT versions.
const (
	V3 cbc.Key = "pl-favat-v3"
)

// KSeF official codes to include.
const (
	StampKSEFNumber cbc.Key = "favat-ksef-number"
	StampHash       cbc.Key = "favat-hash"
	StampQR         cbc.Key = "favat-qr"
)

func init() {
	tax.RegisterAddonDef(newAddonV3())
}

func newAddonV3() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V3,
		Name: i18n.String{
			i18n.EN: "Polish KSeF FA_VAT FA(3)",
		},
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Extensions:  extensionKeys,
		Scenarios:   scenarios,
		Normalizer:  normalize,
		Validator:   validate,
		Corrections: corrections,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Advance:
		normalizePayAdvance(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateBillInvoice(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}

var corrections = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
		},
		Stamps: []cbc.Key{
			StampKSEFNumber,
		},
	},
}
