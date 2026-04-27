package flow6

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Helpers --------------------------------------------------------------

// addonContext activates the Flow 6 rule guard so the addon's validators
// fire even for standalone objects (bill.Reason / org.Party) that do not
// carry an addon themselves.
func addonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(V1), tax.AddonForKey(V1))
	}
}

// runNormalize invokes the addon's registered normalizer on the given
// object, matching what tax.Normalize would do during Calculate.
func runNormalize(t *testing.T, doc any) {
	t.Helper()
	tax.Normalize([]tax.Normalizer{tax.AddonForKey(V1).Normalizer}, doc)
}

func frPartyWithSIREN() *org.Party {
	return &org.Party{
		Name: "Test Platform SARL",
		Identities: []*org.Identity{
			{
				Code: "356000000",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					iso.ExtKeySchemeID: "0002",
				}),
			},
		},
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
		Supplier:  frPartyWithSIREN(),
		Lines: []*bill.StatusLine{
			{
				Key:  bill.StatusEventAccepted,
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
	assert.Equal(t, RoleSE, st.Supplier.Ext.Get(ExtKeyRole))
}

func TestStatusRejectsSystemType(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	st.Type = bill.StatusTypeSystem
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "status type must be one of")
}

func TestStatusSupplierSIRENRequired(t *testing.T) {
	st := testStatus(t)
	st.Supplier.Identities = nil
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "SIREN")
}

func TestStatusTypeMismatchRejected(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	st.Type = bill.StatusTypeUpdate // accepted is a response code
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "Status.Type must match")
}

func TestStatusRejectsMultipleLines(t *testing.T) {
	st := testStatus(t)
	issued := cal.MakeDate(2026, 2, 1)
	st.Lines = append(st.Lines, &bill.StatusLine{
		Key:  bill.StatusEventAccepted,
		Date: &issued,
		Doc: &org.DocumentRef{
			Code:      "INV-2026-002",
			IssueDate: &issued,
		},
	})
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "exactly one status line")
}

func TestStatusRejectsZeroLines(t *testing.T) {
	st := testStatus(t)
	st.Lines = nil
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "exactly one status line")
}

func TestStatusHasExactlyOneLineWrongType(t *testing.T) {
	assert.False(t, statusHasExactlyOneLine("x"))
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
	assert.ErrorContains(t, err, "invoice code is required")
}

func TestStatusLineDocIssueDateRequired(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Doc.IssueDate = nil
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "invoice issue date is required")
}

// --- BR-FR-CDV-15: reason required on rejection-like statuses -----------

func TestStatusRejectedRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventRejected
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "require at least one reason")
}

func TestStatusDisputedRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = StatusEventDisputed
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "require at least one reason")
}

func TestStatusSuspendedRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = StatusEventSuspended
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "require at least one reason")
}

func TestStatusPartiallyAcceptedRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = StatusEventPartiallyAccepted
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "require at least one reason")
}

func TestStatusErrorRequiresReason(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventError
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "require at least one reason")
}

func TestStatusAcceptedDoesNotRequireReason(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
}

// --- Paid: MEN Characteristic required -----------------------------------

func TestStatusPaidRequiresAmount(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventPaid
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "MEN")
}

func TestStatusPaidSatisfiedByComplement(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventPaid
	obj, err := schema.NewObject(&Characteristic{
		TypeCode: TypeCodeAmountReceived,
		Amount: &currency.Amount{
			Currency: "EUR",
			Value:    num.MakeAmount(125000, 2),
		},
	})
	require.NoError(t, err)
	st.Lines[0].Complements = []*schema.Object{obj}
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
}

func TestStatusPaidWithoutMENFailsEvenWithOtherTypes(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventPaid
	obj, err := schema.NewObject(&Characteristic{
		TypeCode: TypeCodeAmountPaid,
		Amount:   &currency.Amount{Currency: "EUR", Value: num.MakeAmount(100, 0)},
	})
	require.NoError(t, err)
	st.Lines[0].Complements = []*schema.Object{obj}
	runNormalize(t, st)
	err = rules.Validate(st)
	assert.ErrorContains(t, err, "MEN")
}

func TestStatusPaidMENMissingCurrencyFails(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventPaid
	obj, err := schema.NewObject(&Characteristic{
		TypeCode: TypeCodeAmountReceived,
		Amount:   &currency.Amount{Value: num.MakeAmount(100, 0)},
	})
	require.NoError(t, err)
	st.Lines[0].Complements = []*schema.Object{obj}
	runNormalize(t, st)
	err = rules.Validate(st)
	assert.ErrorContains(t, err, "MEN")
}

// --- MDT-207 TypeCode whitelist ------------------------------------------

func TestStatusCharacteristicUnknownTypeCodeRejected(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventPaid
	obj, err := schema.NewObject(&Characteristic{
		TypeCode: "BOGUS",
		Amount:   &currency.Amount{Currency: "EUR", Value: num.MakeAmount(100, 0)},
	})
	require.NoError(t, err)
	st.Lines[0].Complements = []*schema.Object{obj}
	runNormalize(t, st)
	err = rules.Validate(st)
	assert.ErrorContains(t, err, "MDT-207")
}

// --- Characteristic ReasonCode link --------------------------------------

func TestStatusCharacteristicReasonLinkMismatch(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventRejected
	st.Lines[0].Reasons = []*bill.Reason{{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(tax.ExtMap{ExtKeyReasonCode: "ART_ERR"}),
	}}
	obj, err := schema.NewObject(&Characteristic{
		ReasonCode: "QTE_ERR", // not matching any sibling reason
		Name:       "description",
		Value:      "wrong",
	})
	require.NoError(t, err)
	st.Lines[0].Complements = []*schema.Object{obj}
	runNormalize(t, st)
	err = rules.Validate(st)
	assert.ErrorContains(t, err, "ReasonCode must match")
}

func TestStatusCharacteristicReasonLinkMatch(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventRejected
	st.Lines[0].Reasons = []*bill.Reason{{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(tax.ExtMap{ExtKeyReasonCode: "ART_ERR"}),
	}}
	obj, err := schema.NewObject(&Characteristic{
		ReasonCode: "ART_ERR",
		Name:       "description",
		Value:      "corrected",
	})
	require.NoError(t, err)
	st.Lines[0].Complements = []*schema.Object{obj}
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
}

// --- bill.Reason validation + normalization ------------------------------

func TestReasonNormalizerFillsKeyFromExt(t *testing.T) {
	r := &bill.Reason{
		Ext: tax.ExtensionsOf(tax.ExtMap{ExtKeyReasonCode: "QTE_ERR"}),
	}
	runNormalize(t, r)
	assert.Equal(t, bill.ReasonKeyQuantity, r.Key)
}

func TestReasonNormalizerFillsExtFromKey(t *testing.T) {
	r := &bill.Reason{Key: bill.ReasonKeyItems}
	runNormalize(t, r)
	assert.Equal(t, "ART_ERR", r.Ext.Get(ExtKeyReasonCode).String())
}

func TestReasonNormalizerLeavesBothWhenSet(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(tax.ExtMap{ExtKeyReasonCode: "ART_ERR"}),
	}
	runNormalize(t, r)
	assert.Equal(t, bill.ReasonKeyItems, r.Key)
	assert.Equal(t, "ART_ERR", r.Ext.Get(ExtKeyReasonCode).String())
}

func TestReasonNormalizerLeavesUnknownExtAlone(t *testing.T) {
	r := &bill.Reason{
		Ext: tax.ExtensionsOf(tax.ExtMap{ExtKeyReasonCode: "NOPE"}),
	}
	runNormalize(t, r)
	assert.Equal(t, cbc.Key(""), r.Key)
}

func TestReasonRulesRejectInconsistentExt(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(tax.ExtMap{ExtKeyReasonCode: "QTE_ERR"}),
	}
	err := rules.Validate(r, addonContext())
	assert.ErrorContains(t, err, "must match reason.key")
}

func TestReasonExtUnknownCodeRejected(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(tax.ExtMap{ExtKeyReasonCode: "NOPE"}),
	}
	err := rules.Validate(r, addonContext())
	assert.ErrorContains(t, err, "must match reason.key")
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

func TestPartyHasSIRENIdentityWrongType(t *testing.T) {
	assert.False(t, partyHasSIRENIdentity("not a party"))
}

func TestPartyHasSIRENIdentityNilParty(t *testing.T) {
	assert.False(t, partyHasSIRENIdentity((*org.Party)(nil)))
}

func TestPartyHasSIRENIdentityWithoutExt(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "X"}}}
	assert.False(t, partyHasSIRENIdentity(p))
}

func TestStatusLineKeyKnownWrongType(t *testing.T) {
	assert.False(t, statusLineKeyKnown("x"))
}

func TestStatusLinePaidHasAmountWrongType(t *testing.T) {
	assert.True(t, statusLinePaidHasAmount(42))
}

func TestStatusLinePaidHasAmountNonPaidLine(t *testing.T) {
	assert.True(t, statusLinePaidHasAmount(&bill.StatusLine{Key: bill.StatusEventAccepted}))
}

func TestStatusLineTypeCodesKnownWrongType(t *testing.T) {
	assert.True(t, statusLineTypeCodesKnown("x"))
}

func TestStatusLineTypeCodesKnownEmptyLine(t *testing.T) {
	assert.True(t, statusLineTypeCodesKnown(&bill.StatusLine{}))
}

func TestStatusLineReasonLinksResolveWrongType(t *testing.T) {
	assert.True(t, statusLineReasonLinksResolve("x"))
}

func TestStatusLineReasonLinksResolveEmptyComplements(t *testing.T) {
	assert.True(t, statusLineReasonLinksResolve(&bill.StatusLine{}))
}

func TestStatusLineRequiresReasonWrongType(t *testing.T) {
	assert.True(t, statusLineRequiresReason("x"))
}

func TestStatusTypeMatchesLinesWrongType(t *testing.T) {
	assert.True(t, statusTypeMatchesLines("x"))
}

func TestStatusTypeMatchesLinesUnknownLineKey(t *testing.T) {
	st := &bill.Status{
		Type:  bill.StatusTypeResponse,
		Lines: []*bill.StatusLine{{Key: "unknown"}},
	}
	assert.True(t, statusTypeMatchesLines(st))
}

func TestLineHasReasonCodeNilReason(t *testing.T) {
	line := &bill.StatusLine{Reasons: []*bill.Reason{nil}}
	assert.False(t, lineHasReasonCode(line, "ART_ERR"))
}

func TestReasonExtMatchesKeyWrongType(t *testing.T) {
	assert.True(t, reasonExtMatchesKey("x"))
}
