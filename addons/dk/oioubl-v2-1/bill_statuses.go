package oioubl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// billStatusRules returns the OIOUBL 2.1 rule set for bill.Status,
// targeting Invoice Response (UBL ApplicationResponse with Type "response").
func billStatusRules() *rules.Set {
	return rules.For(new(bill.Status),
		rules.When(isResponseType,
			rules.Field("code",
				rules.Assert("05", "code is required (F-APR005)", is.Present),
			),
			rules.Field("supplier",
				rules.Field("inboxes",
					rules.Assert("01", "supplier inboxes are required (F-APR008)", is.Present),
				),
				rules.Assert("06", "supplier must have a tax ID or identities (F-APR040)",
					is.Func("has tax id or identities", partyHasTaxIDOrIdentities)),
				rules.Assert("07", "supplier must have a name or identities (F-LIB022)",
					is.Func("has name or identities", partyHasNameOrIdentities)),
			),
			rules.Field("customer",
				rules.Assert("02", "customer is required for a response", is.Present),
				rules.Field("inboxes",
					rules.Assert("03", "customer inboxes are required (F-APR012)", is.Present),
				),
				rules.Assert("08", "customer must have a name or identities (F-LIB022)",
					is.Func("has name or identities", partyHasNameOrIdentities)),
			),
			rules.Field("issuer",
				rules.Field("inboxes",
					rules.Assert("09", "issuer inboxes are required when issuer is set (F-APR047)", is.Present),
				),
				rules.Assert("10", "issuer must have a tax ID or identities when set (F-APR048)",
					is.Func("has tax id or identities", partyHasTaxIDOrIdentities)),
				rules.Assert("11", "issuer must have a name or identities when set (F-LIB022)",
					is.Func("has name or identities", partyHasNameOrIdentities)),
			),
			rules.Field("lines",
				rules.Each(
					rules.Field("doc",
						rules.Assert("04", "line document reference is required for a response (cf. F-APR016, F-APR025)", is.Present),
					),
				),
			),
		),
	)
}

var isResponseType = is.Func("response status type", func(val any) bool {
	st, ok := val.(*bill.Status)
	return ok && st != nil && st.Type == bill.StatusTypeResponse
})

func partyHasTaxIDOrIdentities(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	return p.TaxID != nil || len(p.Identities) > 0
}

func partyHasNameOrIdentities(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	return p.Name != "" || len(p.Identities) > 0
}
