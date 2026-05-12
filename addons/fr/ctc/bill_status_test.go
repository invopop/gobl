package ctc

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

// statusSupplierParty returns a French supplier party with a SIREN
// identity carried via the iso-scheme-id extension. Used as the
// document-level Supplier on a bill.Status.
func statusSupplierParty() *org.Party {
	return &org.Party{
		Name: "Test Platform SARL",
		Identities: []*org.Identity{
			{
				Code: "356000000",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			},
		},
	}
}

// issuerParty returns a buyer-side Issuer (BR-FR-CDV-CL-03 allowed:
// BY/AB/DL/SE/SR/PE/PR/II/IV) with a SIREN identity and inbox so
// BR-FR-CDV-08 is satisfied.
func issuerParty() *org.Party {
	return &org.Party{
		Name: "ACHETEUR",
		Identities: []*org.Identity{{
			Code: "200000008",
			Ext:  tax.MakeExtensions().Set(iso.ExtKeySchemeID, identitySchemeIDSIREN),
		}},
		Inboxes: []*org.Inbox{{Scheme: "0225", Code: "200000008_PEP"}},
		Ext:     tax.MakeExtensions().Set(ExtKeyRole, RoleBY),
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
		Ext:     tax.MakeExtensions().Set(ExtKeyRole, RoleBY),
	}
}

// recipientParty returns the seller-end Recipient counterpart with an
// inbox so BR-FR-CDV-08 is satisfied. Carries the same SIREN as the
// document-level Supplier — both represent the same seller legal
// entity, just at different endpoints — so the normaliser's
// ensureSIRENOnSupplier no-ops instead of duplicating the identity.
func recipientParty() *org.Party {
	return &org.Party{
		Name: "VENDEUR",
		Identities: []*org.Identity{{
			Code: "356000000",
			Ext:  tax.MakeExtensions().Set(iso.ExtKeySchemeID, identitySchemeIDSIREN),
		}},
		Inboxes: []*org.Inbox{{Scheme: "0225", Code: "356000000_PEP"}},
		Ext:     tax.MakeExtensions().Set(ExtKeyRole, RoleSE),
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
		Issuer:    issuerParty(),
		Recipient: recipientParty(),
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
	assert.ErrorContains(t, err, "Flow 6 allow-list")
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
	// Strip the SE party's identity too so the normaliser cannot
	// auto-populate Supplier from it.
	st.Recipient.Identities = nil
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "SIREN")
}

func TestStatusSupplierSIRENFilledFromSEParty(t *testing.T) {
	st := testStatus(t)
	st.Supplier = nil // recipient is SE-roled with SIREN 356000000
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
	require.NotNil(t, st.Supplier)
	require.Len(t, st.Supplier.Identities, 1)
	assert.Equal(t, cbc.Code("356000000"), st.Supplier.Identities[0].Code)
	assert.Equal(t, identitySchemeIDSIREN,
		st.Supplier.Identities[0].Ext.Get(iso.ExtKeySchemeID).String())
}

func TestStatusKeyFilledFromStatusCodeExt(t *testing.T) {
	st := testStatus(t)
	st.Type = ""
	st.Lines[0].Key = ""
	st.Ext = st.Ext.Set(ExtKeyStatusCode, "205")
	runNormalize(t, st)
	require.NoError(t, rules.Validate(st))
	assert.Equal(t, bill.StatusEventAccepted, st.Lines[0].Key)
	assert.Equal(t, bill.StatusTypeResponse, st.Type)
}

func TestStatusTypeMismatchRejected(t *testing.T) {
	st := testStatus(t)
	runNormalize(t, st)
	st.Type = bill.StatusTypeUpdate // accepted is a response code
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "Status.Type must be a valid pair")
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
	st.Lines[0].Key = bill.StatusEventQuerying
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
	st.Type = bill.StatusTypeResponse
	runNormalize(t, st)
	err := rules.Validate(st)
	assert.ErrorContains(t, err, "MEN")
}

func TestStatusPaidSatisfiedByComplement(t *testing.T) {
	st := testStatus(t)
	st.Lines[0].Key = bill.StatusEventPaid
	st.Type = bill.StatusTypeResponse
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
	st.Type = bill.StatusTypeResponse
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
	st.Type = bill.StatusTypeResponse
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
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReasonCode: "TX_TVA_ERR"}),
	}}
	obj, err := schema.NewObject(&Characteristic{
		ReasonCode: "QTE_ERR",
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
		Key: bill.ReasonKeyLegal,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReasonCode: "TX_TVA_ERR"}),
	}}
	obj, err := schema.NewObject(&Characteristic{
		ReasonCode: "TX_TVA_ERR",
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
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReasonCode: "QTE_ERR"}),
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
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReasonCode: "ART_ERR"}),
	}
	runNormalize(t, r)
	assert.Equal(t, bill.ReasonKeyItems, r.Key)
	assert.Equal(t, "ART_ERR", r.Ext.Get(ExtKeyReasonCode).String())
}

func TestReasonNormalizerLeavesUnknownExtAlone(t *testing.T) {
	r := &bill.Reason{
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReasonCode: "NOPE"}),
	}
	runNormalize(t, r)
	assert.Equal(t, cbc.Key(""), r.Key)
}

func TestReasonRulesRejectInconsistentExt(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReasonCode: "QTE_ERR"}),
	}
	err := rules.Validate(r, addonContext())
	assert.ErrorContains(t, err, "must match reason.key")
}

func TestReasonExtUnknownCodeRejected(t *testing.T) {
	r := &bill.Reason{
		Key: bill.ReasonKeyItems,
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyReasonCode: "NOPE"}),
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

func TestStatusPartyHasSIRENIdentityWrongType(t *testing.T) {
	assert.False(t, statusPartyHasSIRENIdentity("not a party"))
}

func TestStatusPartyHasSIRENIdentityNilParty(t *testing.T) {
	assert.False(t, statusPartyHasSIRENIdentity((*org.Party)(nil)))
}

func TestStatusPartyHasSIRENIdentityWithoutExt(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "X"}}}
	assert.False(t, statusPartyHasSIRENIdentity(p))
}

func TestStatusLineKeyKnownWrongType(t *testing.T) {
	assert.False(t, statusLineKeyKnown("x"))
}

func TestStatusPaidResponseHasAmountWrongType(t *testing.T) {
	assert.True(t, statusPaidResponseHasAmount(42))
}

func TestStatusPaidResponseHasAmountNonPaidLine(t *testing.T) {
	st := &bill.Status{
		Type:  bill.StatusTypeResponse,
		Lines: []*bill.StatusLine{{Key: bill.StatusEventAccepted}},
	}
	assert.True(t, statusPaidResponseHasAmount(st))
}

func TestStatusPaidResponseHasAmountUpdateSkips(t *testing.T) {
	st := &bill.Status{
		Type:  bill.StatusTypeUpdate,
		Lines: []*bill.StatusLine{{Key: bill.StatusEventPaid}},
	}
	assert.True(t, statusPaidResponseHasAmount(st))
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

// --- defensive coverage: nil / wrong-type / empty-slice guards --------

func TestSetPartyRoleDefaultNilParty(t *testing.T) {
	assert.NotPanics(t, func() { setPartyRoleDefault(nil, RoleSE) })
}

func TestSetPartyRoleDefaultExistingNotOverridden(t *testing.T) {
	p := &org.Party{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RoleBY})}
	setPartyRoleDefault(p, RoleSE)
	assert.Equal(t, RoleBY, p.Ext.Get(ExtKeyRole))
}

func TestPartyHasRoleWrongType(t *testing.T) {
	assert.False(t, partyHasRole("x"))
}

func TestPartyHasRoleEmptyExt(t *testing.T) {
	assert.False(t, partyHasRole(&org.Party{}))
}

func TestPartyHasInboxWhenRequiredWrongType(t *testing.T) {
	assert.True(t, partyHasInboxWhenRequired("x"))
}

func TestPartyHasInboxWhenRequiredWKRole(t *testing.T) {
	p := &org.Party{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RoleWK})}
	assert.True(t, partyHasInboxWhenRequired(p))
}

func TestStatusPartiesIdentitySchemesAllowedWrongType(t *testing.T) {
	assert.True(t, statusPartiesIdentitySchemesAllowed("x"))
}

func TestStatusReasonCodesAllowedWrongType(t *testing.T) {
	assert.True(t, statusReasonCodesAllowed("x"))
}

func TestStatusReasonCodesAllowedNilReason(t *testing.T) {
	st := &bill.Status{
		Type: bill.StatusTypeResponse,
		Lines: []*bill.StatusLine{{
			Key:     bill.StatusEventRejected,
			Reasons: []*bill.Reason{nil},
		}},
	}
	assert.True(t, statusReasonCodesAllowed(st))
}

// TestEnsureSIRENOnSupplierAlreadyCarries covers the "supplier already
// carries the SIREN" early-return path that the happy-path tests don't
// reach (since the test fixture aligns supplier and recipient SIRENs).
func TestEnsureSIRENOnSupplierAlreadyCarries(t *testing.T) {
	siren := &org.Identity{
		Code: "356000000",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{"iso-scheme-id": "0002"}),
	}
	p := &org.Party{Identities: []*org.Identity{
		{Code: "356000000", Ext: tax.ExtensionsOf(cbc.CodeMap{"iso-scheme-id": "0002"})},
	}}
	got := ensureSIRENOnSupplier(p, siren)
	assert.Same(t, p, got)
	assert.Len(t, got.Identities, 1)
}
