package flow6

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Helpers --------------------------------------------------------------

// statusSupplierParty returns a French supplier party with a SIREN
// identity, an inbox satisfying BR-FR-CDV-08, and an SE role.
func statusSupplierParty() *org.Party {
	return &org.Party{
		Name: "VENDEUR SARL",
		Identities: []*org.Identity{
			{
				Code: "356000000",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			},
		},
		Inboxes: []*org.Inbox{{Scheme: "0225", Code: "356000000_PEP"}},
		Ext:     tax.MakeExtensions().Set(ExtKeyRole, RoleSeller),
	}
}

func statusCustomerParty() *org.Party {
	return &org.Party{
		Name: "ACHETEUR",
		Identities: []*org.Identity{{
			Code: "200000008",
			Ext:  tax.MakeExtensions().Set(iso.ExtKeySchemeID, identitySchemeIDSIREN),
		}},
		Inboxes: []*org.Inbox{{Scheme: "0225", Code: "200000008_PEP"}},
		Ext:     tax.MakeExtensions().Set(ExtKeyRole, RoleBuyer),
	}
}

func testStatus(t *testing.T) *bill.Status {
	t.Helper()
	issued := cal.MakeDate(2026, 2, 1)
	return &bill.Status{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		IssueDate: cal.MakeDate(2026, 2, 2),
		Code:      "STA-2026-0001",
		Supplier:  statusSupplierParty(),
		Customer:  statusCustomerParty(),
		Lines: []*bill.StatusLine{
			{
				Key:  bill.StatusLineAccepted,
				Date: &issued,
				Doc: &org.DocumentRef{
					Code:      "INV-2026-001",
					IssueDate: &issued,
				},
			},
		},
	}
}

// --- bill.Status validation ----------------------------------------------

func TestStatusHappyPath(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
	assert.Equal(t, bill.StatusTypeResponse, st.Type)
}

func TestStatusRejectsSTCIdentityScheme(t *testing.T) {
	st := testStatus(t)
	// Add an STC (0231) identity on the supplier — admissible on a
	// Flow 2 invoice but not on a Flow 6 CDV.
	st.Supplier.Identities = append(st.Supplier.Identities, &org.Identity{
		Code: "12345678",
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			iso.ExtKeySchemeID: "0231",
		}),
	})
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "Flow 6 allowed schemes")
}

func TestStatusRejectsSystemType(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	st.Type = bill.StatusTypeSystem
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "type must be one of")
}

func TestStatusSupplierSIRENRequired(t *testing.T) {
	st := testStatus(t)
	st.Supplier.Identities = nil
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "SIREN")
}

func TestStatusKeyFilledFromStatusCodeExt(t *testing.T) {
	st := testStatus(t)
	st.Type = ""
	st.Lines[0].Key = ""
	st.Lines[0].Ext = st.Lines[0].Ext.Set(ExtKeyStatus, "205")
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
	assert.Equal(t, bill.StatusLineAccepted, st.Lines[0].Key)
	assert.Equal(t, bill.StatusTypeResponse, st.Type)
}

func TestStatusTypeMismatchRejected(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	st.Type = bill.StatusTypeUpdate // accepted is a response code
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "consistent with status type 'update'")
}

func TestStatusRejectsMultipleLines(t *testing.T) {
	st := testStatus(t)
	issued := cal.MakeDate(2026, 2, 1)
	st.Lines = append(st.Lines, &bill.StatusLine{
		Key:  bill.StatusLineAccepted,
		Date: &issued,
		Doc: &org.DocumentRef{
			Code:      "INV-2026-002",
			IssueDate: &issued,
		},
	})
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "status lines must contain exactly one entry")
}

func TestStatusRejectsZeroLines(t *testing.T) {
	st := testStatus(t)
	st.Lines = nil
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "status lines must contain exactly one entry")
}

// --- StatusLine validation -----------------------------------------------

func TestStatusLineUnknownKeyRejected(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = cbc.Key("made-up")
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "recognised Flow 6 event")
}

func TestStatusLineEmptyKeyRejected(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = ""
	err := rules.Validate(st)
	assert.Error(t, err)
}

func TestStatusLineDocCodeRequired(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Doc.Code = ""
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "status line doc code is required")
}

func TestStatusLineDocIssueDateRequired(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Doc.IssueDate = nil
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "status line doc issue_date is required")
}

// --- BR-FR-CDV-15: reason required on rejection-like statuses -----------

func TestStatusRejectedRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusLineRejected
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "reasons require at least one entry")
}

func TestStatusSuspendedRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusLineQuerying
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "reasons require at least one entry")
}

func TestStatusErrorRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusLineError
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "reasons require at least one entry")
}

func TestStatusAcceptedDoesNotRequireReason(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
}

// --- bill.Reason validation + normalization ------------------------------

func TestReasonNormalizerFillsKeyFromExt(t *testing.T) {
	r := &bill.Reason{
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: "QTE_ERR"}),
	}
	runNormalize(t, r)
	assert.Equal(t, bill.ReasonKeyQuantity, r.Key)
}

// TestReasonKeyFromEachReasonCode exercises prepareReasonKey across
// every CDAR ReasonCode bucket (one representative code each), so the
// reverse mapping and the matching forward bucket are both covered.
func TestReasonKeyFromEachReasonCode(t *testing.T) {
	cases := map[cbc.Code]cbc.Key{
		"COORD_BANC_ERR": bill.ReasonKeyFinanceTerms,
		"AUTRE":          bill.ReasonKeyOther,
		"NON_CONFORME":   bill.ReasonKeyLegal,
		"DOUBLON":        bill.ReasonKeyNotRecognized,
		"DEST_INC":       bill.ReasonKeyUnknownReceiver,
		"CMD_ERR":        bill.ReasonKeyReferences,
		"PU_ERR":         bill.ReasonKeyPrices,
		"QTE_ERR":        bill.ReasonKeyQuantity,
		"ART_ERR":        bill.ReasonKeyItems,
		"MODPAI_ERR":     bill.ReasonKeyPaymentTerms,
		"QUALITE_ERR":    bill.ReasonKeyQuality,
		"LIVR_INCOMP":    bill.ReasonKeyDelivery,
	}
	for code, want := range cases {
		t.Run(string(code), func(t *testing.T) {
			r := &bill.Reason{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: code})}
			runNormalize(t, r)
			assert.Equal(t, want, r.Key)
		})
	}
}

func TestReasonNormalizerFillsExtFromKey(t *testing.T) {
	r := &bill.Reason{Key: bill.ReasonKeyItems}
	runNormalize(t, r)
	assert.Equal(t, "ART_ERR", r.Ext.Get(ExtKeyReason).String())
}

func TestReasonNormalizerLeavesBothWhenSet(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: "ART_ERR"}),
	}
	runNormalize(t, r)
	assert.Equal(t, bill.ReasonKeyItems, r.Key)
	assert.Equal(t, "ART_ERR", r.Ext.Get(ExtKeyReason).String())
}

func TestReasonNormalizerLeavesUnknownExtAlone(t *testing.T) {
	r := &bill.Reason{
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: "NOPE"}),
	}
	runNormalize(t, r)
	assert.Equal(t, cbc.Key(""), r.Key)
}

// SetOneOf semantics in the normalizer: a CDAR ReasonCode that does
// not belong to the Reason.Key bucket is replaced with the bucket's
// default rather than failing validation. Mirrors how verifactu
// normalizes inconsistent tax-combo ext values.
func TestReasonNormalizerReplacesInconsistentExt(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: "QTE_ERR"}),
	}
	runNormalize(t, r)
	assert.Equal(t, cbc.Code("ART_ERR"), r.Ext.Get(ExtKeyReason))
}

func TestReasonNormalizerReplacesUnknownExt(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: "NOPE"}),
	}
	runNormalize(t, r)
	assert.Equal(t, cbc.Code("ART_ERR"), r.Ext.Get(ExtKeyReason))
}

// A caller can pick a non-default code from the bucket and the
// normalizer leaves it alone.
func TestReasonNormalizerPreservesNonDefaultCode(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyPrices,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReason: "CALCUL_ERR"}),
	}
	runNormalize(t, r)
	assert.Equal(t, cbc.Code("CALCUL_ERR"), r.Ext.Get(ExtKeyReason))
}

// --- Reason.Ext[fr-ctc-flow6-condition] → CDAR MDT-207 ------------------

// A rejected status with two sibling Reasons — one flagged DIV
// (invalid value), one flagged DVA (expected value) — passes
// validation. The CDAR cardinality (0..1 TypeCode per
// SpecifiedDocumentStatus) is honoured by spreading the two
// characteristics across separate Reasons. The accompanying
// bill.Condition entries are reserved for Peppol cac:Condition-style
// business-rule codes.
func TestStatusRejectedSiblingInvalidAndExpected(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusLineRejected
	st.Lines[0].Reasons = []*bill.Reason{
		{
			Key: bill.ReasonKeyLegal,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyReason:    "TX_TVA_ERR",
				ExtKeyCondition: ConditionInvalidData,
			}),
		},
		{
			Key: bill.ReasonKeyLegal,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyReason:    "TX_TVA_ERR",
				ExtKeyCondition: ConditionExpectedData,
			}),
		},
	}
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
}

// An unknown fr-ctc-flow6-condition value is rejected by
// tax.ExtensionHasValidCode using the registered Values list.
func TestReasonRejectsUnknownConditionExt(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyLegal,
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			ExtKeyReason:    "TX_TVA_ERR",
			ExtKeyCondition: "BOGUS",
		}),
	}
	err := rules.Validate(r, addonContext())
	assert.ErrorContains(t, err, "fr-ctc-flow6-condition")
}

// Every Status-applicable MDT-207 code is accepted on a bill.Reason.
// The 3 payment-amount markers (MEN / MPA / RAP) are explicitly
// excluded — they live on bill.Payment, not bill.Reason — and the
// rule rejects them; see TestReasonRejectsPaymentConditionCodes.
func TestReasonAcceptsAllStatusConditionCodes(t *testing.T) {
	for _, code := range []cbc.Code{
		ConditionBankDetailsUpdate, ConditionInvalidData,
		ConditionExpectedData, ConditionReplacementData,
		ConditionAmountApprovedHT, ConditionAmountApprovedTTC,
		ConditionAmountRejectedHT, ConditionAmountRejectedTTC,
		ConditionDiscount, ConditionRebate, ConditionReduction,
	} {
		r := &bill.Reason{
			Key: bill.ReasonKeyLegal,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyReason:    "TX_TVA_ERR",
				ExtKeyCondition: code,
			}),
		}
		assert.NoError(t, rules.Validate(r, addonContext()), "code %s", code)
	}
}

// Payment-related ProcessConditionCodes (211, 212) are rejected on a
// bill.Status — those belong on bill.Payment.
func TestStatusRejectsPaymentProcessCodes(t *testing.T) {
	for _, code := range []cbc.Code{"211", "212"} {
		st := testStatus(t)
		runNormalize(t, st)
		st.Ext = st.Ext.Set(ExtKeyStatus, code)
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "Status-applicable", "code %s", code)
	}
}

// --- fr-ctc-flow6-action (MDT-121 / BR-FR-CDV-CL-10) -------------------

// Action normalizer fills the ext from the Key bucket.
func TestActionNormalizerFillsExtFromKey(t *testing.T) {
	a := &bill.Action{Key: bill.ActionKeyReissue}
	runNormalize(t, a)
	assert.Equal(t, cbc.Code("NIN"), a.Ext.Get(ExtKeyAction))
}

// Action normalizer reverse-maps the Key from the ext (round-trip).
func TestActionNormalizerFillsKeyFromExt(t *testing.T) {
	a := &bill.Action{
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyAction: "CNP"}),
	}
	runNormalize(t, a)
	assert.Equal(t, bill.ActionKeyCreditPartial, a.Key)
	assert.Equal(t, cbc.Code("CNP"), a.Ext.Get(ExtKeyAction))
}

// Every MDT-121 code defined in actionTable is accepted by the rule.
// The normalizer fills the Action.Key from the ext so core
// bill.Action validation (which requires Key) also passes.
func TestActionAcceptsAllMDT121Codes(t *testing.T) {
	for _, code := range []cbc.Code{"NOA", "PIN", "NIN", "CNF", "CNP", "CNA", "OTH"} {
		a := &bill.Action{
			Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyAction: code}),
		}
		runNormalize(t, a)
		assert.NoError(t, rules.Validate(a, addonContext()), "code %s", code)
	}
}

// Unknown action codes are rejected by tax.ExtensionHasValidCode using
// the registered Values list. We pin a valid Key so the failure
// surfaces from the flow6 rule, not the core "action key required".
func TestActionRejectsUnknownCode(t *testing.T) {
	a := &bill.Action{
		Key: bill.ActionKeyOther,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyAction: "BOGUS"}),
	}
	err := rules.Validate(a, addonContext())
	assert.ErrorContains(t, err, "fr-ctc-flow6-action")
}

// Payment-amount codes (MEN / MPA / RAP) are rejected on a Reason.
func TestReasonRejectsPaymentConditionCodes(t *testing.T) {
	for _, code := range []cbc.Code{
		ConditionAmountReceived, ConditionAmountPaid, ConditionAmountRemaining,
	} {
		r := &bill.Reason{
			Key: bill.ReasonKeyLegal,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyReason:    "TX_TVA_ERR",
				ExtKeyCondition: code,
			}),
		}
		err := rules.Validate(r, addonContext())
		assert.ErrorContains(t, err, "belong on bill.Payment", "code %s", code)
	}
}

// The normalizer defaults Reason.Ext[fr-ctc-flow6-condition] per
// bucket: finance-terms → CBB, everything else → DIV.
func TestReasonNormalizerDefaultsConditionExt(t *testing.T) {
	cases := []struct {
		key  cbc.Key
		want cbc.Code
	}{
		{bill.ReasonKeyFinanceTerms, ConditionBankDetailsUpdate},
		{bill.ReasonKeyLegal, ConditionInvalidData},
		{bill.ReasonKeyPrices, ConditionInvalidData},
		{bill.ReasonKeyItems, ConditionInvalidData},
		{bill.ReasonKeyQuality, ConditionInvalidData},
		{bill.ReasonKeyDelivery, ConditionInvalidData},
		{bill.ReasonKeyOther, ConditionInvalidData},
	}
	for _, c := range cases {
		r := &bill.Reason{Key: c.key}
		runNormalize(t, r)
		assert.Equal(t, c.want, r.Ext.Get(ExtKeyCondition), "key %s", c.key)
	}
}

// Caller-supplied DVA / MAJ survive normalization (SetOneOf keeps any
// value that is already one of the bucket's allowed codes).
func TestReasonNormalizerPreservesExplicitConditionExt(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyPrices,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyCondition: ConditionExpectedData}),
	}
	runNormalize(t, r)
	assert.Equal(t, ConditionExpectedData, r.Ext.Get(ExtKeyCondition))
}

// none / partial / empty keys carry no characteristic context, so the
// normalizer leaves the ext blank.
func TestReasonNormalizerSkipsConditionExtForNonCorrectiveKeys(t *testing.T) {
	for _, k := range []cbc.Key{bill.ReasonKeyNone, bill.ReasonKeyPartial, ""} {
		r := &bill.Reason{Key: k}
		runNormalize(t, r)
		assert.Empty(t, r.Ext.Get(ExtKeyCondition).String(), "key %s", k)
	}
}

// --- Internal helper coverage (nil / wrong-type defensive branches) -----

func TestNormalizeStatusNilSafe(t *testing.T) {
	assert.NotPanics(t, func() { normalizeStatus(nil) })
}

func TestNormalizeStatusAllLinesNil(t *testing.T) {
	st := &bill.Status{Lines: []*bill.StatusLine{nil}}
	normalizeStatus(st)
	assert.Equal(t, cbc.Key(""), st.Type)
}

func TestNormalizeReasonNilSafe(t *testing.T) {
	assert.NotPanics(t, func() { normalizeReason(nil) })
}

// --- defensive coverage: nil / wrong-type / empty-slice guards --------

func TestSetPartyRoleDefaultNilParty(t *testing.T) {
	assert.NotPanics(t, func() { setPartyRoleDefault(nil, RoleSeller) })
}

func TestSetPartyRoleDefaultExistingNotOverridden(t *testing.T) {
	p := &org.Party{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RoleBuyer})}
	setPartyRoleDefault(p, RoleSeller)
	assert.Equal(t, RoleBuyer, p.Ext.Get(ExtKeyRole))
}

func TestPartyHasInboxWhenRequiredWrongType(t *testing.T) {
	assert.True(t, partyHasInboxWhenRequired("x"))
}

func TestPartyHasInboxWhenRequiredWKRole(t *testing.T) {
	p := &org.Party{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RolePlatform})}
	assert.True(t, partyHasInboxWhenRequired(p))
}
