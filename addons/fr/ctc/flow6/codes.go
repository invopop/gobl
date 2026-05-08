package flow6

import (
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

// Extended bill.StatusLine.Key values added by Flow 6. We reuse stock
// GOBL keys wherever the (key, Status.Type) pair is unambiguous —
// notably `paid + update` (CDV-211 Paiement Transmis) and
// `paid + response` (CDV-212 Encaissée), which share the "paid"
// semantic but distinguish transmission vs treatment phase via Type.
//
// Stock keys carrying CDAR codes:
//   issued       + update   → 200 (Déposée)
//   acknowledged + response → 202 (Reçue par PA)
//   processing   + response → 204 (Prise en charge)
//   accepted     + response → 205 (Approuvée)
//   querying     + response → 208 (Suspendue) — "buyer will not proceed
//                                                 without additional info"
//   rejected     + response → 210 (Refusée)
//   paid         + update   → 211 (Paiement transmis)
//   paid         + response → 212 (Encaissée)
//   error        + response → 213 (Rejetée sémantique)
//
// The keys below cover Flow 6 events that have no stock equivalent.
const (
	StatusEventIssuedByPlatform  cbc.Key = "issued-by-platform"
	StatusEventMadeAvailable     cbc.Key = "made-available"
	StatusEventPartiallyAccepted cbc.Key = "partially-accepted"
	StatusEventDisputed          cbc.Key = "disputed"
	StatusEventCompleted         cbc.Key = "completed"
)

// processEntry pairs a bill.StatusLine.Key with the bill.Status.Type
// the CDV expects, alongside the wire ProcessConditionCode. A key may
// appear in multiple entries when its (key, type) pairs map to
// different CDAR codes — that is what allows us to share `paid`
// across CDV-211 (update) and CDV-212 (response).
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
	{bill.StatusEventAcknowledged, bill.StatusTypeResponse, "202"},
	{StatusEventMadeAvailable, bill.StatusTypeResponse, "203"},
	{bill.StatusEventProcessing, bill.StatusTypeResponse, "204"},
	{bill.StatusEventAccepted, bill.StatusTypeResponse, "205"},
	{StatusEventPartiallyAccepted, bill.StatusTypeResponse, "206"},
	{StatusEventDisputed, bill.StatusTypeResponse, "207"},
	{bill.StatusEventQuerying, bill.StatusTypeResponse, "208"},
	{StatusEventCompleted, bill.StatusTypeResponse, "209"},
	{bill.StatusEventRejected, bill.StatusTypeResponse, "210"},
	{bill.StatusEventPaid, bill.StatusTypeUpdate, "211"},
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

// statusTypeForKey returns the Status.Type associated with a
// StatusLine.Key for Flow 6 *if the key has exactly one*. Returns
// ("", false) when the key is unknown OR when the same key is shared
// across multiple types (e.g. `paid` covers both update/211 and
// response/212) — in that case the caller must specify Type explicitly.
func statusTypeForKey(key cbc.Key) (cbc.Key, bool) {
	var found cbc.Key
	for _, e := range processTable {
		if e.Key != key {
			continue
		}
		if found != "" && found != e.Type {
			return "", false
		}
		found = e.Type
	}
	if found == "" {
		return "", false
	}
	return found, true
}

// statusKeyKnown reports whether the key appears in the Flow 6
// process table at least once (regardless of how many types it
// pairs with).
func statusKeyKnown(key cbc.Key) bool {
	for _, e := range processTable {
		if e.Key == key {
			return true
		}
	}
	return false
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

// CDVSide reports which end-party plays the Issuer role on a CDV
// message of the given process code. Used by the cii writer to
// auto-fill IssuerTradeParty / RecipientTradeParty from Supplier and
// Customer when the caller has not set them explicitly. The mapping
// follows Annexe A "Acteurs CDV".
type CDVSide string

const (
	// CDVSideBuyer — the buyer-side end-party issues the CDV
	// (Issuer = Customer, Recipient = Supplier).
	CDVSideBuyer CDVSide = "buyer"
	// CDVSideSeller — the seller-side end-party issues the CDV
	// (Issuer = Supplier, Recipient = Customer).
	CDVSideSeller CDVSide = "seller"
	// CDVSidePlatform — the message is issued by a platform (PA-E,
	// PA-R) or addressed to the PPF, so neither end-party plays the
	// issuer role. The caller must supply st.Issuer (and typically
	// st.Recipient) explicitly.
	CDVSidePlatform CDVSide = "platform"
)

// SideForCode returns which end-party issues a CDV with the given
// CDAR ProcessConditionCode (per Annexe A "Acteurs CDV", treatment
// phase). Returns CDVSideUnknown for codes not in the table.
func SideForCode(code string) CDVSide {
	switch code {
	case "204", "205", "206", "207", "208", "210", "211":
		return CDVSideBuyer
	case "209", "212":
		return CDVSideSeller
	case "200", "201", "202", "203", "213":
		return CDVSidePlatform
	}
	return CDVSidePlatform
}

// SideForKeyType is a convenience wrapper around SideForCode that
// looks up the process code for a (StatusLine.Key, Status.Type) pair
// first.
func SideForKeyType(key, typ cbc.Key) CDVSide {
	if code, ok := CDARProcessCodeFor(key, typ); ok {
		return SideForCode(code)
	}
	return CDVSidePlatform
}

// allowedReasonsByProcessCode is the BR-FR-CDV-CL-09 table — for each
// CDAR process code that admits Reasons, the set of CDAR ReasonCodes
// the schematron will accept. Codes not listed here either don't carry
// reasons (200, 201, 202, 203, 204, 205, 209, 211, 212) or carry any
// reason (the table is the strict list per Annexe A "Tableau des motifs
// de STATUTS").
var allowedReasonsByProcessCode = map[string][]string{
	"200": {"NON_TRANSMISE"},
	"206": {
		"AUTRE", "CMD_ERR", "SIRET_ERR", "CODE_ROUTAGE_ERR",
		"REF_CT_ABSENT", "REF_ERR", "PU_ERR", "REM_ERR", "QTE_ERR",
		"ART_ERR", "MODPAI_ERR", "QUALITE_ERR", "LIVR_INCOMP",
	},
	"207": {
		"AUTRE", "COORD_BANC_ERR", "TX_TVA_ERR", "MONTANTTOTAL_ERR",
		"CALCUL_ERR", "NON_CONFORME", "DOUBLON", "DEST_ERR",
		"TRANSAC_INC", "EMMET_INC", "CONTRAT_TERM", "DOUBLE_FACT",
		"CMD_ERR", "ADR_ERR", "SIRET_ERR", "CODE_ROUTAGE_ERR",
		"REF_CT_ABSENT", "REF_ERR", "PU_ERR", "REM_ERR", "QTE_ERR",
		"ART_ERR", "MODPAI_ERR", "QUALITE_ERR", "LIVR_INCOMP",
	},
	"208": {
		"JUSTIF_ABS", "COORD_BANC_ERR", "CMD_ERR", "SIRET_ERR",
		"CODE_ROUTAGE_ERR", "REF_CT_ABSENT", "REF_ERR",
	},
	"210": {
		"TX_TVA_ERR", "MONTANTTOTAL_ERR", "CALCUL_ERR", "NON_CONFORME",
		"DOUBLON", "DEST_ERR", "TRANSAC_INC", "EMMET_INC", "CONTRAT_TERM",
		"DOUBLE_FACT", "CMD_ERR", "ADR_ERR", "REF_CT_ABSENT",
	},
	"213": {
		"MONTANTTOTAL_ERR", "CALCUL_ERR", "DOUBLON", "ADR_ERR",
		"REJ_SEMAN", "REJ_UNI", "REJ_COH", "REJ_ADR", "REJ_CONT_B2G",
		"REJ_REF_PJ", "REJ_ASS_PJ",
	},
}

// ReasonCodeAllowedForProcessCode reports whether the given CDAR
// ReasonCode is permitted on a status line whose ProcessConditionCode is
// processCode. Returns true when the process code does not constrain the
// reason set (i.e. it isn't in the BR-FR-CDV-CL-09 table).
func ReasonCodeAllowedForProcessCode(reasonCode, processCode string) bool {
	allowed, ok := allowedReasonsByProcessCode[processCode]
	if !ok {
		return true
	}
	return slices.Contains(allowed, reasonCode)
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
