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
	// Default party role for the two structural slots. Issuer and
	// Recipient are left untouched: their role is context-dependent.
	setPartyRoleDefault(st.Supplier, RoleSE)
	setPartyRoleDefault(st.Customer, RoleBY)
}

func setPartyRoleDefault(p *org.Party, role cbc.Code) {
	if p == nil || p.Ext.Get(ExtKeyRole) != "" {
		return
	}
	p.Ext = p.Ext.Set(ExtKeyRole, role)
}

func billStatusRules() *rules.Set {
	return rules.For(new(bill.Status),
		rules.Field("type",
			rules.Assert("01", "status type must be one of: response, update",
				is.In(bill.StatusTypeResponse, bill.StatusTypeUpdate),
			),
		),
		rules.Field("supplier",
			rules.Assert("02", "supplier is required on Flow 6 status messages",
				is.Present,
			),
			rules.Assert("03", "supplier must have an identity with ISO/IEC 6523 scheme 0002 (SIREN)",
				is.Func("supplier has SIREN", partyHasSIRENIdentity),
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
				rules.Assert("07", "status line with key 'paid' (CDAR 212) must carry a Characteristic complement with Amount (value + currency) set — this is the MEN",
					is.Func("amount received set when paid", statusLinePaidHasAmount),
				),
				rules.Assert("09", "Characteristic.ReasonCode must match the fr-ctc-reason-code of some sibling Reason on the same status line",
					is.Func("characteristic reason link resolves", statusLineReasonLinksResolve),
				),
				rules.Assert("10", "Characteristic.TypeCode must be one of the MDT-207 values: MEN, MPA, RAP, ESC, RAB, REM, MAP, MAPTTC, MNA, MNATTC, CBB, DIV, DVA, MAJ",
					is.Func("characteristic type code known", statusLineTypeCodesKnown),
				),
			),
		),
		rules.Assert("08", "Status.Type must match the Type implied by each StatusLine.Key",
			is.Func("status type consistent with line keys", statusTypeMatchesLines),
		),
	)
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
	_, ok = statusTypeForKey(line.Key)
	return ok
}

// statusLinePaidHasAmount checks that a paid StatusLine carries a
// Characteristic complement with TypeCode=MEN and Amount populated
// (both value and currency). Other payment-related TypeCodes (MPA,
// RAP, etc.) may coexist on the same line but do not substitute for
// the mandatory MEN.
func statusLinePaidHasAmount(v any) bool {
	line, ok := v.(*bill.StatusLine)
	if !ok || line == nil {
		return true
	}
	if line.Key != bill.StatusEventPaid {
		return true
	}
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
	StatusEventDisputed,
	StatusEventPartiallyAccepted,
	StatusEventSuspended,
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
