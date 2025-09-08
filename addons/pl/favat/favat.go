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
	V2 cbc.Key = "pl-favat-v2"
	V3 cbc.Key = "pl-favat-v3"
)

// KSeF official codes to include.
const (
	StampID   cbc.Key = "favat-id"
	StampHash cbc.Key = "favat-hash"
	StampQR   cbc.Key = "favat-qr"
)

func init() {
	tax.RegisterAddonDef(newAddonV2())
	// V3 coming soon...
}

func newAddonV2() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V2,
		Name: i18n.String{
			i18n.EN: "Polish KSeF FA_VAT v2.x",
		},
		Tags: []*tax.TagSet{
			invoiceTags, // scenarios.go
		},
		Extensions:  extensionKeys,
		Scenarios:   scenarios, // scenarios.go
		Normalizer:  normalize,
		Validator:   validate,
		Corrections: corrections,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Advance:
		normalizePayAdvance(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *pay.Instructions:
		return validatePayInstructions(obj)
	case *pay.Advance:
		return validatePayAdvance(obj)
	case *pay.Terms:
		return validatePayTerms(obj)
	}
	return nil
}

var corrections = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
		},
		ReasonRequired: true,
		Stamps: []cbc.Key{
			StampID,
		},
		Extensions: []cbc.Key{
			ExtKeyEffectiveDate,
		},
	},
}
