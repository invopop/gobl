// Package sdi handles the extensions and validation rules in order to use
// GOBL with the Italian SDI and FatturaPA format.
package sdi

import (
	"errors"

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
	// Key identifies the SDI addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "it-sdi"

	// V1 for SDI's FatturaPA verions 1.x
	V1 cbc.Key = Key + "-v1"

	// KeyFundContribution is the key for the Fund Contribution charge
	KeyFundContribution cbc.Key = "fund-contribution"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("IT-SDI"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		billChargeRules(),
		orgAddressRules(),
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

// validateLatin1String ensures that the item name only contains characters
// from Latin and Latin-1 range (ASCII 0-127 and extended Latin-1 128-255).
func validateLatin1String(val any) error {
	name, _ := val.(string)

	for _, r := range name {
		// Check if the character is outside Latin and Latin-1 range
		// Latin and Latin-1 includes ASCII (0-127) and extended Latin-1 (128-255)
		if r > 255 {
			return errors.New("contains characters outside of Latin and Latin-1 range")
		}
	}
	return nil
}
