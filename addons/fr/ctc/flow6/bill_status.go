package flow6

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeStatus(st *bill.Status) {
	if st == nil {
		return
	}

	// Ensure the supplier and customer parties carry the expected CDV role when missing.
	if st.Supplier != nil {
		st.Supplier.Ext = st.Supplier.Ext.SetIfEmpty(ExtKeyRole, RoleSeller)
	}
	if st.Customer != nil {
		st.Customer.Ext = st.Customer.Ext.SetIfEmpty(ExtKeyRole, RoleBuyer)
	}

	for _, line := range st.Lines {
		// There should only be one line
		normalizeStatusLine(st, line)
	}
}

func normalizeStatusLine(s *bill.Status, line *bill.StatusLine) {
	if line == nil {
		return
	}

	// If the line already has an extension, this will set the key.
	prepareStatusWithLine(s, line)

	switch s.Type {
	case bill.StatusTypeUpdate, bill.StatusTypeSystem:
		// Issue and System status types
		switch line.Key {
		case bill.StatusLineIssued:
			line.Ext = line.Ext.Set(ExtKeyStatus, "200")
		default:
			line.Ext = line.Ext.Delete(ExtKeyStatus)
		}
	case bill.StatusTypeResponse:
		// Response types from the invoice recipient.
		switch line.Key {
		case bill.StatusLineIssued:
			line.Ext = line.Ext.Set(ExtKeyStatus, "201")
		case bill.StatusLineAcknowledged:
			line.Ext = line.Ext.Set(ExtKeyStatus, "202")
		case bill.StatusLineProcessing:
			line.Ext = line.Ext.Set(ExtKeyStatus, "204")
		case bill.StatusLineAccepted:
			line.Ext = line.Ext.Set(ExtKeyStatus, "205")
		case bill.StatusLineQuerying:
			line.Ext = line.Ext.Set(ExtKeyStatus, "208")
		case bill.StatusLineRejected:
			line.Ext = line.Ext.SetOneOf(ExtKeyStatus,
				"210", // rejected (default)
				"207", // disputed
				"206", // partially accepted
			)
		default:
			line.Ext = line.Ext.Delete(ExtKeyStatus)
		}
	}
}

func prepareStatusWithLine(s *bill.Status, line *bill.StatusLine) {
	sk := line.Ext.Get(ExtKeyStatus)
	if !sk.IsEmpty() {
		switch sk {
		case "200":
			s.Type = bill.StatusTypeUpdate
			line.Key = bill.StatusLineIssued
		case "201":
			s.Type = bill.StatusTypeResponse
			line.Key = bill.StatusLineIssued
		case "202":
			s.Type = bill.StatusTypeResponse
			line.Key = bill.StatusLineAcknowledged
		case "204":
			s.Type = bill.StatusTypeResponse
			line.Key = bill.StatusLineProcessing
		case "205":
			s.Type = bill.StatusTypeResponse
			line.Key = bill.StatusLineAccepted
		case "206", "207", "210":
			s.Type = bill.StatusTypeResponse
			line.Key = bill.StatusLineRejected
		case "208":
			s.Type = bill.StatusTypeResponse
			line.Key = bill.StatusLineQuerying
		case "213":
			s.Type = bill.StatusTypeResponse
			line.Key = bill.StatusLineError
		}
		return
	}
	if !s.Type.IsEmpty() {
		return
	}
	// No ext, no Type: derive Type from the line.Key. Issued is the
	// only line key that's ambiguous (update→200 vs response→201);
	// default to update there.
	switch line.Key {
	case bill.StatusLineIssued:
		s.Type = bill.StatusTypeUpdate
	case bill.StatusLineAcknowledged, bill.StatusLineProcessing,
		bill.StatusLineAccepted, bill.StatusLineQuerying,
		bill.StatusLineRejected, bill.StatusLineError:
		s.Type = bill.StatusTypeResponse
	}
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

// billStatusRules validates the integrity of the addon's own extensions
// and the supported document shape. French CTC format/business rules
// (BR-FR-CDV-*) are the converter's responsibility — see the package doc.
func billStatusRules() *rules.Set {
	return rules.For(new(bill.Status),
		rules.Field("type",
			rules.Assert("01", "status type must be one of: response, update",
				is.In(bill.StatusTypeResponse, bill.StatusTypeUpdate),
			),
		),
		rules.Field("ext",
			rules.Assert("02", "status ext fr-ctc-flow6-status must be a Status-applicable ProcessConditionCode (200-210 or 213); codes 211, 212 belong on bill.Payment",
				tax.ExtensionsHasCodes(ExtKeyStatus, statusProcessCodes...),
			),
		),
		rules.Field("lines",
			rules.Each(
				rules.Field("key",
					rules.Assert("03", "status line key must be a recognised Flow 6 event",
						is.In(
							bill.StatusLineIssued, bill.StatusLineAcknowledged,
							bill.StatusLineProcessing, bill.StatusLineAccepted,
							bill.StatusLineQuerying, bill.StatusLineRejected,
							bill.StatusLineError,
						),
					),
				),
			),
		),
	)
}

// -- bill.Reason --------------------------------------------------------

// normalizeReason maps between bill.Reason.Key (Peppol-aligned
// rejection bucket) and the CDAR extensions on the Reason. Following
// addons/es/verifactu/tax.go's normalizeTaxCombo pattern: a reverse
// step via prepareReasonKey recovers Key from a previously-set
// ReasonCode extension, then a forward switch chains SetOneOf calls
// so the CDAR ReasonCode (fr-ctc-flow6-reason) and CharacteristicType
// (fr-ctc-flow6-condition) defaults are populated when missing,
// while an explicit caller pick within the bucket is preserved.
//
// The Reason carries one fr-ctc-flow6-condition value (CDAR cardinality
// 0..1 per SpecifiedDocumentStatus). For multiple kinds of
// characteristic on the same status line — e.g. DIV + DVA — add
// multiple Reasons. The bill.Condition entries on each Reason are
// reserved for Peppol cac:Condition-style business-rule codes
// describing the affected field and value.
func normalizeReason(r *bill.Reason) {
	if r == nil {
		return
	}

	// Reverse step: fill Reason.Key from the CDAR ReasonCode ext
	// when only the ext is set (round-tripping a parsed CDV).
	prepareReasonKey(r)

	// Forward step: per bucket, SetOneOf defaults each ext to the
	// first listed CDAR code and preserves any caller-set value that
	// already matches one of the bucket's other allowed codes.
	switch r.Key {
	case bill.ReasonKeyFinanceTerms:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason, "COORD_BANC_ERR").
			SetOneOf(ExtKeyCondition,
				ConditionBankDetailsUpdate, ConditionInvalidData,
				ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyOther:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason, "AUTRE").
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyLegal:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason,
				"NON_CONFORME", "TX_TVA_ERR",
				"REJ_SEMAN", "REJ_COH", "REJ_CONT_B2G",
				"IRR_VIDE_F", "IRR_TYPE_F", "IRR_SYNTAX",
				"IRR_TAILLE_PJ", "IRR_NOM_PJ", "IRR_VID_PJ",
				"IRR_EXT_DOC", "IRR_TAILLE_F", "IRR_ANTIVIRUS",
			).
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyNotRecognized:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason,
				"DOUBLON", "TRANSAC_INC", "EMMET_INC",
				"CONTRAT_TERM", "DOUBLE_FACT", "REJ_UNI",
			).
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyUnknownReceiver:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason,
				"DEST_INC", "NON_TRANSMISE", "ROUTAGE_ERR",
			).
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyReferences:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason,
				"CMD_ERR", "JUSTIF_ABS", "DEST_ERR",
				"ADR_ERR", "SIRET_ERR", "CODE_ROUTAGE_ERR",
				"REF_CT_ABSENT", "REF_ERR",
				"REJ_ADR", "REJ_REF_PJ", "REJ_ASS_PJ",
			).
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyPrices:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason,
				"PU_ERR", "MONTANTTOTAL_ERR", "CALCUL_ERR", "REM_ERR",
			).
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyQuantity:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason, "QTE_ERR").
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyItems:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason, "ART_ERR").
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyPaymentTerms:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason, "MODPAI_ERR").
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyQuality:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason, "QUALITE_ERR").
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	case bill.ReasonKeyDelivery:
		r.Ext = r.Ext.
			SetOneOf(ExtKeyReason, "LIVR_INCOMP").
			SetOneOf(ExtKeyCondition,
				ConditionInvalidData, ConditionExpectedData, ConditionReplacementData,
			)
	}
}

// normalizeAction maps between bill.Action.Key (Peppol-aligned) and
// the CDAR RequestedActionCode (MDT-121) extension. Mirrors
// normalizeReason: a reverse step via prepareActionKey recovers Key
// from a previously-set ext, then a forward switch chains SetOneOf
// calls to default the ext when missing while preserving any
// caller-set value.
func normalizeAction(a *bill.Action) {
	if a == nil {
		return
	}
	prepareActionKey(a)
	switch a.Key {
	case bill.ActionKeyNone:
		a.Ext = a.Ext.Set(ExtKeyAction, "NOA")
	case bill.ActionKeyProvide:
		a.Ext = a.Ext.Set(ExtKeyAction, "PIN")
	case bill.ActionKeyReissue:
		a.Ext = a.Ext.Set(ExtKeyAction, "NIN")
	case bill.ActionKeyCreditFull:
		a.Ext = a.Ext.Set(ExtKeyAction, "CNF")
	case bill.ActionKeyCreditPartial:
		a.Ext = a.Ext.Set(ExtKeyAction, "CNP")
	case bill.ActionKeyCreditAmount:
		a.Ext = a.Ext.Set(ExtKeyAction, "CNA")
	case bill.ActionKeyOther:
		a.Ext = a.Ext.Set(ExtKeyAction, "OTH")
	}
}

// prepareActionKey reverse-maps the CDAR RequestedActionCode
// extension to its bill.Action.Key when the caller has only set the
// ext (e.g. when round-tripping a parsed CDV).
func prepareActionKey(a *bill.Action) {
	if !a.Key.IsEmpty() {
		return
	}
	switch a.Ext.Get(ExtKeyAction) {
	case "NOA":
		a.Key = bill.ActionKeyNone
	case "PIN":
		a.Key = bill.ActionKeyProvide
	case "NIN":
		a.Key = bill.ActionKeyReissue
	case "CNF":
		a.Key = bill.ActionKeyCreditFull
	case "CNP":
		a.Key = bill.ActionKeyCreditPartial
	case "CNA":
		a.Key = bill.ActionKeyCreditAmount
	case "OTH":
		a.Key = bill.ActionKeyOther
	}
}

// billActionRules validates the fr-ctc-flow6-action extension carried
// on a bill.Action under a bill.Status. tax.ExtensionHasValidCode
// reads the registered Values list of ExtKeyAction.
func billActionRules() *rules.Set {
	return rules.For(new(bill.Action),
		rules.Field("ext",
			rules.Assert("01", "action ext fr-ctc-flow6-action must be a known CDAR RequestedActionCode (MDT-121)",
				tax.ExtensionHasValidCode(ExtKeyAction),
			),
		),
	)
}

func billReasonRules() *rules.Set {
	return rules.For(new(bill.Reason),
		rules.Field("ext",
			rules.Assert("01", "reason ext fr-ctc-flow6-reason must be a known CDAR ReasonCode",
				tax.ExtensionHasValidCode(ExtKeyReason),
			),
			rules.Assert("02", "reason ext fr-ctc-flow6-condition must be a Status-applicable CharacteristicTypeCode (CBB, DIV, DVA, MAJ, MAP, MAPTTC, MNA, MNATTC, ESC, RAB, REM); MEN, MPA, RAP belong on bill.Payment",
				tax.ExtensionsHasCodes(ExtKeyCondition, statusConditionCodes...),
			),
		),
	)
}

// prepareReasonKey reverse-maps the CDAR ReasonCode extension to its
// bill.Reason.Key bucket when the caller has only set the ext (e.g.
// when round-tripping a parsed CDV).
func prepareReasonKey(r *bill.Reason) {
	if !r.Key.IsEmpty() {
		return
	}
	switch r.Ext.Get(ExtKeyReason) {
	case "COORD_BANC_ERR":
		r.Key = bill.ReasonKeyFinanceTerms
	case "AUTRE":
		r.Key = bill.ReasonKeyOther
	case "NON_CONFORME", "TX_TVA_ERR",
		"REJ_SEMAN", "REJ_COH", "REJ_CONT_B2G",
		"IRR_VIDE_F", "IRR_TYPE_F", "IRR_SYNTAX",
		"IRR_TAILLE_PJ", "IRR_NOM_PJ", "IRR_VID_PJ",
		"IRR_EXT_DOC", "IRR_TAILLE_F", "IRR_ANTIVIRUS":
		r.Key = bill.ReasonKeyLegal
	case "DOUBLON", "TRANSAC_INC", "EMMET_INC",
		"CONTRAT_TERM", "DOUBLE_FACT", "REJ_UNI":
		r.Key = bill.ReasonKeyNotRecognized
	case "DEST_INC", "NON_TRANSMISE", "ROUTAGE_ERR":
		r.Key = bill.ReasonKeyUnknownReceiver
	case "CMD_ERR", "JUSTIF_ABS", "DEST_ERR",
		"ADR_ERR", "SIRET_ERR", "CODE_ROUTAGE_ERR",
		"REF_CT_ABSENT", "REF_ERR",
		"REJ_ADR", "REJ_REF_PJ", "REJ_ASS_PJ":
		r.Key = bill.ReasonKeyReferences
	case "PU_ERR", "MONTANTTOTAL_ERR", "CALCUL_ERR", "REM_ERR":
		r.Key = bill.ReasonKeyPrices
	case "QTE_ERR":
		r.Key = bill.ReasonKeyQuantity
	case "ART_ERR":
		r.Key = bill.ReasonKeyItems
	case "MODPAI_ERR":
		r.Key = bill.ReasonKeyPaymentTerms
	case "QUALITE_ERR":
		r.Key = bill.ReasonKeyQuality
	case "LIVR_INCOMP":
		r.Key = bill.ReasonKeyDelivery
	}
}
