// Package testaddons provides test-only addons for exercising generic GOBL
// behaviour that is driven by addon definitions, without depending on any
// real addon. Import it from tests only.
package testaddons

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const (
	// CorrectionStamps is the key of a test addon whose correction rules
	// only allow credit notes and require a stamp, as some tax authorities do.
	CorrectionStamps cbc.Key = "test-correction-stamps-v1"
	// StampRef is the stamp required by the CorrectionStamps addon.
	StampRef cbc.Key = "test-ref"
)

func init() {
	tax.RegisterAddonDef(&tax.AddonDef{
		Key:  CorrectionStamps,
		Name: i18n.String{i18n.EN: "Test Correction Stamps"},
		Corrections: tax.CorrectionSet{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Stamps: []cbc.Key{
					StampRef,
				},
			},
		},
	})
}

// Invoice provides a minimal invoice with no tax regime that uses the given
// addons, so their rules are the only ones applied.
func Invoice(addons ...cbc.Key) *bill.Invoice {
	return &bill.Invoice{
		Addons:    tax.WithAddons(addons...),
		Currency:  currency.EUR,
		Series:    "TEST",
		Code:      "123",
		IssueDate: cal.MakeDate(2024, 1, 15),
		Supplier: &org.Party{
			Name: "Test Supplier",
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
			},
		},
	}
}
