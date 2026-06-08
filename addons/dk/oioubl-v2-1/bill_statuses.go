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
			// The converter sends the response FROM the responder TO the
			// originator: the customer (or issuer) becomes the SenderParty
			// (F-APR008 endpoint, F-APR040 one PartyLegalEntity) and the supplier
			// (or recipient) becomes the ReceiverParty (F-APR012 endpoint,
			// F-APR041 at most one PartyLegalEntity). The citations below follow
			// that mapping rather than the GOBL field name.
			rules.Field("supplier",
				rules.Field("inboxes",
					rules.Assert("01", "supplier inboxes are required (F-APR012)", is.Present),
				),
				rules.Assert("06", "supplier must have a tax ID or identities (F-APR041)",
					is.Func("has tax id or identities", partyHasTaxIDOrIdentities)),
				rules.Assert("07", "supplier must have a name or legal identity (F-LIB022)",
					is.Func("has name or legal identity", partyHasNameOrLegalIdentity)),
			),
			rules.Field("customer",
				rules.Assert("02", "customer is required for a response", is.Present),
				rules.Field("inboxes",
					rules.Assert("03", "customer inboxes are required (F-APR008)", is.Present),
				),
				rules.Assert("08", "customer must have a name or legal identity (F-LIB022)",
					is.Func("has name or legal identity", partyHasNameOrLegalIdentity)),
			),
			rules.Field("issuer",
				rules.Field("inboxes",
					rules.Assert("09", "issuer inboxes are required when issuer is set (F-APR008)", is.Present),
				),
				rules.Assert("10", "issuer must have a tax ID or identities when set (F-APR040)",
					is.Func("has tax id or identities", partyHasTaxIDOrIdentities)),
				rules.Assert("11", "issuer must have a name or legal identity when set (F-LIB022)",
					is.Func("has name or legal identity", partyHasNameOrLegalIdentity)),
			),
			// A recipient, when present, replaces the supplier as the document's
			// ReceiverParty, so it carries the same requirements as the supplier.
			rules.Field("recipient",
				rules.Field("inboxes",
					rules.Assert("12", "recipient inboxes are required when recipient is set (F-APR012)", is.Present),
				),
				rules.Assert("13", "recipient must have a tax ID or identities when set (F-APR041)",
					is.Func("has tax id or identities", partyHasTaxIDOrIdentities)),
				rules.Assert("14", "recipient must have a name or legal identity when set (F-LIB022)",
					is.Func("has name or legal identity", partyHasNameOrLegalIdentity)),
			),
			rules.Field("lines",
				// An OIOUBL ApplicationResponse carries exactly one Response for one
				// referenced document (F-APR051 / F-APR054); the converter maps each
				// status line to a DocumentResponse, so it must hold a single line.
				rules.Assert("16", "a response carries exactly one document response (F-APR051 / F-APR054)", is.Length(1, 1)),
				rules.Each(
					// Only the four events the responsecode-1.1 code list accepts are
					// representable; everything else (issued, processing, paid, …) has
					// no OIOUBL response code (F-APR018).
					rules.Field("key",
						rules.Assert("15", "response status event must be one OIOUBL supports (F-APR018)",
							is.In(
								bill.StatusEventAccepted,
								bill.StatusEventRejected,
								bill.StatusEventAcknowledged,
								bill.StatusEventError,
							)),
					),
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

// partyHasNameOrLegalIdentity mirrors what the gobl.ubl converter can turn into
// a cac:PartyLegalEntity (F-APR040/044/048): a registration name (from the party
// name) or a CompanyID (from a legal-scope identity). A party carrying only a
// tax-scope identity and no name produces no PartyLegalEntity and would fail the
// schematron, so it must be rejected here.
func partyHasNameOrLegalIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if p.Name != "" {
		return true
	}
	for _, id := range p.Identities {
		if id != nil && id.Scope == org.IdentityScopeLegal {
			return true
		}
	}
	return false
}
