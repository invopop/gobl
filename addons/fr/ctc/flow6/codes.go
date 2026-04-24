package flow6

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

// Extended bill.StatusLine.Key values added by Flow 6 so every CDAR
// ProcessConditionCode maps 1:1 to a (key, type) pair. GOBL ships the
// "plain" keys (issued, processing, accepted, rejected, paid, error);
// the ones marked here are France-specific additions needed for CDAR.
const (
	StatusEventIssuedByPlatform   cbc.Key = "issued-by-platform"
	StatusEventReceivedByPlatform cbc.Key = "received-by-platform"
	StatusEventMadeAvailable      cbc.Key = "made-available"
	StatusEventPartiallyAccepted  cbc.Key = "partially-accepted"
	StatusEventDisputed           cbc.Key = "disputed"
	StatusEventSuspended          cbc.Key = "suspended"
	StatusEventCompleted          cbc.Key = "completed"
	StatusEventPaymentForwarded   cbc.Key = "payment-forwarded"
)

// processEntry pairs a bill.StatusLine.Key with the bill.Status.Type it
// implies. For Flow 6 the pair is always fixed per key: the Type is a
// property of the CDAR code, not a disambiguator.
type processEntry struct {
	Key  cbc.Key
	Type cbc.Key
	Code string
}

// processTable is the authoritative ProcessConditionCode mapping for
// Flow 6 CDAR messages. Order is stable and matches the spec table.
var processTable = []processEntry{
	{bill.StatusEventIssued, bill.StatusTypeUpdate, "200"},
	{StatusEventIssuedByPlatform, bill.StatusTypeUpdate, "201"},
	{StatusEventReceivedByPlatform, bill.StatusTypeResponse, "202"},
	{StatusEventMadeAvailable, bill.StatusTypeResponse, "203"},
	{bill.StatusEventProcessing, bill.StatusTypeResponse, "204"},
	{bill.StatusEventAccepted, bill.StatusTypeResponse, "205"},
	{StatusEventPartiallyAccepted, bill.StatusTypeResponse, "206"},
	{StatusEventDisputed, bill.StatusTypeResponse, "207"},
	{StatusEventSuspended, bill.StatusTypeResponse, "208"},
	{StatusEventCompleted, bill.StatusTypeResponse, "209"},
	{bill.StatusEventRejected, bill.StatusTypeResponse, "210"},
	{StatusEventPaymentForwarded, bill.StatusTypeUpdate, "211"},
	{bill.StatusEventPaid, bill.StatusTypeResponse, "212"},
	{bill.StatusEventError, bill.StatusTypeResponse, "213"},
}

// CDARProcessCodeFor returns the CDAR ProcessConditionCode for a bill
// StatusLine.Key + Status.Type pair. Returns ("", false) if the pair is
// unknown or the Type does not match the fixed Type for the key.
func CDARProcessCodeFor(key cbc.Key, typ cbc.Key) (string, bool) {
	for _, e := range processTable {
		if e.Key == key && e.Type == typ {
			return e.Code, true
		}
	}
	return "", false
}

// StatusKeyFor returns the (StatusLine.Key, Status.Type) pair for a CDAR
// ProcessConditionCode. Returns ("", "", false) if the code is unknown.
func StatusKeyFor(code string) (cbc.Key, cbc.Key, bool) {
	for _, e := range processTable {
		if e.Code == code {
			return e.Key, e.Type, true
		}
	}
	return "", "", false
}

// statusTypeForKey returns the fixed Status.Type associated with a
// StatusLine.Key for Flow 6. The second return is false if the key has
// no CDAR entry.
func statusTypeForKey(key cbc.Key) (cbc.Key, bool) {
	for _, e := range processTable {
		if e.Key == key {
			return e.Type, true
		}
	}
	return "", false
}

// reasonEntry pairs a CDAR ReasonCode with its bucket bill.Reason.Key
// and flags whether this code is the default emitted when the caller
// has not pinned an exact ReasonCode via the extension.
type reasonEntry struct {
	Code      string
	Key       cbc.Key
	IsDefault bool
}

// reasonTable lists all 45 French CDAR reason codes and the bill.Reason
// bucket they roll up to. IsDefault marks the code the generator should
// emit when the caller only sets Reason.Key (see CDARReasonCodeFor).
var reasonTable = []reasonEntry{
	// Business rejection reasons (codes carried on 206 / 207 / 208 / 210).
	{"NON_TRANSMISE", bill.ReasonKeyUnknownReceiver, false},
	{"JUSTIF_ABS", bill.ReasonKeyReferences, false},
	{"ROUTAGE_ERR", bill.ReasonKeyUnknownReceiver, false},
	{"AUTRE", bill.ReasonKeyOther, true},
	{"COORD_BANC_ERR", bill.ReasonKeyFinanceTerms, true},
	{"TX_TVA_ERR", bill.ReasonKeyLegal, false},
	{"MONTANTTOTAL_ERR", bill.ReasonKeyPrices, false},
	{"CALCUL_ERR", bill.ReasonKeyPrices, false},
	{"NON_CONFORME", bill.ReasonKeyLegal, true},
	{"DOUBLON", bill.ReasonKeyNotRecognized, true},
	{"DEST_INC", bill.ReasonKeyUnknownReceiver, true},
	{"DEST_ERR", bill.ReasonKeyReferences, false},
	{"TRANSAC_INC", bill.ReasonKeyNotRecognized, false},
	{"EMMET_INC", bill.ReasonKeyNotRecognized, false},
	{"CONTRAT_TERM", bill.ReasonKeyNotRecognized, false},
	{"DOUBLE_FACT", bill.ReasonKeyNotRecognized, false},
	{"CMD_ERR", bill.ReasonKeyReferences, true},
	{"ADR_ERR", bill.ReasonKeyReferences, false},
	{"SIRET_ERR", bill.ReasonKeyReferences, false},
	{"CODE_ROUTAGE_ERR", bill.ReasonKeyReferences, false},
	{"REF_CT_ABSENT", bill.ReasonKeyReferences, false},
	{"REF_ERR", bill.ReasonKeyReferences, false},
	{"PU_ERR", bill.ReasonKeyPrices, true},
	{"REM_ERR", bill.ReasonKeyPrices, false},
	{"QTE_ERR", bill.ReasonKeyQuantity, true},
	{"ART_ERR", bill.ReasonKeyItems, true},
	{"MODPAI_ERR", bill.ReasonKeyPaymentTerms, true},
	{"QUALITE_ERR", bill.ReasonKeyQuality, true},
	{"LIVR_INCOMP", bill.ReasonKeyDelivery, true},

	// Technical / platform rejection reasons (code 213 only).
	{"REJ_SEMAN", bill.ReasonKeyLegal, false},
	{"REJ_UNI", bill.ReasonKeyNotRecognized, false},
	{"REJ_COH", bill.ReasonKeyLegal, false},
	{"REJ_ADR", bill.ReasonKeyReferences, false},
	{"REJ_CONT_B2G", bill.ReasonKeyLegal, false},
	{"REJ_REF_PJ", bill.ReasonKeyReferences, false},
	{"REJ_ASS_PJ", bill.ReasonKeyReferences, false},
	{"IRR_VIDE_F", bill.ReasonKeyLegal, false},
	{"IRR_TYPE_F", bill.ReasonKeyLegal, false},
	{"IRR_SYNTAX", bill.ReasonKeyLegal, false},
	{"IRR_TAILLE_PJ", bill.ReasonKeyLegal, false},
	{"IRR_NOM_PJ", bill.ReasonKeyLegal, false},
	{"IRR_VID_PJ", bill.ReasonKeyLegal, false},
	{"IRR_EXT_DOC", bill.ReasonKeyLegal, false},
	{"IRR_TAILLE_F", bill.ReasonKeyLegal, false},
	{"IRR_ANTIVIRUS", bill.ReasonKeyLegal, false},
}

// CDARReasonCodeFor returns the default CDAR ReasonCode for a
// bill.Reason.Key. Used on generate when the caller did not pin an
// exact code via Reason.Ext["fr-ctc-reason-code"].
func CDARReasonCodeFor(key cbc.Key) (string, bool) {
	for _, e := range reasonTable {
		if e.Key == key && e.IsDefault {
			return e.Code, true
		}
	}
	return "", false
}

// ReasonKeyFor returns the bucket bill.Reason.Key for a CDAR
// ReasonCode. Used on parse and by the normalizer to fill Reason.Key
// from the extension.
func ReasonKeyFor(code string) (cbc.Key, bool) {
	for _, e := range reasonTable {
		if e.Code == code {
			return e.Key, true
		}
	}
	return "", false
}

// actionTable maps CDAR RequestedActionCode 1:1 to bill.Action.Key.
var actionTable = []struct {
	Code string
	Key  cbc.Key
}{
	{"NOA", bill.ActionKeyNone},
	{"PIN", bill.ActionKeyProvide},
	{"NIN", bill.ActionKeyReissue},
	{"CNF", bill.ActionKeyCreditFull},
	{"CNP", bill.ActionKeyCreditPartial},
	{"CNA", bill.ActionKeyCreditAmount},
	{"OTH", bill.ActionKeyOther},
}

// CDARActionCodeFor returns the CDAR RequestedActionCode for a
// bill.Action.Key.
func CDARActionCodeFor(key cbc.Key) (string, bool) {
	for _, e := range actionTable {
		if e.Key == key {
			return e.Code, true
		}
	}
	return "", false
}

// ActionKeyFor returns the bill.Action.Key for a CDAR
// RequestedActionCode.
func ActionKeyFor(code string) (cbc.Key, bool) {
	for _, e := range actionTable {
		if e.Code == code {
			return e.Key, true
		}
	}
	return "", false
}
