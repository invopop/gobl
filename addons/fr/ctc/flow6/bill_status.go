package flow6

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
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

// partyHasInboxWhenRequired enforces BR-FR-CDV-08: a party whose role
// is not WK (dematerialisation platform / operator) and not DFH
// (PPF) must carry a URIID (electronic inbox).
func partyHasInboxWhenRequired(v any) bool {
	p, ok := v.(*org.Party)
	if !ok || p == nil {
		return true
	}
	role := p.Ext.Get(ExtKeyRole)
	if role == RolePlatform || role == RolePPF {
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
		rules.Field("ext",
			rules.Assert("02", "status ext fr-ctc-flow6-status must be a Status-applicable ProcessConditionCode (200-210 or 213); codes 211, 212 belong on bill.Payment",
				tax.ExtensionsHasCodes(ExtKeyStatus, statusProcessCodes...),
			),
		),
		rules.Field("supplier",
			rules.Assert("03", "status supplier is required (BR-FR-CDV-13)",
				is.Present,
			),
			rules.Assert("04", "status supplier must have an inbox when its role is not WK or DFH (BR-FR-CDV-08)",
				is.Func("supplier has inbox unless WK/DFH", partyHasInboxWhenRequired),
			),
			rules.Field("ext",
				rules.Assert("05", fmt.Sprintf("status supplier ext %s is required (BR-FR-CDV-CL-03)", ExtKeyRole),
					tax.ExtensionsRequire(ExtKeyRole),
				),
			),
			rules.Field("identities",
				rules.Assert("06", "status supplier must have an identity with ISO/IEC 6523 scheme 0002 (SIREN)",
					org.IdentitiesExtensionIn(iso.ExtKeySchemeID, identitySchemeIDSIREN),
				),
			),
		),
		rules.Field("customer",
			rules.Assert("07", "status customer is required (BR-FR-CDV-CL-04)",
				is.Present,
			),
			rules.Assert("08", "status customer must have an inbox when its role is not WK or DFH (BR-FR-CDV-08)",
				is.Func("customer has inbox unless WK/DFH", partyHasInboxWhenRequired),
			),
			rules.Field("ext",
				rules.Assert("09", fmt.Sprintf("status customer ext %s is required (BR-FR-CDV-CL-04)", ExtKeyRole),
					tax.ExtensionsRequire(ExtKeyRole),
				),
			),
			rules.Field("identities",
				rules.Assert("10", "status customer must have at least one identity with an iso-scheme-id in the Flow 6 allow-list; STC 0231 is a Flow 2 invoice concept",
					org.IdentitiesExtensionIn(iso.ExtKeySchemeID, allowedFlow6IdentitySchemes...),
				),
			),
		),
		rules.Field("lines",
			rules.Assert("11", "status lines must contain exactly one entry",
				is.Present, is.Length(1, 1),
			),
			rules.Each(
				rules.Field("doc",
					rules.Assert("12", "status line doc is required (BR-FR-CDV-10)",
						is.Present,
					),
					rules.Field("code",
						rules.Assert("13", "status line doc code is required (BR-FR-CDV-10)",
							is.Present,
						),
					),
					rules.Field("issue_date",
						rules.Assert("14", "status line doc issue_date is required (BR-FR-CDV-11)",
							is.Present,
						),
					),
				),
				rules.Field("key",
					rules.Assert("15", "status line key must be a recognised Flow 6 event",
						is.In(
							bill.StatusLineIssued, bill.StatusLineAcknowledged,
							bill.StatusLineProcessing, bill.StatusLineAccepted,
							bill.StatusLineQuerying, bill.StatusLineRejected,
							bill.StatusLineError,
						),
					),
				),
				rules.When(
					bill.StatusLineKeyIn(bill.StatusLineRejected, bill.StatusLineQuerying, bill.StatusLineError),
					rules.Field("reasons",
						rules.Assert("16", "status line reasons require at least one entry when key is rejected, querying or error (BR-FR-CDV-14)",
							is.Present,
						),
					),
				),
				// Each Reason's CDAR ReasonCode must be in the allow-list
				// for the line's ProcessConditionCode (the
				// line.Ext[ExtKeyStatus] value derived by
				// normalizeStatusLine).
				rules.When(
					lineHasStatusCode("200"),
					rules.Field("reasons",
						rules.Each(
							rules.Field("ext",
								rules.Assert("17", "status line reason ext fr-ctc-flow6-reason for status code 200 (Déposée — transmission rejection) must be NON_TRANSMISE (BR-FR-CDV-CL-09)",
									tax.ExtensionsHasCodes(ExtKeyReason, "NON_TRANSMISE"),
								),
							),
						),
					),
				),
				rules.When(
					lineHasStatusCode("206"),
					rules.Field("reasons",
						rules.Each(
							rules.Field("ext",
								rules.Assert("18", "status line reason ext fr-ctc-flow6-reason for status code 206 (Approuvée partiellement) must be one of AUTRE, CMD_ERR, SIRET_ERR, CODE_ROUTAGE_ERR, REF_CT_ABSENT, REF_ERR, PU_ERR, REM_ERR, QTE_ERR, ART_ERR, MODPAI_ERR, QUALITE_ERR, LIVR_INCOMP (BR-FR-CDV-CL-09)",
									tax.ExtensionsHasCodes(ExtKeyReason,
										"AUTRE", "CMD_ERR", "SIRET_ERR", "CODE_ROUTAGE_ERR",
										"REF_CT_ABSENT", "REF_ERR", "PU_ERR", "REM_ERR", "QTE_ERR",
										"ART_ERR", "MODPAI_ERR", "QUALITE_ERR", "LIVR_INCOMP",
									),
								),
							),
						),
					),
				),
				rules.When(
					lineHasStatusCode("207"),
					rules.Field("reasons",
						rules.Each(
							rules.Field("ext",
								rules.Assert("19", "status line reason ext fr-ctc-flow6-reason for status code 207 (En litige) must be one of AUTRE, COORD_BANC_ERR, TX_TVA_ERR, MONTANTTOTAL_ERR, CALCUL_ERR, NON_CONFORME, DOUBLON, DEST_ERR, TRANSAC_INC, EMMET_INC, CONTRAT_TERM, DOUBLE_FACT, CMD_ERR, ADR_ERR, SIRET_ERR, CODE_ROUTAGE_ERR, REF_CT_ABSENT, REF_ERR, PU_ERR, REM_ERR, QTE_ERR, ART_ERR, MODPAI_ERR, QUALITE_ERR, LIVR_INCOMP (BR-FR-CDV-CL-09)",
									tax.ExtensionsHasCodes(ExtKeyReason,
										"AUTRE", "COORD_BANC_ERR", "TX_TVA_ERR", "MONTANTTOTAL_ERR",
										"CALCUL_ERR", "NON_CONFORME", "DOUBLON", "DEST_ERR",
										"TRANSAC_INC", "EMMET_INC", "CONTRAT_TERM", "DOUBLE_FACT",
										"CMD_ERR", "ADR_ERR", "SIRET_ERR", "CODE_ROUTAGE_ERR",
										"REF_CT_ABSENT", "REF_ERR", "PU_ERR", "REM_ERR", "QTE_ERR",
										"ART_ERR", "MODPAI_ERR", "QUALITE_ERR", "LIVR_INCOMP",
									),
								),
							),
						),
					),
				),
				rules.When(
					lineHasStatusCode("208"),
					rules.Field("reasons",
						rules.Each(
							rules.Field("ext",
								rules.Assert("20", "status line reason ext fr-ctc-flow6-reason for status code 208 (Suspendue) must be one of JUSTIF_ABS, COORD_BANC_ERR, CMD_ERR, SIRET_ERR, CODE_ROUTAGE_ERR, REF_CT_ABSENT, REF_ERR (BR-FR-CDV-CL-09)",
									tax.ExtensionsHasCodes(ExtKeyReason,
										"JUSTIF_ABS", "COORD_BANC_ERR", "CMD_ERR", "SIRET_ERR",
										"CODE_ROUTAGE_ERR", "REF_CT_ABSENT", "REF_ERR",
									),
								),
							),
						),
					),
				),
				rules.When(
					lineHasStatusCode("210"),
					rules.Field("reasons",
						rules.Each(
							rules.Field("ext",
								rules.Assert("21", "status line reason ext fr-ctc-flow6-reason for status code 210 (Refusée) must be one of TX_TVA_ERR, MONTANTTOTAL_ERR, CALCUL_ERR, NON_CONFORME, DOUBLON, DEST_ERR, TRANSAC_INC, EMMET_INC, CONTRAT_TERM, DOUBLE_FACT, CMD_ERR, ADR_ERR, REF_CT_ABSENT (BR-FR-CDV-CL-09)",
									tax.ExtensionsHasCodes(ExtKeyReason,
										"TX_TVA_ERR", "MONTANTTOTAL_ERR", "CALCUL_ERR", "NON_CONFORME",
										"DOUBLON", "DEST_ERR", "TRANSAC_INC", "EMMET_INC", "CONTRAT_TERM",
										"DOUBLE_FACT", "CMD_ERR", "ADR_ERR", "REF_CT_ABSENT",
									),
								),
							),
						),
					),
				),
				rules.When(
					lineHasStatusCode("213"),
					rules.Field("reasons",
						rules.Each(
							rules.Field("ext",
								rules.Assert("22", "status line reason ext fr-ctc-flow6-reason for status code 213 (Rejetée sémantique) must be one of MONTANTTOTAL_ERR, CALCUL_ERR, DOUBLON, ADR_ERR, REJ_SEMAN, REJ_UNI, REJ_COH, REJ_ADR, REJ_CONT_B2G, REJ_REF_PJ, REJ_ASS_PJ (BR-FR-CDV-CL-09)",
									tax.ExtensionsHasCodes(ExtKeyReason,
										"MONTANTTOTAL_ERR", "CALCUL_ERR", "DOUBLON", "ADR_ERR",
										"REJ_SEMAN", "REJ_UNI", "REJ_COH", "REJ_ADR", "REJ_CONT_B2G",
										"REJ_REF_PJ", "REJ_ASS_PJ",
									),
								),
							),
						),
					),
				),
			),
		),
		rules.When(
			bill.StatusTypeIn(bill.StatusTypeResponse),
			rules.Field("lines",
				rules.Each(
					rules.Field("key",
						rules.Assert("23", "status line key must be consistent with status type 'response'",
							is.In(
								bill.StatusLineAcknowledged, bill.StatusLineProcessing,
								bill.StatusLineAccepted, bill.StatusLineQuerying,
								bill.StatusLineRejected, bill.StatusLineError,
							),
						),
					),
				),
			),
		),
		rules.When(
			bill.StatusTypeIn(bill.StatusTypeUpdate),
			rules.Field("lines",
				rules.Each(
					rules.Field("key",
						rules.Assert("24", "status line key must be consistent with status type 'update'",
							is.In(
								bill.StatusLineIssued,
							),
						),
					),
				),
			),
		),
	)
}

// lineHasStatusCode gates a rules.When on the line's CDAR
// ProcessConditionCode (line.Ext[ExtKeyStatus] — set by
// normalizeStatusLine from the (Status.Type, line.Key) pair). Used to
// branch BR-FR-CDV-CL-09's per-process-code reason allow-lists.
func lineHasStatusCode(code cbc.Code) rules.Test {
	return is.Func(fmt.Sprintf("line status code %s", code), func(v any) bool {
		line, ok := v.(*bill.StatusLine)
		return ok && line != nil && line.Ext.Get(ExtKeyStatus) == code
	})
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
