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

	// ExtKeyReasonCode pins the exact CDAR ReasonCode for a bill.Reason
	// on a Flow 6 message.
	ExtKeyReasonCode cbc.Key = "fr-ctc-flow6-reason-code"

	// ExtKeyStatusCode surfaces the CDAR ProcessConditionCode (MDT-9)
	// on a Flow 6 bill.Status.
	ExtKeyStatusCode cbc.Key = "fr-ctc-flow6-status-code"
)

// Flow 6 party role codes — UNCL 3035 subset repurposed by CDAR
// (MDT-158).
const (
	RoleBY  cbc.Code = "BY"  // Acheteur (Buyer)
	RoleDL  cbc.Code = "DL"  // Affactureur (Factor)
	RoleSE  cbc.Code = "SE"  // Vendeur (Seller)
	RoleAB  cbc.Code = "AB"  // Agent d'acheteur (Buyer's agent)
	RoleSR  cbc.Code = "SR"  // Agent de vendeur (Seller's agent)
	RoleWK  cbc.Code = "WK"  // Plateforme / opérateur de dématérialisation
	RoleDFH cbc.Code = "DFH" // Portail public de facturation (PPF)
	RolePE  cbc.Code = "PE"  // Bénéficiaire (Payee)
	RolePR  cbc.Code = "PR"  // Payeur (Payer)
	RoleII  cbc.Code = "II"  // Invoicer (issuer of invoice)
	RoleIV  cbc.Code = "IV"  // Invoicee
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
				for each populated party on a Flow 6 lifecycle message.
				Labels follow the French CTC specification, which assigns
				CDAR-specific meanings to WK (dematerialisation platform /
				operator) and DFH (Portail Public de Facturation).
				The normalizer fills the obvious defaults (Supplier → SE,
				Customer → BY) and leaves the rest for the caller to set
				explicitly.
			`),
		},
		Values: []*cbc.Definition{
			{Code: RoleBY, Name: i18n.String{i18n.EN: "Buyer", i18n.FR: "Acheteur"}},
			{Code: RoleDL, Name: i18n.String{i18n.EN: "Factor", i18n.FR: "Affactureur"}},
			{Code: RoleSE, Name: i18n.String{i18n.EN: "Seller", i18n.FR: "Vendeur"}},
			{Code: RoleAB, Name: i18n.String{i18n.EN: "Buyer's agent", i18n.FR: "Agent d'acheteur"}},
			{Code: RoleSR, Name: i18n.String{i18n.EN: "Seller's agent", i18n.FR: "Agent de vendeur"}},
			{Code: RoleWK, Name: i18n.String{i18n.EN: "Dematerialisation platform or operator", i18n.FR: "Plateforme ou opérateur de dématérialisation"}},
			{Code: RoleDFH, Name: i18n.String{i18n.EN: "Portail Public de Facturation (PPF)", i18n.FR: "Portail Public de Facturation"}},
			{Code: RolePE, Name: i18n.String{i18n.EN: "Payee", i18n.FR: "Bénéficiaire"}},
			{Code: RolePR, Name: i18n.String{i18n.EN: "Payer", i18n.FR: "Payeur"}},
			{Code: RoleII, Name: i18n.String{i18n.EN: "Invoicer", i18n.FR: "Émetteur de la facture"}},
			{Code: RoleIV, Name: i18n.String{i18n.EN: "Invoicee", i18n.FR: "Destinataire de la facture"}},
		},
	},
	{
		Key: ExtKeyReasonCode,
		Name: i18n.String{
			i18n.EN: "CDAR Reason Code",
			i18n.FR: "Code motif CDAR",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Exact CDAR ReasonCode pinned on a bill.Reason for Flow 6
				lifecycle messages. The CDAR ReasonCode dimension is 1:N
				with bill.Reason.Key: this extension lets the caller pick
				the precise code within a bucket. When absent, the
				converter falls back to the default_for_key code for
				Reason.Key.
			`),
		},
		Values: reasonCodeDefinitions(),
	},
	{
		Key: ExtKeyStatusCode,
		Name: i18n.String{
			i18n.EN: "CDAR Process Condition Code",
			i18n.FR: "Code condition processus CDAR",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				CDAR ProcessConditionCode (MDT-9) identifying the lifecycle
				event reported by the Flow 6 message. The normalizer derives
				it from the (StatusLine.Key, Status.Type) pair; callers can
				pre-set it to pin a specific code (e.g. when round-tripping
				a parsed CDV).
			`),
		},
		Values: statusCodeDefinitions(),
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

// reasonCodeDefinitions builds the value list for the
// fr-ctc-flow6-reason-code extension from the authoritative
// reasonTable — avoids drift between the helper table and the
// extension's accepted value set.
func reasonCodeDefinitions() []*cbc.Definition {
	out := make([]*cbc.Definition, len(reasonTable))
	for i, e := range reasonTable {
		out[i] = &cbc.Definition{
			Code: cbc.Code(e.Code),
			Name: i18n.String{i18n.EN: string(e.Key)},
		}
	}
	return out
}

// processCodeLabels carries the official CDAR libellé for each
// ProcessConditionCode. Includes 211 / 212 (payment-related) which
// are emitted from bill.Payment documents but still need to appear
// in the fr-ctc-flow6-status-code extension's allowed-value list.
var processCodeLabels = map[string]string{
	"200": "Déposée",
	"201": "Émise par la plateforme",
	"202": "Reçue par PA",
	"203": "Mise à disposition",
	"204": "Prise en charge",
	"205": "Approuvée",
	"206": "Approuvée partiellement",
	"207": "En litige",
	"208": "Suspendue",
	"209": "Complétée",
	"210": "Refusée",
	"211": "Paiement transmis",
	"212": "Encaissée",
	"213": "Rejetée sémantique",
}

// paymentProcessCodes lists the payment-related CDAR codes that are
// allowed values for fr-ctc-flow6-status-code but do not appear in the
// bill.Status processTable.
var paymentProcessCodes = []string{"211", "212"}

// statusCodeDefinitions builds the value list for
// fr-ctc-flow6-status-code from the authoritative processTable plus
// the payment-only codes (211, 212) carried on bill.Payment.
func statusCodeDefinitions() []*cbc.Definition {
	out := make([]*cbc.Definition, 0, len(processTable)+len(paymentProcessCodes))
	for _, e := range processTable {
		out = append(out, &cbc.Definition{
			Code: cbc.Code(e.Code),
			Name: i18n.String{i18n.EN: processCodeLabels[e.Code]},
		})
	}
	for _, code := range paymentProcessCodes {
		out = append(out, &cbc.Definition{
			Code: cbc.Code(code),
			Name: i18n.String{i18n.EN: processCodeLabels[code]},
		})
	}
	return out
}
