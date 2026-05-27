package flow6

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Flow 6 extension keys.
const (
	// ExtKeyRole carries the CDAR RoleCode for a party (UNCL 3035
	// subset) on a Flow 6 bill.Status message.
	ExtKeyRole cbc.Key = "fr-ctc-flow6-role"

	// ExtKeyReason pins the exact CDAR ReasonCode for a bill.Reason
	// on a Flow 6 message.
	ExtKeyReason cbc.Key = "fr-ctc-flow6-reason"

	// ExtKeyStatus surfaces the CDAR ProcessConditionCode (MDT-9)
	// on a Flow 6 bill.Status. Determined from the Status Line's Key.
	ExtKeyStatus cbc.Key = "fr-ctc-flow6-status"

	// ExtKeyAction pins the CDAR RequestedActionCode (MDT-121) on a
	// bill.Action under a bill.Status. The normalizer fills it from
	// bill.Action.Key (using the actionTable round-trip mapping);
	// callers can override or set the ext directly on round-trip
	// from a parsed CDV.
	ExtKeyAction cbc.Key = "fr-ctc-flow6-action"

	// ExtKeyCondition pins the CDAR CharacteristicTypeCode (MDT-207)
	// on a bill.Reason attached to a Flow 6 rejection / dispute /
	// partially-accepted / completed status line. The cardinality
	// mirrors the spec: one TypeCode per Reason (CDAR's
	// SpecifiedDocumentStatus). To express multiple kinds of
	// characteristic on the same status line — e.g. a "DIV" alongside
	// a sibling "DVA" describing the same field — add multiple
	// Reasons to bill.StatusLine.Reasons, each with its own
	// fr-ctc-flow6-condition value. The accompanying bill.Condition
	// entries on the Reason carry the business-rule codes (Peppol
	// cac:Condition style) plus the affected paths / values.
	ExtKeyCondition cbc.Key = "fr-ctc-flow6-condition"
)

// CDAR CharacteristicTypeCode values (MDT-207) — the controlled
// vocabulary for fr-ctc-flow6-condition on a bill.Reason.
const (
	// -- Field-level correction markers --------------------------------

	// ConditionBankDetailsUpdate (CBB) — coordonnées bancaires
	// du bénéficiaire à modifier sur la facture.
	ConditionBankDetailsUpdate cbc.Code = "CBB"
	// ConditionInvalidData (DIV) — a field on the referenced
	// document carries an invalid value.
	ConditionInvalidData cbc.Code = "DIV"
	// ConditionExpectedData (DVA) — pairs with a sibling DIV
	// Condition to carry the value the receiver expected.
	ConditionExpectedData cbc.Code = "DVA"
	// ConditionReplacementData (MAJ) — carries the value the
	// receiver wants the issuer to apply on the next revision.
	ConditionReplacementData cbc.Code = "MAJ"

	// -- Amount markers (partially-accepted / rejected) ---------------

	// ConditionAmountApprovedHT (MAP) — montant HT approuvé.
	ConditionAmountApprovedHT cbc.Code = "MAP"
	// ConditionAmountApprovedTTC (MAPTTC) — montant TTC
	// approuvé.
	ConditionAmountApprovedTTC cbc.Code = "MAPTTC"
	// ConditionAmountRejectedHT (MNA) — montant HT non
	// approuvé.
	ConditionAmountRejectedHT cbc.Code = "MNA"
	// ConditionAmountRejectedTTC (MNATTC) — montant TTC non
	// approuvé.
	ConditionAmountRejectedTTC cbc.Code = "MNATTC"

	// -- Discount / rebate markers ------------------------------------

	// ConditionDiscount (ESC) — escompte accordé.
	ConditionDiscount cbc.Code = "ESC"
	// ConditionRebate (RAB) — rabais accordé.
	ConditionRebate cbc.Code = "RAB"
	// ConditionReduction (REM) — remise accordée.
	ConditionReduction cbc.Code = "REM"

	// -- Payment-related amounts (round-tripped from bill.Payment) ----

	// ConditionAmountReceived (MEN) — montant encaissé (TTC).
	ConditionAmountReceived cbc.Code = "MEN"
	// ConditionAmountPaid (MPA) — montant payé.
	ConditionAmountPaid cbc.Code = "MPA"
	// ConditionAmountRemaining (RAP) — reste à payer.
	ConditionAmountRemaining cbc.Code = "RAP"
)

// statusProcessCodes lists the ProcessConditionCodes (MDT-9) valid on
// bill.Status.Ext[fr-ctc-flow6-status]. Payment-related codes 211 /
// 212 live on bill.Payment — see paymentProcessCodes.
var statusProcessCodes = []cbc.Code{
	"200", "201", "202", "203", "204", "205",
	"206", "207", "208", "209", "210", "213",
}

// paymentProcessCodes lists the ProcessConditionCodes (MDT-9) valid on
// bill.Payment.Ext[fr-ctc-flow6-status]. The normalizer derives the
// value from bill.Payment.Type: advice → 211, receipt → 212.
var paymentProcessCodes = []cbc.Code{
	"211", "212",
}

// statusConditionCodes lists the CharacteristicTypeCodes (MDT-207)
// valid on bill.Reason.Ext[fr-ctc-flow6-condition] under bill.Status.
// MEN / MPA / RAP live on bill.Payment — see paymentConditionCodes.
var statusConditionCodes = []cbc.Code{
	ConditionBankDetailsUpdate, ConditionInvalidData,
	ConditionExpectedData, ConditionReplacementData,
	ConditionAmountApprovedHT, ConditionAmountApprovedTTC,
	ConditionAmountRejectedHT, ConditionAmountRejectedTTC,
	ConditionDiscount, ConditionRebate, ConditionReduction,
}

// paymentConditionCodes lists the CharacteristicTypeCodes (MDT-207)
// valid on bill.Payment.Ext[fr-ctc-flow6-condition]. The normalizer
// defaults the value from bill.Payment.Type: receipt → MEN, advice →
// MPA. Partial payments can override to RAP.
var paymentConditionCodes = []cbc.Code{
	ConditionAmountReceived,  // MEN
	ConditionAmountPaid,      // MPA
	ConditionAmountRemaining, // RAP
}

// Flow 6 party role codes — UNCL 3035 subset repurposed by CDAR
// (MDT-158).
const (
	RoleBuyer       cbc.Code = "BY"  // Acheteur (Buyer)
	RoleBuyerAgent  cbc.Code = "AB"  // Agent d'acheteur (Buyer's agent)
	RoleFactor      cbc.Code = "DL"  // Affactureur (Factoring company)
	RoleSeller      cbc.Code = "SE"  // Vendeur (Seller)
	RoleSellerAgent cbc.Code = "SR"  // Agent de vendeur (Seller's agent)
	RolePlatform    cbc.Code = "WK"  // Plateforme / opérateur de dématérialisation (Value added network provider)
	RolePPF         cbc.Code = "DFH" // Portail public de facturation (PPF)
	RolePayee       cbc.Code = "PE"  // Bénéficiaire (Payee)
	RolePayer       cbc.Code = "PR"  // Payeur (Payer)
	RoleIssuer      cbc.Code = "II"  // Invoicer (issuer of invoice)
	RoleInvoicee    cbc.Code = "IV"  // Invoicee
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyRole,
		Name: i18n.String{
			i18n.EN: "Party Role Code",
			i18n.FR: "Code rôle partie",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				UNCL 3035 role code carried as the CDAR RoleCode (MDT-158)
				on each populated party of any Flow 6 lifecycle message
				— set on Supplier / Customer of both bill.Status and
				bill.Payment. Labels follow the French CTC
				specification, which assigns CDAR-specific meanings to
				WK (dematerialisation platform / operator) and DFH
				(Portail Public de Facturation). The normalizer fills
				the obvious defaults (Supplier → SE, Customer → BY)
				and leaves the rest for the caller to set explicitly.
			`),
		},
		Values: []*cbc.Definition{
			{Code: RoleBuyer, Name: i18n.String{i18n.EN: "Buyer", i18n.FR: "Acheteur"}},
			{Code: RoleBuyerAgent, Name: i18n.String{i18n.EN: "Buyer's agent", i18n.FR: "Agent d'acheteur"}},
			{Code: RoleFactor, Name: i18n.String{i18n.EN: "Factor", i18n.FR: "Affactureur"}},
			{Code: RoleSeller, Name: i18n.String{i18n.EN: "Seller", i18n.FR: "Vendeur"}},
			{Code: RoleSellerAgent, Name: i18n.String{i18n.EN: "Seller's agent", i18n.FR: "Agent de vendeur"}},
			{Code: RolePlatform, Name: i18n.String{i18n.EN: "Dematerialisation platform or operator", i18n.FR: "Plateforme ou opérateur de dématérialisation"}},
			{Code: RolePPF, Name: i18n.String{i18n.EN: "Portail Public de Facturation (PPF)", i18n.FR: "Portail Public de Facturation"}},
			{Code: RolePayee, Name: i18n.String{i18n.EN: "Payee", i18n.FR: "Bénéficiaire"}},
			{Code: RolePayer, Name: i18n.String{i18n.EN: "Payer", i18n.FR: "Payeur"}},
			{Code: RoleIssuer, Name: i18n.String{i18n.EN: "Invoicer", i18n.FR: "Émetteur de la facture"}},
			{Code: RoleInvoicee, Name: i18n.String{i18n.EN: "Invoicee", i18n.FR: "Destinataire de la facture"}},
		},
	},
	{
		Key: ExtKeyReason,
		Name: i18n.String{
			i18n.EN: "CDAR Reason Code",
			i18n.FR: "Code motif CDAR",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Exact CDAR ReasonCode pinned on a bill.Reason under a
				bill.Status — Status-only (bill.Payment has no
				Reasons). The CDAR ReasonCode dimension is 1:N with
				bill.Reason.Key: this extension lets the caller pick
				the precise code within a bucket. The normalizer
				derives a default from Reason.Key (see
				normalizeReason); callers can override to one of the
				bucket's other allowed codes.
			`),
		},
		Values: []*cbc.Definition{
			// Business rejection reasons (codes carried on 206 / 207 / 208 / 210).
			{Code: "NON_TRANSMISE", Name: i18n.String{i18n.EN: "Not transmitted", i18n.FR: "Non transmise"}},
			{Code: "JUSTIF_ABS", Name: i18n.String{i18n.EN: "Justification absent", i18n.FR: "Justificatif absent"}},
			{Code: "ROUTAGE_ERR", Name: i18n.String{i18n.EN: "Routing error", i18n.FR: "Erreur de routage"}},
			{Code: "AUTRE", Name: i18n.String{i18n.EN: "Other", i18n.FR: "Autre"}},
			{Code: "COORD_BANC_ERR", Name: i18n.String{i18n.EN: "Bank account details error", i18n.FR: "Erreur de coordonnées bancaires"}},
			{Code: "TX_TVA_ERR", Name: i18n.String{i18n.EN: "VAT rate error", i18n.FR: "Erreur de taux de TVA"}},
			{Code: "MONTANTTOTAL_ERR", Name: i18n.String{i18n.EN: "Total amount error", i18n.FR: "Erreur de montant total"}},
			{Code: "CALCUL_ERR", Name: i18n.String{i18n.EN: "Calculation error", i18n.FR: "Erreur de calcul"}},
			{Code: "NON_CONFORME", Name: i18n.String{i18n.EN: "Non-compliant", i18n.FR: "Non conforme"}},
			{Code: "DOUBLON", Name: i18n.String{i18n.EN: "Duplicate", i18n.FR: "Doublon"}},
			{Code: "DEST_INC", Name: i18n.String{i18n.EN: "Unknown recipient", i18n.FR: "Destinataire inconnu"}},
			{Code: "DEST_ERR", Name: i18n.String{i18n.EN: "Recipient error", i18n.FR: "Erreur de destinataire"}},
			{Code: "TRANSAC_INC", Name: i18n.String{i18n.EN: "Unknown transaction", i18n.FR: "Transaction inconnue"}},
			{Code: "EMMET_INC", Name: i18n.String{i18n.EN: "Unknown issuer", i18n.FR: "Émetteur inconnu"}},
			{Code: "CONTRAT_TERM", Name: i18n.String{i18n.EN: "Contract terminated", i18n.FR: "Contrat terminé"}},
			{Code: "DOUBLE_FACT", Name: i18n.String{i18n.EN: "Double invoicing", i18n.FR: "Double facturation"}},
			{Code: "CMD_ERR", Name: i18n.String{i18n.EN: "Order reference error", i18n.FR: "Erreur de référence de commande"}},
			{Code: "ADR_ERR", Name: i18n.String{i18n.EN: "Address error", i18n.FR: "Erreur d'adresse"}},
			{Code: "SIRET_ERR", Name: i18n.String{i18n.EN: "SIRET error", i18n.FR: "Erreur de SIRET"}},
			{Code: "CODE_ROUTAGE_ERR", Name: i18n.String{i18n.EN: "Routing code error", i18n.FR: "Erreur de code de routage"}},
			{Code: "REF_CT_ABSENT", Name: i18n.String{i18n.EN: "Contract reference absent", i18n.FR: "Référence contrat absente"}},
			{Code: "REF_ERR", Name: i18n.String{i18n.EN: "Reference error", i18n.FR: "Erreur de référence"}},
			{Code: "PU_ERR", Name: i18n.String{i18n.EN: "Unit price error", i18n.FR: "Erreur de prix unitaire"}},
			{Code: "REM_ERR", Name: i18n.String{i18n.EN: "Discount error", i18n.FR: "Erreur de remise"}},
			{Code: "QTE_ERR", Name: i18n.String{i18n.EN: "Quantity error", i18n.FR: "Erreur de quantité"}},
			{Code: "ART_ERR", Name: i18n.String{i18n.EN: "Item error", i18n.FR: "Erreur d'article"}},
			{Code: "MODPAI_ERR", Name: i18n.String{i18n.EN: "Payment method error", i18n.FR: "Erreur de mode de paiement"}},
			{Code: "QUALITE_ERR", Name: i18n.String{i18n.EN: "Quality issue", i18n.FR: "Problème de qualité"}},
			{Code: "LIVR_INCOMP", Name: i18n.String{i18n.EN: "Incomplete delivery", i18n.FR: "Livraison incomplète"}},
			// Technical / platform rejection reasons (code 213 only).
			{Code: "REJ_SEMAN", Name: i18n.String{i18n.EN: "Semantic rejection", i18n.FR: "Rejet sémantique"}},
			{Code: "REJ_UNI", Name: i18n.String{i18n.EN: "Uniqueness violation", i18n.FR: "Rejet pour unicité"}},
			{Code: "REJ_COH", Name: i18n.String{i18n.EN: "Coherence rejection", i18n.FR: "Rejet de cohérence"}},
			{Code: "REJ_ADR", Name: i18n.String{i18n.EN: "Address rejection", i18n.FR: "Rejet d'adresse"}},
			{Code: "REJ_CONT_B2G", Name: i18n.String{i18n.EN: "B2G context rejection", i18n.FR: "Rejet contexte B2G"}},
			{Code: "REJ_REF_PJ", Name: i18n.String{i18n.EN: "Attachment reference rejection", i18n.FR: "Rejet référence pièce jointe"}},
			{Code: "REJ_ASS_PJ", Name: i18n.String{i18n.EN: "Attachment association rejection", i18n.FR: "Rejet association pièce jointe"}},
			{Code: "IRR_VIDE_F", Name: i18n.String{i18n.EN: "Empty file", i18n.FR: "Fichier vide"}},
			{Code: "IRR_TYPE_F", Name: i18n.String{i18n.EN: "Invalid file type", i18n.FR: "Type de fichier invalide"}},
			{Code: "IRR_SYNTAX", Name: i18n.String{i18n.EN: "Syntax error", i18n.FR: "Erreur de syntaxe"}},
			{Code: "IRR_TAILLE_PJ", Name: i18n.String{i18n.EN: "Attachment size", i18n.FR: "Taille de pièce jointe"}},
			{Code: "IRR_NOM_PJ", Name: i18n.String{i18n.EN: "Attachment name", i18n.FR: "Nom de pièce jointe"}},
			{Code: "IRR_VID_PJ", Name: i18n.String{i18n.EN: "Empty attachment", i18n.FR: "Pièce jointe vide"}},
			{Code: "IRR_EXT_DOC", Name: i18n.String{i18n.EN: "External document", i18n.FR: "Document externe"}},
			{Code: "IRR_TAILLE_F", Name: i18n.String{i18n.EN: "File size", i18n.FR: "Taille de fichier"}},
			{Code: "IRR_ANTIVIRUS", Name: i18n.String{i18n.EN: "Antivirus", i18n.FR: "Antivirus"}},
		},
	},
	{
		Key: ExtKeyStatus,
		Name: i18n.String{
			i18n.EN: "CDAR Process Condition Code",
			i18n.FR: "Code condition processus CDAR",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				CDAR ProcessConditionCode (MDT-9) identifying the
				lifecycle event reported by the Flow 6 message.
				Unified vocabulary; the rules narrow the allow-list
				per document type:

				  - bill.Status accepts 200 / 201 / 202 / 203 / 204 /
				    205 / 206 / 207 / 208 / 209 / 210 / 213. The
				    normalizer derives the value from the
				    (StatusLine.Key, Status.Type) pair via processTable.
				  - bill.Payment accepts 211 (Paiement transmis) and
				    212 (Encaissée). The normalizer derives the value
				    from bill.Payment.Type (advice → 211, receipt →
				    212).

				Callers can pre-set the ext directly to pin a specific
				code (e.g. when round-tripping a parsed CDV).
			`),
		},
		Values: []*cbc.Definition{
			{Code: "200", Name: i18n.String{i18n.EN: "Deposited", i18n.FR: "Déposée"}},
			{Code: "201", Name: i18n.String{i18n.EN: "Issued by platform", i18n.FR: "Émise par la plateforme"}},
			{Code: "202", Name: i18n.String{i18n.EN: "Received by PA", i18n.FR: "Reçue par PA"}},
			{Code: "203", Name: i18n.String{i18n.EN: "Made available", i18n.FR: "Mise à disposition"}},
			{Code: "204", Name: i18n.String{i18n.EN: "Taken into account", i18n.FR: "Prise en charge"}},
			{Code: "205", Name: i18n.String{i18n.EN: "Approved", i18n.FR: "Approuvée"}},
			{Code: "206", Name: i18n.String{i18n.EN: "Partially approved", i18n.FR: "Approuvée partiellement"}},
			{Code: "207", Name: i18n.String{i18n.EN: "In dispute", i18n.FR: "En litige"}},
			{Code: "208", Name: i18n.String{i18n.EN: "Suspended", i18n.FR: "Suspendue"}},
			{Code: "209", Name: i18n.String{i18n.EN: "Completed", i18n.FR: "Complétée"}},
			{Code: "210", Name: i18n.String{i18n.EN: "Rejected", i18n.FR: "Refusée"}},
			{Code: "211", Name: i18n.String{i18n.EN: "Payment transmitted", i18n.FR: "Paiement transmis"}},
			{Code: "212", Name: i18n.String{i18n.EN: "Cashed in", i18n.FR: "Encaissée"}},
			{Code: "213", Name: i18n.String{i18n.EN: "Semantically rejected", i18n.FR: "Rejetée sémantique"}},
		},
	},
	{
		Key: ExtKeyCondition,
		Name: i18n.String{
			i18n.EN: "Flow 6 Condition Code",
			i18n.FR: "Code condition Flux 6",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				CDAR CharacteristicTypeCode (MDT-207). Unified
				vocabulary; the rules narrow the allow-list per
				document type:

				  - On bill.Reason under bill.Status: CBB, DIV, DVA,
				    MAJ (field-level corrections), MAP / MAPTTC /
				    MNA / MNATTC (partial-approval / rejection
				    amounts), ESC / RAB / REM (discount markers).
				    Each Reason carries 0..1 condition (CDAR
				    cardinality); to convey multiple characteristics
				    for the same status line, add multiple Reasons.
				    bill.Condition entries on each Reason are reserved
				    for Peppol cac:Condition-style business-rule
				    codes describing the affected field and value.
				  - On bill.Payment.Ext: MEN (Encaissé), MPA (Payé),
				    RAP (Reste à payer). The normalizer defaults the
				    value from bill.Payment.Type (receipt → MEN,
				    advice → MPA); partial payments can override to
				    RAP.
			`),
		},
		Values: []*cbc.Definition{
			{Code: ConditionBankDetailsUpdate, Name: i18n.String{i18n.EN: "Bank details to modify", i18n.FR: "Coordonnées bancaires à modifier"}},
			{Code: ConditionInvalidData, Name: i18n.String{i18n.EN: "Invalid data", i18n.FR: "Donnée invalide"}},
			{Code: ConditionExpectedData, Name: i18n.String{i18n.EN: "Expected valid data", i18n.FR: "Donnée valide attendue"}},
			{Code: ConditionReplacementData, Name: i18n.String{i18n.EN: "Replacement value", i18n.FR: "Donnée à prendre en compte"}},
			{Code: ConditionAmountApprovedHT, Name: i18n.String{i18n.EN: "Approved amount (excl. VAT)", i18n.FR: "Montant HT approuvé"}},
			{Code: ConditionAmountApprovedTTC, Name: i18n.String{i18n.EN: "Approved amount (incl. VAT)", i18n.FR: "Montant TTC approuvé"}},
			{Code: ConditionAmountRejectedHT, Name: i18n.String{i18n.EN: "Rejected amount (excl. VAT)", i18n.FR: "Montant HT non approuvé"}},
			{Code: ConditionAmountRejectedTTC, Name: i18n.String{i18n.EN: "Rejected amount (incl. VAT)", i18n.FR: "Montant TTC non approuvé"}},
			{Code: ConditionDiscount, Name: i18n.String{i18n.EN: "Discount granted", i18n.FR: "Escompte accordé"}},
			{Code: ConditionRebate, Name: i18n.String{i18n.EN: "Rebate granted", i18n.FR: "Rabais accordé"}},
			{Code: ConditionReduction, Name: i18n.String{i18n.EN: "Reduction granted", i18n.FR: "Remise accordée"}},
			{Code: ConditionAmountReceived, Name: i18n.String{i18n.EN: "Amount received", i18n.FR: "Montant encaissé"}},
			{Code: ConditionAmountPaid, Name: i18n.String{i18n.EN: "Amount paid", i18n.FR: "Montant payé"}},
			{Code: ConditionAmountRemaining, Name: i18n.String{i18n.EN: "Amount remaining", i18n.FR: "Reste à payer"}},
		},
	},
	{
		Key: ExtKeyAction,
		Name: i18n.String{
			i18n.EN: "Flow 6 Action Code",
			i18n.FR: "Code action Flux 6",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				CDAR RequestedActionCode (MDT-121) carried on a
				bill.Action under a bill.Status line. Tells the issuer
				what to do next with the referenced invoice — reissue
				it, send a credit note, provide more information, etc.
				The normalizer derives the default from bill.Action.Key
				(NOA / PIN / NIN / CNF / CNP / CNA / OTH);
				callers can override directly on the extension when
				round-tripping a parsed CDV.
			`),
		},
		Values: []*cbc.Definition{
			{Code: "NOA", Name: i18n.String{i18n.EN: "No action", i18n.FR: "Aucune action"}},
			{Code: "PIN", Name: i18n.String{i18n.EN: "Provide information", i18n.FR: "Fournir des informations"}},
			{Code: "NIN", Name: i18n.String{i18n.EN: "Reissue invoice", i18n.FR: "Réémettre la facture"}},
			{Code: "CNF", Name: i18n.String{i18n.EN: "Credit note (full)", i18n.FR: "Avoir total"}},
			{Code: "CNP", Name: i18n.String{i18n.EN: "Credit note (partial)", i18n.FR: "Avoir partiel"}},
			{Code: "CNA", Name: i18n.String{i18n.EN: "Credit note (amount)", i18n.FR: "Avoir d'un montant"}},
			{Code: "OTH", Name: i18n.String{i18n.EN: "Other", i18n.FR: "Autre"}},
		},
	},
}

// extValue unwraps a tax.Extensions value whether the rules engine has
// passed it to us by value or by pointer.
func extValue(v any) tax.Extensions {
	switch e := v.(type) {
	case tax.Extensions:
		return e
	case *tax.Extensions:
		if e == nil {
			return tax.Extensions{}
		}
		return *e
	}
	return tax.Extensions{}
}
