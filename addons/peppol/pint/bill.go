package pint

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("supplier",
			rules.Assert("01", "invoice supplier must have a tax ID or identity (IBR-CO-26)",
				is.Func("has identifier", partyHasIdentifier),
			),
			rules.Field("inboxes",
				rules.Assert("02", "invoice supplier electronic address is required (IBR-081)", is.Present),
			),
		),
		rules.Field("customer",
			rules.Assert("03", "invoice customer is required (IBR-007)", is.Present),
			rules.Field("inboxes",
				rules.Assert("04", "invoice customer electronic address is required (IBR-080)", is.Present),
			),
		),
	)
}

// partyHasIdentifier reports whether the party can be identified through a tax
// identity code or at least one organisation identity.
func partyHasIdentifier(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true // presence of the party is asserted elsewhere
	}
	if p.TaxID != nil && p.TaxID.Code != "" {
		return true
	}
	return len(p.Identities) > 0
}
