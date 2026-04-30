// Package favat provides support for the Polish KSeF FA_VAT format.
package favat

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Polish FA_VAT versions.
const (
	// Key identifies the FA_VAT addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "pl-favat"

	// V3 for FA_VAT version 3 (FA(3)).
	V3 cbc.Key = Key + "-v3"
)

// KSeF official codes to include.
const (
	StampKSeFNumber          cbc.Key = "favat-ksef-number"
	StampKSeFAcquisitionDate cbc.Key = "favat-ksef-acquisition-date"
	StampQR                  cbc.Key = "favat-qr"
)

func init() {
	tax.RegisterAddonDef(newAddonV3())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("PL-FAVAT"),
		is.InContext(tax.AddonIn(V3)),
		billInvoiceRules(),
		taxComboRules(),
		payAdvanceRules(),
	)
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
		Corrections: corrections,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Record:
		normalizePayRecord(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

var corrections = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
		},
		Stamps: []cbc.Key{
			StampKSeFNumber,
		},
	},
}
