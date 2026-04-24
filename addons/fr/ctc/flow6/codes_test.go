package flow6

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

// assertProcessRoundTrip verifies that a CDAR ProcessConditionCode
// resolves to the expected (key, type) pair and that the pair resolves
// back to the same code.
func assertProcessRoundTrip(t *testing.T, code string, wantKey, wantType cbc.Key) {
	t.Helper()
	key, typ, ok := StatusKeyFor(code)
	assert.True(t, ok, "StatusKeyFor should resolve")
	assert.Equal(t, wantKey, key)
	assert.Equal(t, wantType, typ)
	got, ok := CDARProcessCodeFor(key, typ)
	assert.True(t, ok, "CDARProcessCodeFor should resolve")
	assert.Equal(t, code, got)
}

func TestProcessCode200Issued(t *testing.T) {
	assertProcessRoundTrip(t, "200", bill.StatusEventIssued, bill.StatusTypeUpdate)
}

func TestProcessCode201IssuedByPlatform(t *testing.T) {
	assertProcessRoundTrip(t, "201", StatusEventIssuedByPlatform, bill.StatusTypeUpdate)
}

func TestProcessCode202ReceivedByPlatform(t *testing.T) {
	assertProcessRoundTrip(t, "202", StatusEventReceivedByPlatform, bill.StatusTypeResponse)
}

func TestProcessCode203MadeAvailable(t *testing.T) {
	assertProcessRoundTrip(t, "203", StatusEventMadeAvailable, bill.StatusTypeResponse)
}

func TestProcessCode204Processing(t *testing.T) {
	assertProcessRoundTrip(t, "204", bill.StatusEventProcessing, bill.StatusTypeResponse)
}

func TestProcessCode205Accepted(t *testing.T) {
	assertProcessRoundTrip(t, "205", bill.StatusEventAccepted, bill.StatusTypeResponse)
}

func TestProcessCode206PartiallyAccepted(t *testing.T) {
	assertProcessRoundTrip(t, "206", StatusEventPartiallyAccepted, bill.StatusTypeResponse)
}

func TestProcessCode207Disputed(t *testing.T) {
	assertProcessRoundTrip(t, "207", StatusEventDisputed, bill.StatusTypeResponse)
}

func TestProcessCode208Suspended(t *testing.T) {
	assertProcessRoundTrip(t, "208", StatusEventSuspended, bill.StatusTypeResponse)
}

func TestProcessCode209Completed(t *testing.T) {
	assertProcessRoundTrip(t, "209", StatusEventCompleted, bill.StatusTypeResponse)
}

func TestProcessCode210Rejected(t *testing.T) {
	assertProcessRoundTrip(t, "210", bill.StatusEventRejected, bill.StatusTypeResponse)
}

func TestProcessCode211PaymentForwarded(t *testing.T) {
	assertProcessRoundTrip(t, "211", StatusEventPaymentForwarded, bill.StatusTypeUpdate)
}

func TestProcessCode212Paid(t *testing.T) {
	assertProcessRoundTrip(t, "212", bill.StatusEventPaid, bill.StatusTypeResponse)
}

func TestProcessCode213Error(t *testing.T) {
	assertProcessRoundTrip(t, "213", bill.StatusEventError, bill.StatusTypeResponse)
}

func TestProcessCodeUnknownReturnsFalse(t *testing.T) {
	_, _, ok := StatusKeyFor("999")
	assert.False(t, ok)
}

func TestProcessKeyTypeMismatchReturnsFalse(t *testing.T) {
	// paid is a response code; querying with Type=update must miss.
	_, ok := CDARProcessCodeFor(bill.StatusEventPaid, bill.StatusTypeUpdate)
	assert.False(t, ok)
}

// assertActionRoundTrip verifies that an action code resolves to the
// expected bill.Action.Key and round-trips back.
func assertActionRoundTrip(t *testing.T, code string, wantKey cbc.Key) {
	t.Helper()
	key, ok := ActionKeyFor(code)
	assert.True(t, ok)
	assert.Equal(t, wantKey, key)
	got, ok := CDARActionCodeFor(key)
	assert.True(t, ok)
	assert.Equal(t, code, got)
}

func TestActionNOA(t *testing.T) { assertActionRoundTrip(t, "NOA", bill.ActionKeyNone) }
func TestActionPIN(t *testing.T) { assertActionRoundTrip(t, "PIN", bill.ActionKeyProvide) }
func TestActionNIN(t *testing.T) { assertActionRoundTrip(t, "NIN", bill.ActionKeyReissue) }
func TestActionCNF(t *testing.T) { assertActionRoundTrip(t, "CNF", bill.ActionKeyCreditFull) }
func TestActionCNP(t *testing.T) { assertActionRoundTrip(t, "CNP", bill.ActionKeyCreditPartial) }
func TestActionCNA(t *testing.T) { assertActionRoundTrip(t, "CNA", bill.ActionKeyCreditAmount) }
func TestActionOTH(t *testing.T) { assertActionRoundTrip(t, "OTH", bill.ActionKeyOther) }

func TestActionUnknownCodeMisses(t *testing.T) {
	_, ok := ActionKeyFor("XYZ")
	assert.False(t, ok)
}

func TestActionUnknownKeyMisses(t *testing.T) {
	_, ok := CDARActionCodeFor("never-heard-of")
	assert.False(t, ok)
}

// assertReasonBucket verifies that a CDAR reason code buckets into the
// expected bill.Reason.Key.
func assertReasonBucket(t *testing.T, code string, wantKey cbc.Key) {
	t.Helper()
	got, ok := ReasonKeyFor(code)
	assert.True(t, ok)
	assert.Equal(t, wantKey, got)
}

// Business-rejection reasons ---------------------------------------------

func TestReasonNON_TRANSMISE(t *testing.T) {
	assertReasonBucket(t, "NON_TRANSMISE", bill.ReasonKeyUnknownReceiver)
}
func TestReasonJUSTIF_ABS(t *testing.T) {
	assertReasonBucket(t, "JUSTIF_ABS", bill.ReasonKeyReferences)
}
func TestReasonROUTAGE_ERR(t *testing.T) {
	assertReasonBucket(t, "ROUTAGE_ERR", bill.ReasonKeyUnknownReceiver)
}
func TestReasonAUTRE(t *testing.T) {
	assertReasonBucket(t, "AUTRE", bill.ReasonKeyOther)
}
func TestReasonCOORD_BANC_ERR(t *testing.T) {
	assertReasonBucket(t, "COORD_BANC_ERR", bill.ReasonKeyFinanceTerms)
}
func TestReasonTX_TVA_ERR(t *testing.T) {
	assertReasonBucket(t, "TX_TVA_ERR", bill.ReasonKeyLegal)
}
func TestReasonMONTANTTOTAL_ERR(t *testing.T) {
	assertReasonBucket(t, "MONTANTTOTAL_ERR", bill.ReasonKeyPrices)
}
func TestReasonCALCUL_ERR(t *testing.T) {
	assertReasonBucket(t, "CALCUL_ERR", bill.ReasonKeyPrices)
}
func TestReasonNON_CONFORME(t *testing.T) {
	assertReasonBucket(t, "NON_CONFORME", bill.ReasonKeyLegal)
}
func TestReasonDOUBLON(t *testing.T) {
	assertReasonBucket(t, "DOUBLON", bill.ReasonKeyNotRecognized)
}
func TestReasonDEST_INC(t *testing.T) {
	assertReasonBucket(t, "DEST_INC", bill.ReasonKeyUnknownReceiver)
}
func TestReasonDEST_ERR(t *testing.T) {
	assertReasonBucket(t, "DEST_ERR", bill.ReasonKeyReferences)
}
func TestReasonTRANSAC_INC(t *testing.T) {
	assertReasonBucket(t, "TRANSAC_INC", bill.ReasonKeyNotRecognized)
}
func TestReasonEMMET_INC(t *testing.T) {
	assertReasonBucket(t, "EMMET_INC", bill.ReasonKeyNotRecognized)
}
func TestReasonCONTRAT_TERM(t *testing.T) {
	assertReasonBucket(t, "CONTRAT_TERM", bill.ReasonKeyNotRecognized)
}
func TestReasonDOUBLE_FACT(t *testing.T) {
	assertReasonBucket(t, "DOUBLE_FACT", bill.ReasonKeyNotRecognized)
}
func TestReasonCMD_ERR(t *testing.T) {
	assertReasonBucket(t, "CMD_ERR", bill.ReasonKeyReferences)
}
func TestReasonADR_ERR(t *testing.T) {
	assertReasonBucket(t, "ADR_ERR", bill.ReasonKeyReferences)
}
func TestReasonSIRET_ERR(t *testing.T) {
	assertReasonBucket(t, "SIRET_ERR", bill.ReasonKeyReferences)
}
func TestReasonCODE_ROUTAGE_ERR(t *testing.T) {
	assertReasonBucket(t, "CODE_ROUTAGE_ERR", bill.ReasonKeyReferences)
}
func TestReasonREF_CT_ABSENT(t *testing.T) {
	assertReasonBucket(t, "REF_CT_ABSENT", bill.ReasonKeyReferences)
}
func TestReasonREF_ERR(t *testing.T) {
	assertReasonBucket(t, "REF_ERR", bill.ReasonKeyReferences)
}
func TestReasonPU_ERR(t *testing.T) {
	assertReasonBucket(t, "PU_ERR", bill.ReasonKeyPrices)
}
func TestReasonREM_ERR(t *testing.T) {
	assertReasonBucket(t, "REM_ERR", bill.ReasonKeyPrices)
}
func TestReasonQTE_ERR(t *testing.T) {
	assertReasonBucket(t, "QTE_ERR", bill.ReasonKeyQuantity)
}
func TestReasonART_ERR(t *testing.T) {
	assertReasonBucket(t, "ART_ERR", bill.ReasonKeyItems)
}
func TestReasonMODPAI_ERR(t *testing.T) {
	assertReasonBucket(t, "MODPAI_ERR", bill.ReasonKeyPaymentTerms)
}
func TestReasonQUALITE_ERR(t *testing.T) {
	assertReasonBucket(t, "QUALITE_ERR", bill.ReasonKeyQuality)
}
func TestReasonLIVR_INCOMP(t *testing.T) {
	assertReasonBucket(t, "LIVR_INCOMP", bill.ReasonKeyDelivery)
}

// Technical / platform rejection reasons (code 213) ---------------------

func TestReasonREJ_SEMAN(t *testing.T) {
	assertReasonBucket(t, "REJ_SEMAN", bill.ReasonKeyLegal)
}
func TestReasonREJ_UNI(t *testing.T) {
	assertReasonBucket(t, "REJ_UNI", bill.ReasonKeyNotRecognized)
}
func TestReasonREJ_COH(t *testing.T) {
	assertReasonBucket(t, "REJ_COH", bill.ReasonKeyLegal)
}
func TestReasonREJ_ADR(t *testing.T) {
	assertReasonBucket(t, "REJ_ADR", bill.ReasonKeyReferences)
}
func TestReasonREJ_CONT_B2G(t *testing.T) {
	assertReasonBucket(t, "REJ_CONT_B2G", bill.ReasonKeyLegal)
}
func TestReasonREJ_REF_PJ(t *testing.T) {
	assertReasonBucket(t, "REJ_REF_PJ", bill.ReasonKeyReferences)
}
func TestReasonREJ_ASS_PJ(t *testing.T) {
	assertReasonBucket(t, "REJ_ASS_PJ", bill.ReasonKeyReferences)
}
func TestReasonIRR_VIDE_F(t *testing.T) {
	assertReasonBucket(t, "IRR_VIDE_F", bill.ReasonKeyLegal)
}
func TestReasonIRR_TYPE_F(t *testing.T) {
	assertReasonBucket(t, "IRR_TYPE_F", bill.ReasonKeyLegal)
}
func TestReasonIRR_SYNTAX(t *testing.T) {
	assertReasonBucket(t, "IRR_SYNTAX", bill.ReasonKeyLegal)
}
func TestReasonIRR_TAILLE_PJ(t *testing.T) {
	assertReasonBucket(t, "IRR_TAILLE_PJ", bill.ReasonKeyLegal)
}
func TestReasonIRR_NOM_PJ(t *testing.T) {
	assertReasonBucket(t, "IRR_NOM_PJ", bill.ReasonKeyLegal)
}
func TestReasonIRR_VID_PJ(t *testing.T) {
	assertReasonBucket(t, "IRR_VID_PJ", bill.ReasonKeyLegal)
}
func TestReasonIRR_EXT_DOC(t *testing.T) {
	assertReasonBucket(t, "IRR_EXT_DOC", bill.ReasonKeyLegal)
}
func TestReasonIRR_TAILLE_F(t *testing.T) {
	assertReasonBucket(t, "IRR_TAILLE_F", bill.ReasonKeyLegal)
}
func TestReasonIRR_ANTIVIRUS(t *testing.T) {
	assertReasonBucket(t, "IRR_ANTIVIRUS", bill.ReasonKeyLegal)
}

func TestReasonUnknownCodeMisses(t *testing.T) {
	_, ok := ReasonKeyFor("NONEXISTENT")
	assert.False(t, ok)
}

// Default-for-key: one per bucket with codes.

func TestReasonDefaultForUnknownReceiver(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyUnknownReceiver)
	assert.True(t, ok)
	assert.Equal(t, "DEST_INC", got)
}

func TestReasonDefaultForReferences(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyReferences)
	assert.True(t, ok)
	assert.Equal(t, "CMD_ERR", got)
}

func TestReasonDefaultForOther(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyOther)
	assert.True(t, ok)
	assert.Equal(t, "AUTRE", got)
}

func TestReasonDefaultForFinanceTerms(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyFinanceTerms)
	assert.True(t, ok)
	assert.Equal(t, "COORD_BANC_ERR", got)
}

func TestReasonDefaultForLegal(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyLegal)
	assert.True(t, ok)
	assert.Equal(t, "NON_CONFORME", got)
}

func TestReasonDefaultForPrices(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyPrices)
	assert.True(t, ok)
	assert.Equal(t, "PU_ERR", got)
}

func TestReasonDefaultForNotRecognized(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyNotRecognized)
	assert.True(t, ok)
	assert.Equal(t, "DOUBLON", got)
}

func TestReasonDefaultForQuantity(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyQuantity)
	assert.True(t, ok)
	assert.Equal(t, "QTE_ERR", got)
}

func TestReasonDefaultForItems(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyItems)
	assert.True(t, ok)
	assert.Equal(t, "ART_ERR", got)
}

func TestReasonDefaultForPaymentTerms(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyPaymentTerms)
	assert.True(t, ok)
	assert.Equal(t, "MODPAI_ERR", got)
}

func TestReasonDefaultForQuality(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyQuality)
	assert.True(t, ok)
	assert.Equal(t, "QUALITE_ERR", got)
}

func TestReasonDefaultForDelivery(t *testing.T) {
	got, ok := CDARReasonCodeFor(bill.ReasonKeyDelivery)
	assert.True(t, ok)
	assert.Equal(t, "LIVR_INCOMP", got)
}

func TestReasonDefaultForKeyUnknownMisses(t *testing.T) {
	_, ok := CDARReasonCodeFor("made-up-key")
	assert.False(t, ok)
}

// --- Internal helper coverage -------------------------------------------

func TestStatusTypeForKeyUnknown(t *testing.T) {
	_, ok := statusTypeForKey("unknown")
	assert.False(t, ok)
}
