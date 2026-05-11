// Package tbai provides the TicketBAI addon
package tbai

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the TicketBAI addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "es-tbai"

	// V1 for TicketBAI versions 1.x
	V1 cbc.Key = Key + "-v1"
)

// Official stamps or codes validated by government agencies
const (
	// StampCode contains the code required to be presented alongside
	// the QR code.
	StampCode cbc.Key = "tbai-code"
	// StampQR contains the URL included in the QR code.
	StampQR cbc.Key = "tbai-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("ES-TBAI"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		taxComboRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Spain TicketBAI",
		},
		Extensions:  extensions,
		Normalizer:  normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *bill.Line:
		normalizeBillLine(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
