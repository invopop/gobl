// Package sdi handles the extensions and validation rules in order to use
// GOBL with the Italian SDI and FatturaPA format.
package sdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for SDI's FatturaPA verions 1.x
	V1 cbc.Key = "it-sdi-v1"

	// KeyFundContribution is the key for the Fund Contribution charge
	KeyFundContribution cbc.Key = "fund-contribution"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		V1.String(),
		rules.GOBL.Add("IT-SDI-V1"),
		is.HasContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		billChargeRules(),
		taxComboRules(),
		payInstructionsRules(),
		payAdvanceRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Italy SDI FatturaPA v1.x",
		},
		Extensions: extensions,
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Inboxes:    inboxes,
		Normalizer: normalize,
		Scenarios:  scenarios,
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
	case *org.Address:
		normalizeAddress(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

