package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			// Seller identified by org number, via tax ID or ON identity; a VAT
			// registration is not required (bokføringsforskriften § 5-1-2).
			rules.Field("supplier",
				rules.Assert("01", "supplier must have a tax ID or an organization number identity",
					is.Func("has tax ID or org identity", hasSupplierTaxIDOrIdentity),
				),
			),
			// Customer required on non-simplified invoices (§ 5-1-2).
			rules.When(
				is.Func("not simplified", isNotSimplified),
				rules.Field("customer",
					rules.Assert("02", "customer is required on standard invoices", is.Present),
				),
			),
			// Norwegian law has only the kreditnota, no debit note (§ 5-2-7).
			rules.When(
				bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote),
				rules.Field("preceding",
					rules.Assert("03", "preceding document is required for credit notes", is.Present),
				),
			),
		),
	)
}

func hasSupplierTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	if party == nil {
		return false
	}
	if party.TaxID != nil && party.TaxID.Code != "" {
		return true
	}
	for _, id := range party.Identities {
		if id == nil {
			continue
		}
		if id.Type == IdentityTypeOrgNr && id.Code != "" {
			return true
		}
	}
	return false
}

func isNotSimplified(value any) bool {
	inv, ok := value.(*bill.Invoice)
	return ok && inv != nil && !inv.HasTags(tax.TagSimplified)
}
