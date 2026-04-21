// Package bis provides an addon that enforces the Peppol BIS Billing 3.0
// rule set on top of the EN 16931 ruleset. It covers the 55 base Peppol rules
// (PEPPOL-COMMON and PEPPOL-EN16931) as well as the national CIUS rules for
// Denmark, Germany, Greece, Iceland, Italy, the Netherlands, Norway and Sweden.
//
// Sibling packages will cover other Peppol profiles (e.g. `peppol/pint` for
// Peppol International).
package bis

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// V3 is the key for the Peppol BIS Billing 3.0 addon.
	V3 cbc.Key = "peppol-bis-v3"
)

func init() {
	tax.RegisterAddonDef(newAddon())

	rules.RegisterWithGuard(
		V3.String(),
		rules.GOBL.Add("PEPPOL-BIS-V3"),
		is.InContext(tax.AddonIn(V3)),
		// Base rules
		billInvoiceRules(),
		billLineRules(),
		payInstructionsRules(),
		orgPartyRules(),
		orgIdentityRules(),
		orgInboxRules(),
		taxComboRules(),
		// National CIUS rule sets — each guards internally on supplier country.
		billInvoiceRulesDK(), orgPartyRulesDK(), payInstructionsRulesDK(), orgItemRulesDK(),
		billInvoiceRulesDE(), payInstructionsRulesDE(),
		billInvoiceRulesGR(), orgPartyRulesGR(),
		billInvoiceRulesIS(), payInstructionsRulesIS(),
		orgPartyRulesIT(),
		billInvoiceRulesNL(), orgPartyRulesNL(), payInstructionsRulesNL(),
		orgPartyRulesNO(),
		billInvoiceRulesSE(), orgPartyRulesSE(), payInstructionsRulesSE(),
	)
}

// supplierCountryIs returns a Test that passes when the given invoice's
// supplier country matches the provided code. Used as an inner guard in
// country-specific rule sets so they only fire for the relevant suppliers.
func supplierCountryIs(country l10n.Code) rules.Test {
	return is.Func("supplier country "+country.String(), func(val any) bool {
		return invoiceSupplierCountry(val) == country
	})
}

// invoiceSupplierCountry extracts the supplier country from an invoice,
// preferring TaxID.Country and falling back to the first address's country.
func invoiceSupplierCountry(val any) l10n.Code {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return ""
	}
	return partyCountry(inv.Supplier)
}

// partyCountry returns the party's country, checking TaxID first then first
// address. Used for country guards that target the supplier's tax jurisdiction.
func partyCountry(p *org.Party) l10n.Code {
	if p == nil {
		return ""
	}
	if p.TaxID != nil && !p.TaxID.Country.Empty() {
		return p.TaxID.Country.Code()
	}
	return partyAddressCountry(p)
}

// partyAddressCountry returns the party's first-address country, ignoring the
// TaxID. Use this for rules whose schematron targets cac:PostalAddress/cac:Country
// (e.g. PEPPOL-EN16931-R002's DK exception).
func partyAddressCountry(p *org.Party) l10n.Code {
	if p == nil {
		return ""
	}
	if len(p.Addresses) > 0 && p.Addresses[0] != nil {
		return p.Addresses[0].Country.Code()
	}
	return ""
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V3,
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Name: i18n.String{
			i18n.EN: "Peppol BIS Billing 3.0",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Enforces the Peppol BIS Billing 3.0 rule set on top of the EN 16931 base.

				Covers the 55 base Peppol rules (PEPPOL-COMMON-Rxxx and PEPPOL-EN16931-*) as
				well as the national CIUS rule sets for Denmark, Germany, Greece, Iceland,
				Italy, the Netherlands, Norway and Sweden. National rules are applied
				automatically based on the supplier's country.

				This addon does not define the Peppol ProfileID or CustomizationID — those
				are emitted by the converter (` + "`gobl.ubl`" + `) based on the Context
				selected at serialization time.

				For the authoritative rule list, see:
				https://docs.peppol.eu/poacc/billing/3.0/rules/ubl-peppol/
			`),
		},
		Identities: identities,
		Normalizer: normalize,
	}
}

// normalize dispatches addon-specific normalization to the right handler.
func normalize(doc any) {
	switch obj := doc.(type) {
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
