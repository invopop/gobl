package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Note for reviewers: the rules below intentionally cover only the Norway-
// specific requirements. GOBL core already validates the universal fields
// (type, dates, supplier presence + name, customer name when a tax ID is set,
// line prices, totals), so they are not repeated here.
//
// The following genuine legal requirements are deliberately NOT enforced at the
// regime layer:
//   - Seller address: bokføringsforskriften § 5-1-2 requires the seller's head-
//     office address only for AS/ASA and foreign branches (NUF); the minimum
//     for a seller is name + organisasjonsnummer. We do not require an address
//     generally, as that would over-enforce against e.g. sole proprietorships.
//   - Foretaksregisteret: § 5-1-2 / foretaksregisterloven § 10-2 require the
//     word "Foretaksregisteret" on the document, again only for AS/ASA and
//     foreign branches.
//   - VAT amount in NOK: § 5-1-1 nr. 6 requires the VAT amount to be stated in
//     NOK even on foreign-currency invoices.
//
// The first two depend on the supplier's legal form, which GOBL core does not
// model; the third is a currency-conversion concern. gobl expects all three to
// be handled by the EHF/SAF-T addon rather than the regime.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			// The seller must be identified by their organisasjonsnummer
			// (bokføringsforskriften § 5-1-2). Every Norwegian business has one
			// even when not VAT-registered, so we require either a tax ID or an
			// `ON` organization identity — NOT specifically a VAT registration.
			// Businesses below the NOK 50,000 turnover threshold are not VAT-
			// registered and cannot charge MVA, but may still issue invoices.
			rules.Field("supplier",
				rules.Assert("01", "supplier must have a tax ID or an organization number identity",
					is.Func("has tax ID or org identity", hasSupplierTaxIDOrIdentity),
				),
			),
			// Standard invoices must identify the customer (bokføringsforskriften
			// § 5-1-2 first paragraph). Simplified invoices (e.g. cash sales)
			// relax this.
			rules.When(
				is.Func("not simplified", isNotSimplified),
				rules.Field("customer",
					rules.Assert("02", "customer is required on standard invoices", is.Present),
				),
			),
			// Corrections: Norwegian bookkeeping law recognises only the
			// kreditnota (bokføringsforskriften § 5-2-7). There is no debit-note
			// concept, so only credit notes must reference the document they
			// reverse.
			rules.When(
				bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote),
				rules.Field("preceding",
					rules.Assert("03", "preceding document is required for credit notes", is.Present),
				),
			),
		),
	)
}

// hasSupplierTaxIDOrIdentity reports whether the supplier carries either a tax
// ID code or an `ON` organisasjonsnummer identity, satisfying the seller
// identification requirement of bokføringsforskriften § 5-1-2.
func hasSupplierTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	if party == nil {
		return false
	}
	if party.TaxID != nil && party.TaxID.Code != "" {
		return true
	}
	for _, id := range party.Identities {
		if id.Type == IdentityTypeOrgNr {
			return true
		}
	}
	return false
}

func isNotSimplified(value any) bool {
	inv, ok := value.(*bill.Invoice)
	return ok && inv != nil && !inv.HasTags(tax.TagSimplified)
}
