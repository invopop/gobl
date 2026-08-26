// Package sdi handles the extensions and validation rules in order to use
// GOBL with the Italian SDI and FatturaPA format.
package sdi

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the SDI addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "it-sdi"

	// V1 for SDI's FatturaPA versions 1.x
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
		orgAttributeRules(),
		taxComboRules(),
		payInstructionsRules(),
		payAdvanceRules(),
	)
	norm.RegisterWithGuard(
		is.InContext(tax.AddonIn(V1)),
		norm.For(normalizeInvoice),
		norm.For(normalizePayInstructions),
		norm.For(normalizePayRecord),
		norm.For(normalizeAddress),
		norm.For(normalizeTaxCombo),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Italy SDI FatturaPA v1.x",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Italy exchanges electronic invoices in the FatturaPA XML format through the tax
				authority's Sistema di Interscambio (SDI). This addon ensures GOBL documents carry
				the fields and extensions needed to produce valid FatturaPA files with
				[gobl.fatturapa](https://github.com/invopop/gobl.fatturapa).

				## Customer identification

				An Italian customer is identified by a partita IVA (VAT number, in the party's
				~tax_id~) or a codice fiscale (an identity with the key ~it-fiscal-code~), and
				routed by a codice destinatario (~it-sdi-code~ inbox) or a PEC address
				(~it-sdi-pec~ inbox). A customer outside Italy needs only the country in ~tax_id~.

				Where a party has no code of its own, FatturaPA expects a placeholder rather than
				an empty field, and the conversion supplies it: ~0000000~ for a customer with no
				tax code, or an Italian customer with no inbox; ~XXXXXXX~ as the recipient code for
				a customer outside Italy; and ~OO99999999999~ in place of a non-EU business's own
				tax number, since only EU VAT numbers are meaningful to SDI. Leave the field out
				rather than writing these values yourself.
			`),
		},
		Extensions: extensions,
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Inboxes:   inboxes,
		Scenarios: scenarios,
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
