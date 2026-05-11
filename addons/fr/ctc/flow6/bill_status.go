package flow6

import (
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// schemeIDSIREN is the ISO/IEC 6523 scheme for SIREN identities.
const schemeIDSIREN = "0002"

func normalizeStatus(st *bill.Status) {
	if st == nil {
		return
	}
	// Default Type from the first line's Key — each Flow 6 line key has
	// exactly one associated Status.Type in the process table.
	if st.Type == "" {
		for _, line := range st.Lines {
			if line == nil {
				continue
			}
			if typ, ok := statusTypeForKey(line.Key); ok {
				st.Type = typ
				break
			}
		}
	}

	// Deduce the fr-ctc-role on Issuer / Recipient from the line's
	// (key, type) → side mapping per Annexe A "Acteurs CDV". Saves the
	// caller from spelling out a role that's already implied by the
	// process code. If the caller already set a role, leave it alone.
	if len(st.Lines) > 0 && st.Lines[0] != nil {
		issuerRole, recipientRole := rolesForSide(SideForKeyType(st.Lines[0].Key, st.Type))
		if issuerRole != "" {
			setPartyRoleDefault(st.Issuer, issuerRole)
		}
		if recipientRole != "" {
			setPartyRoleDefault(st.Recipient, recipientRole)
		}
	}

	// Propagate the SE-roled party's SIREN onto Supplier when missing.
	// The seller's SIREN is what populates ref.IssuerTradeParty
	// (MDT-129, BR-FR-CDV-13); when the seller already shows up as
	// Issuer or Recipient, the caller shouldn't have to repeat the
	// identity on Supplier. Only copies the SIREN identity (other
	// fields stay caller-controlled).
	if siren := siRENFromSEParty(st.Issuer, st.Recipient); siren != nil {
		st.Supplier = ensureSIRENOnSupplier(st.Supplier, siren)
	}

	// Surface the CDAR ProcessConditionCode on the document so the
	// wire-level event identifier is visible without consulting the
	// converter. Skip when the caller already pinned a code (e.g.
	// round-tripping a parsed CDV) — the validator catches mismatches.
	if len(st.Lines) > 0 && st.Lines[0] != nil && st.Ext.Get(ExtKeyStatusCode) == "" {
		if code, ok := CDARProcessCodeFor(st.Lines[0].Key, st.Type); ok {
			st.Ext = st.Ext.Set(ExtKeyStatusCode, cbc.Code(code))
		}
	}
}

// siRENFromSEParty returns the first SIREN identity carried by an
// SE-roled party among the given candidates, or nil.
func siRENFromSEParty(candidates ...*org.Party) *org.Identity {
	for _, p := range candidates {
		if p == nil {
			continue
		}
		if p.Ext.Get(ExtKeyRole) != RoleSE {
			continue
		}
		for _, id := range p.Identities {
			if id == nil || id.Ext.IsZero() {
				continue
			}
			if id.Ext.Get(iso.ExtKeySchemeID).String() == schemeIDSIREN {
				return id
			}
		}
	}
	return nil
}

// ensureSIRENOnSupplier returns a Supplier party that carries the
// given SIREN identity, creating one if it was nil and appending the
// identity if the existing Supplier doesn't already carry the same
// SIREN. The identity is shallow-copied so caller-side mutations on
// the source don't leak.
func ensureSIRENOnSupplier(p *org.Party, siren *org.Identity) *org.Party {
	clone := *siren
	if p == nil {
		return &org.Party{Identities: []*org.Identity{&clone}}
	}
	for _, id := range p.Identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(iso.ExtKeySchemeID).String() == schemeIDSIREN &&
			id.Code == siren.Code {
			return p
		}
	}
	p.Identities = append(p.Identities, &clone)
	return p
}

// rolesForSide returns the (Issuer.role, Recipient.role) pair implied
// by an Annexe A side. Empty strings mean "no derivation possible" —
// the caller must supply the role explicitly (e.g. platform-issued
// codes, where Issuer is WK and the recipient role varies).
func rolesForSide(side CDVSide) (issuer, recipient cbc.Code) {
	switch side {
	case CDVSideBuyer:
		return RoleBY, RoleSE
	case CDVSideSeller:
		return RoleSE, RoleBY
	}
	return "", ""
}

func setPartyRoleDefault(p *org.Party, role cbc.Code) {
	if p == nil {
		return
	}
	if !p.Ext.IsZero() && p.Ext.Get(ExtKeyRole) != "" {
		return
	}
	p.Ext = p.Ext.Set(ExtKeyRole, role)
}

func partyHasRole(v any) bool {
	p, ok := v.(*org.Party)
	if !ok || p == nil {
		return false
	}
	if p.Ext.IsZero() {
		return false
	}
	return p.Ext.Get(ExtKeyRole) != ""
}

// partyHasInboxWhenRequired enforces BR-FR-CDV-08: a party whose role
// is not WK (legal representative) or DFH (declarant for VAT grouping)
// must carry a URIID (electronic inbox). We enforce for both issuer
// and supplier to simplify emmission.
func partyHasInboxWhenRequired(v any) bool {
	p, ok := v.(*org.Party)
	if !ok || p == nil {
		return true
	}
	role := p.Ext.Get(ExtKeyRole)
	if role == RoleWK || role == RoleDFH {
		return true
	}
	for _, ib := range p.Inboxes {
		if ib != nil && ib.Code != "" {
			return true
		}
	}
	return false
}

func billStatusRules() *rules.Set {
	return rules.For(new(bill.Status),
		rules.Field("type",
			rules.Assert("01", "status type must be one of: response, update",
				is.In(bill.StatusTypeResponse, bill.StatusTypeUpdate),
			),
		),
		rules.Field("supplier",
			rules.Assert("02", "supplier is required — its SIREN populates ref.IssuerTradeParty (MDT-129, BR-FR-CDV-13)",
				is.Present,
			),
			rules.Assert("03", "supplier must have an identity with ISO/IEC 6523 scheme 0002 (SIREN)",
				is.Func("supplier has SIREN", partyHasSIRENIdentity),
			),
		),
		rules.Field("issuer",
			rules.Assert("14", "issuer is required — maps to ExchangedDocument/IssuerTradeParty (MDG-16) per BR-FR-CDV-CL-03",
				is.Present,
			),
			rules.Assert("15", "issuer.ext.fr-ctc-role must be set; the allowed values depend on ack TypeCode (BR-FR-CDV-CL-03)",
				is.Func("issuer has fr-ctc-role", partyHasRole),
			),
			rules.Assert("20", "issuer must have an electronic address (inbox) when its role is not WK or DFH (BR-FR-CDV-08)",
				is.Func("issuer has inbox unless WK/DFH", partyHasInboxWhenRequired),
			),
		),
		rules.Field("recipient",
			rules.Assert("16", "recipient is required — maps to ExchangedDocument/RecipientTradeParty (MDG-23) per BR-FR-CDV-CL-04",
				is.Present,
			),
			rules.Assert("17", "recipient.ext.fr-ctc-role must be set (BR-FR-CDV-CL-04: BY/DL/SE/AB/SR/PE/PR/II/IV/WK/DFH)",
				is.Func("recipient has fr-ctc-role", partyHasRole),
			),
			rules.Assert("18", "recipient must have an electronic address (inbox) when its role is not WK or DFH (BR-FR-CDV-08)",
				is.Func("recipient has inbox unless WK/DFH", partyHasInboxWhenRequired),
			),
		),
		rules.Field("lines",
			rules.Assert("04", "exactly one status line is required (CDAR carries a single status per CDV message)",
				is.Func("exactly one line", statusHasExactlyOneLine),
			),
			rules.Each(
				rules.Field("doc",
					rules.Assert("05", "status line must reference a document (BR-FR-CDV-10)",
						is.Present,
					),
					rules.Field("code",
						rules.Assert("11", "referenced invoice code is required (BR-FR-CDV-10)",
							is.Present,
						),
					),
					rules.Field("issue_date",
						rules.Assert("12", "referenced invoice issue date is required (BR-FR-CDV-11)",
							is.Present,
						),
					),
				),
				rules.Assert("06", "status line key must be a recognised Flow 6 event",
					is.Func("known Flow 6 status event", statusLineKeyKnown),
				),
				rules.Assert("13", "status lines with key rejected / error / disputed / partially-accepted / suspended require at least one reason (BR-FR-CDV-15)",
					is.Func("reason required for rejection-like statuses", statusLineRequiresReason),
				),
				rules.Assert("09", "Characteristic.ReasonCode must match the fr-ctc-reason-code of some sibling Reason on the same status line",
					is.Func("characteristic reason link resolves", statusLineReasonLinksResolve),
				),
				rules.Assert("10", "Characteristic.TypeCode must be one of the MDT-207 values: MEN, MPA, RAP, ESC, RAB, REM, MAP, MAPTTC, MNA, MNATTC, CBB, DIV, DVA, MAJ",
					is.Func("characteristic type code known", statusLineTypeCodesKnown),
				),
			),
		),
		rules.Assert("08", "Status.Type must be a valid pair with each StatusLine.Key in the Flow 6 process table",
			is.Func("status type consistent with line keys", statusTypeMatchesLines),
		),
		rules.Assert("19", "each Reason's fr-ctc-reason-code must be allowed for the line's CDAR ProcessConditionCode (BR-FR-CDV-CL-09)",
			is.Func("reason codes allowed for status", statusReasonCodesAllowed),
		),
		rules.Assert("07", "status line with key 'paid' on a response status (CDAR 212) must carry a Characteristic complement with Amount (value + currency) set — this is the MEN (BR-FR-CDV-14)",
			is.Func("amount received set when paid response", statusPaidResponseHasAmount),
		),
		rules.Assert("21", "ext.fr-ctc-status-code must match the CDAR ProcessConditionCode implied by (line.Key, Status.Type)",
			is.Func("status code matches key/type", statusCodeMatchesLine),
		),
	)
}

// statusCodeMatchesLine ensures the fr-ctc-status-code ext, when set,
// is consistent with the (line.Key, Status.Type) pair. Empty ext is
// permitted on input — the normalizer fills it.
func statusCodeMatchesLine(v any) bool {
	st, ok := v.(*bill.Status)
	if !ok || st == nil {
		return true
	}
	code := st.Ext.Get(ExtKeyStatusCode).String()
	if code == "" {
		return true
	}
	if len(st.Lines) == 0 || st.Lines[0] == nil {
		return true
	}
	expected, ok := CDARProcessCodeFor(st.Lines[0].Key, st.Type)
	if !ok {
		return true
	}
	return code == expected
}

// statusHasExactlyOneLine enforces the CDAR invariant that a CDV
// message carries one and only one status — a single line on the
// bill.Status. Multiple lines would map to multiple CDARs and must be
// split into separate documents.
func statusHasExactlyOneLine(v any) bool {
	lines, ok := v.([]*bill.StatusLine)
	if !ok {
		return false
	}
	return len(lines) == 1
}

func partyHasSIRENIdentity(v any) bool {
	p, ok := v.(*org.Party)
	if !ok || p == nil {
		return false
	}
	for _, id := range p.Identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(iso.ExtKeySchemeID).String() == schemeIDSIREN {
			return true
		}
	}
	return false
}

func statusLineKeyKnown(v any) bool {
	line, ok := v.(*bill.StatusLine)
	if !ok || line == nil {
		return false
	}
	return statusKeyKnown(line.Key)
}

// statusPaidResponseHasAmount checks BR-FR-CDV-14: every line with
// key=paid on a response status (CDAR 212 Encaissée) must carry a
// Characteristic with TypeCode=MEN and Amount populated. The same
// `paid` key on an update status (CDAR 211 Paiement transmis) does
// not require the MEN.
func statusPaidResponseHasAmount(v any) bool {
	st, ok := v.(*bill.Status)
	if !ok || st == nil {
		return true
	}
	if st.Type != bill.StatusTypeResponse {
		return true
	}
	for _, line := range st.Lines {
		if line == nil || line.Key != bill.StatusEventPaid {
			continue
		}
		if !lineHasMENAmount(line) {
			return false
		}
	}
	return true
}

// lineHasMENAmount returns true if the given line carries a
// flow6.Characteristic complement with TypeCode=MEN and a populated
// Amount (value + currency).
func lineHasMENAmount(line *bill.StatusLine) bool {
	for _, obj := range line.Complements {
		if obj == nil {
			continue
		}
		c, ok := obj.Instance().(*Characteristic)
		if !ok || c == nil {
			continue
		}
		if c.TypeCode != TypeCodeAmountReceived {
			continue
		}
		if c.Amount == nil || c.Amount.Currency == "" {
			continue
		}
		return true
	}
	return false
}

// statusLineTypeCodesKnown ensures every Characteristic.TypeCode on
// the line is one of the MDT-207 controlled values.
func statusLineTypeCodesKnown(v any) bool {
	line, ok := v.(*bill.StatusLine)
	if !ok || line == nil {
		return true
	}
	for _, obj := range line.Complements {
		if obj == nil {
			continue
		}
		c, ok := obj.Instance().(*Characteristic)
		if !ok || c == nil || c.TypeCode == "" {
			continue
		}
		if !slices.Contains(typeCodes, c.TypeCode) {
			return false
		}
	}
	return true
}

// statusLineReasonLinksResolve ensures that every Characteristic on the
// line whose ReasonCode is set matches the fr-ctc-reason-code of some
// sibling bill.Reason on the same line. An unset ReasonCode is allowed.
func statusLineReasonLinksResolve(v any) bool {
	line, ok := v.(*bill.StatusLine)
	if !ok || line == nil {
		return true
	}
	if len(line.Complements) == 0 {
		return true
	}
	for _, obj := range line.Complements {
		if obj == nil {
			continue
		}
		c, ok := obj.Instance().(*Characteristic)
		if !ok || c == nil || c.ReasonCode == "" {
			continue
		}
		if !lineHasReasonCode(line, c.ReasonCode) {
			return false
		}
	}
	return true
}

func lineHasReasonCode(line *bill.StatusLine, code cbc.Code) bool {
	for _, r := range line.Reasons {
		if r == nil {
			continue
		}
		if r.Ext.Get(ExtKeyReasonCode) == code {
			return true
		}
	}
	return false
}

// reasonRequiredStatusKeys lists the Flow 6 status-line keys that BR-FR-CDV-15
// designates as carrying mandatory motifs. The 501 "IRRECEVABLE" status
// from the CSV is not in our process table (it's PPF-ingress-only) and
// is deliberately omitted — if we ever model it, add it here.
var reasonRequiredStatusKeys = []cbc.Key{
	bill.StatusEventRejected,
	bill.StatusEventError,
	bill.StatusEventQuerying,
	StatusEventDisputed,
	StatusEventPartiallyAccepted,
}

// statusReasonCodesAllowed enforces BR-FR-CDV-CL-09 at the
// bill.Status level: each Reason on each line must carry an
// fr-ctc-reason-code permitted for the (line.Key, st.Type) →
// ProcessConditionCode pair. Lives at the status level because the
// pair-lookup needs Type — line-only keys like `paid` are ambiguous
// (211 update vs 212 response) without it.
func statusReasonCodesAllowed(v any) bool {
	st, ok := v.(*bill.Status)
	if !ok || st == nil {
		return true
	}
	for _, line := range st.Lines {
		if line == nil {
			continue
		}
		processCode, ok := CDARProcessCodeFor(line.Key, st.Type)
		if !ok {
			continue
		}
		for _, r := range line.Reasons {
			if r == nil {
				continue
			}
			code := r.Ext.Get(ExtKeyReasonCode).String()
			if code == "" {
				continue
			}
			if !ReasonCodeAllowedForProcessCode(code, processCode) {
				return false
			}
		}
	}
	return true
}

func statusLineRequiresReason(v any) bool {
	line, ok := v.(*bill.StatusLine)
	if !ok || line == nil {
		return true
	}
	if !slices.Contains(reasonRequiredStatusKeys, line.Key) {
		return true
	}
	return len(line.Reasons) > 0
}

func statusTypeMatchesLines(v any) bool {
	st, ok := v.(*bill.Status)
	if !ok || st == nil {
		return true
	}
	for _, line := range st.Lines {
		if line == nil {
			continue
		}
		expected, ok := statusTypeForKey(line.Key)
		if !ok {
			continue
		}
		if expected != st.Type {
			return false
		}
	}
	return true
}

// -- bill.Reason --------------------------------------------------------

// normalizeReason fills in the other side of the Reason.Key ↔ Ext
// relationship when exactly one side is set. The extension carries the
// exact CDAR ReasonCode; the Key is the bucket.
func normalizeReason(r *bill.Reason) {
	if r == nil {
		return
	}
	ext := r.Ext.Get(ExtKeyReasonCode).String()
	switch {
	case r.Key == "" && ext != "":
		if key, ok := ReasonKeyFor(ext); ok {
			r.Key = key
		}
	case r.Key != "" && ext == "":
		if code, ok := CDARReasonCodeFor(r.Key); ok {
			r.Ext = r.Ext.Set(ExtKeyReasonCode, cbc.Code(code))
		}
	}
}

func billReasonRules() *rules.Set {
	return rules.For(new(bill.Reason),
		rules.Field("key",
			rules.AssertIfPresent("01", "reason key is not a recognised bill.ReasonKeys value",
				is.In(reasonKeyAnySlice()...),
			),
		),
		rules.Assert("02", "fr-ctc-reason-code must be a known CDAR code and its bucket must match reason.key",
			is.Func("reason ext code consistent with key", reasonExtMatchesKey),
		),
	)
}

var validReasonKeys = func() []cbc.Key {
	keys := make([]cbc.Key, 0, len(bill.ReasonKeys))
	for _, def := range bill.ReasonKeys {
		keys = append(keys, def.Key)
	}
	return keys
}()

func reasonKeyAnySlice() []any {
	out := make([]any, len(validReasonKeys))
	for i, k := range validReasonKeys {
		out[i] = k
	}
	return out
}

func reasonExtMatchesKey(v any) bool {
	r, ok := v.(*bill.Reason)
	if !ok || r == nil {
		return true
	}
	ext := r.Ext.Get(ExtKeyReasonCode).String()
	if ext == "" {
		return true
	}
	bucket, ok := ReasonKeyFor(ext)
	if !ok {
		return false
	}
	// Key may be empty when normalization has not yet run; the
	// normalizer fills it from the ext.
	if r.Key == "" {
		return true
	}
	return bucket == r.Key
}

// -- bill.Action --------------------------------------------------------

var validActionKeys = func() []cbc.Key {
	keys := make([]cbc.Key, 0, len(bill.ActionKeys))
	for _, def := range bill.ActionKeys {
		keys = append(keys, def.Key)
	}
	return keys
}()

func actionKeyAnySlice() []any {
	out := make([]any, len(validActionKeys))
	for i, k := range validActionKeys {
		out[i] = k
	}
	return out
}

func billActionRules() *rules.Set {
	return rules.For(new(bill.Action),
		rules.Field("key",
			rules.AssertIfPresent("01", "action key is not a recognised bill.ActionKeys value",
				is.In(actionKeyAnySlice()...),
			),
		),
	)
}
