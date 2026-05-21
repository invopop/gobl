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

func normalizeStatus(st *bill.Status) {
	if st == nil {
		return
	}

	// If the caller pinned fr-ctc-flow6-status-code but left
	// Line.Key / Status.Type blank, fill them from the process table.
	if code := st.Ext.Get(ExtKeyStatusCode).String(); code != "" {
		if key, typ, ok := StatusKeyFor(code); ok {
			if len(st.Lines) > 0 && st.Lines[0] != nil && st.Lines[0].Key == "" {
				st.Lines[0].Key = key
			}
			if st.Type == "" {
				st.Type = typ
			}
		}
	}

	// Default Type from the first line's Key.
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

	// Deduce the fr-ctc-flow6-role on Issuer / Recipient from the
	// line's (key, type) → side mapping per Annexe A "Acteurs CDV".
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
	if siren := sirenFromSEParty(st.Issuer, st.Recipient); siren != nil {
		st.Supplier = ensureSIRENOnSupplier(st.Supplier, siren)
	}

	// Surface the CDAR ProcessConditionCode on the document.
	if len(st.Lines) > 0 && st.Lines[0] != nil && st.Ext.Get(ExtKeyStatusCode) == "" {
		if code, ok := CDARProcessCodeFor(st.Lines[0].Key, st.Type); ok {
			st.Ext = st.Ext.Set(ExtKeyStatusCode, cbc.Code(code))
		}
	}
}

// sirenFromSEParty returns the first SIREN identity carried by an
// SE-roled party among the given candidates, or nil. Relies on the
// addon normaliser having set iso-scheme-id=0002 on SIREN identities
// (see normalizeIdentity in org.go).
func sirenFromSEParty(candidates ...*org.Party) *org.Identity {
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
			if id.Ext.Get(iso.ExtKeySchemeID).String() == identitySchemeIDSIREN {
				return id
			}
		}
	}
	return nil
}

// ensureSIRENOnSupplier returns a Supplier party that carries the
// given SIREN identity, creating one if it was nil and appending the
// identity if the existing Supplier doesn't already carry the same
// SIREN.
func ensureSIRENOnSupplier(p *org.Party, siren *org.Identity) *org.Party {
	clone := *siren
	if p == nil {
		return &org.Party{Identities: []*org.Identity{&clone}}
	}
	for _, id := range p.Identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(iso.ExtKeySchemeID).String() == identitySchemeIDSIREN &&
			id.Code == siren.Code {
			return p
		}
	}
	p.Identities = append(p.Identities, &clone)
	return p
}

// rolesForSide returns the (Issuer.role, Recipient.role) pair implied
// by an Annexe A side.
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
// is not WK (dematerialisation platform / operator) and not DFH
// (PPF) must carry a URIID (electronic inbox).
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
			rules.Assert("02", "supplier is required (BR-FR-CDV-13)",
				is.Present,
			),
			rules.Assert("03", "supplier must have an identity with ISO/IEC 6523 scheme 0002 (SIREN)",
				is.Func("supplier has SIREN", statusPartyHasSIRENIdentity),
			),
		),
		rules.Field("issuer",
			rules.Assert("14", "issuer is required (BR-FR-CDV-CL-03)",
				is.Present,
			),
			rules.Assert("15", "issuer.ext.fr-ctc-flow6-role must be set. Allowed values are BY/DL/SE/AB/SR/PE/PR/II/IV/WK/DFH (BR-FR-CDV-CL-03)",
				is.Func("issuer has fr-ctc-flow6-role", partyHasRole),
			),
			rules.Assert("20", "issuer must have an electronic address (inbox) when its role is not WK or DFH (BR-FR-CDV-08)",
				is.Func("issuer has inbox unless WK/DFH", partyHasInboxWhenRequired),
			),
		),
		rules.Field("recipient",
			rules.Assert("16", "recipient is required (BR-FR-CDV-CL-04)",
				is.Present,
			),
			rules.Assert("17", "recipient.ext.fr-ctc-flow6-role must be set. Allowed values are BY/DL/SE/AB/SR/PE/PR/II/IV/WK/DFH (BR-FR-CDV-CL-04)",
				is.Func("recipient has fr-ctc-flow6-role", partyHasRole),
			),
			rules.Assert("18", "recipient must have an electronic address (inbox) when its role is not WK or DFH (BR-FR-CDV-08)",
				is.Func("recipient has inbox unless WK/DFH", partyHasInboxWhenRequired),
			),
		),
		rules.Field("lines",
			rules.Assert("04", "exactly one status line is required",
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
			),
		),
		rules.Assert("08", "Status.Type must be a valid pair with each StatusLine.Key in the Flow 6 process table",
			is.Func("status type consistent with line keys", statusTypeMatchesLines),
		),
		rules.Assert("19", "each Reason's fr-ctc-flow6-reason-code must be allowed for the line's CDAR ProcessConditionCode (BR-FR-CDV-CL-09)",
			is.Func("reason codes allowed for status", statusReasonCodesAllowed),
		),
		rules.Assert("21", "ext.fr-ctc-flow6-status-code must match the CDAR ProcessConditionCode implied by (line.Key, Status.Type)",
			is.Func("status code matches key/type", statusCodeMatchesLine),
		),
		rules.Assert("22", "every CDV party identity scheme must be in the Flow 6 allow-list (STC 0231 is rejected — it is a Flow 2 invoice concept)",
			is.Func("CDV identity schemes allowed", statusPartiesIdentitySchemesAllowed),
		),
	)
}

// statusPartiesIdentitySchemesAllowed rejects any identity whose
// iso-scheme-id falls outside allowedFlow6IdentitySchemes on any of
// the four party slots on a bill.Status.
func statusPartiesIdentitySchemesAllowed(v any) bool {
	st, ok := v.(*bill.Status)
	if !ok || st == nil {
		return true
	}
	for _, p := range []*org.Party{st.Supplier, st.Customer, st.Issuer, st.Recipient} {
		if p == nil {
			continue
		}
		for _, id := range p.Identities {
			if id == nil || id.Ext.IsZero() {
				continue
			}
			scheme := id.Ext.Get(iso.ExtKeySchemeID).String()
			if scheme != "" && !slices.Contains(allowedFlow6IdentitySchemes, scheme) {
				return false
			}
		}
	}
	return true
}

// statusCodeMatchesLine ensures the fr-ctc-flow6-status-code ext, when
// set, is consistent with the (line.Key, Status.Type) pair.
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
// message carries one and only one status.
func statusHasExactlyOneLine(v any) bool {
	lines, ok := v.([]*bill.StatusLine)
	if !ok {
		return false
	}
	return len(lines) == 1
}

func statusPartyHasSIRENIdentity(v any) bool {
	p, ok := v.(*org.Party)
	if !ok || p == nil {
		return false
	}
	for _, id := range p.Identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(iso.ExtKeySchemeID).String() == identitySchemeIDSIREN {
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

// reasonRequiredStatusKeys lists the Flow 6 status-line keys that
// BR-FR-CDV-15 designates as carrying mandatory motifs.
var reasonRequiredStatusKeys = []cbc.Key{
	bill.StatusEventRejected,
	bill.StatusEventError,
	bill.StatusEventQuerying,
	StatusEventDisputed,
	StatusEventPartiallyAccepted,
}

// statusReasonCodesAllowed enforces BR-FR-CDV-CL-09 at the
// bill.Status level.
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
// relationship when exactly one side is set.
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
		rules.Assert("02", "fr-ctc-flow6-reason-code must be a known CDAR code and its bucket must match reason.key",
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
