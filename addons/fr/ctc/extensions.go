package ctc

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// French CTC extension keys.
const (
	// ExtKeyBillingMode defines the billing framework mode (B1-B8, S1-S8, M1-M8).
	// Applies to both Flow 2 invoices and Flow 10 B2B reporting invoices.
	ExtKeyBillingMode cbc.Key = "fr-ctc-billing-mode"

	// ExtKeyB2CCategory classifies the B2C transaction for PPF reporting
	// (G1.68). Required on Flow 10 B2C invoices.
	ExtKeyB2CCategory cbc.Key = "fr-ctc-b2c-category"

	// ExtKeyRole carries the CDAR RoleCode for a party (UNCL 3035 subset)
	// on a Flow 6 bill.Status message.
	ExtKeyRole cbc.Key = "fr-ctc-role"

	// ExtKeyReasonCode pins the exact CDAR ReasonCode for a bill.Reason
	// on a Flow 6 message.
	ExtKeyReasonCode cbc.Key = "fr-ctc-reason-code"

	// ExtKeyStatusCode surfaces the CDAR ProcessConditionCode (MDT-9)
	// on a Flow 6 bill.Status.
	ExtKeyStatusCode cbc.Key = "fr-ctc-status-code"
)

// B2C transaction category codes (G1.68).
const (
	// B2CCategoryGoods — deliveries of goods subject to VAT.
	B2CCategoryGoods cbc.Code = "TLB1"
	// B2CCategoryServices — services subject to VAT.
	B2CCategoryServices cbc.Code = "TPS1"
	// B2CCategoryNotTaxable — deliveries / services not subject to VAT in
	// France, including intra-EU distance sales under CGI art. 258 A / 259 B.
	B2CCategoryNotTaxable cbc.Code = "TNT1"
	// B2CCategoryMargin — operations under the VAT-on-margin regime
	// (CGI art. 266-1-e, 268, 297 A).
	B2CCategoryMargin cbc.Code = "TMA1"
)

// Billing mode codes (Cadre de Facturation). Prefix denotes invoice
// nature (B/S/M); numeric suffix encodes payment context.
const (
	BillingModeB1 cbc.Code = "B1"
	BillingModeB2 cbc.Code = "B2"
	BillingModeB4 cbc.Code = "B4"
	BillingModeB7 cbc.Code = "B7"
	BillingModeS1 cbc.Code = "S1"
	BillingModeS2 cbc.Code = "S2"
	BillingModeS4 cbc.Code = "S4"
	BillingModeS5 cbc.Code = "S5"
	BillingModeS6 cbc.Code = "S6"
	BillingModeS7 cbc.Code = "S7"
	BillingModeM1 cbc.Code = "M1"
	BillingModeM2 cbc.Code = "M2"
	BillingModeM4 cbc.Code = "M4"
)

// Flow 6 party role codes — UNCL 3035 subset repurposed by CDAR
// (MDT-158). Labels follow the French specification, NOT the generic
// UNCL 3035 names: WK and DFH in particular have CDAR-specific
// meanings (dematerialisation platform / French public portal).
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
		Key: ExtKeyBillingMode,
		Name: i18n.String{
			i18n.EN: "Billing Mode",
			i18n.FR: "Cadre de Facturation",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to describe the billing framework of the invoice. The billing mode
				indicates the nature of goods/services and the payment context.

				Code prefixes indicate the invoice nature:
				- "B": Goods invoice (Biens)
				- "S": Services invoice
				- "M": Mixed/dual invoice (goods and services that are not accessory to each other)

				The numeric suffix indicates the payment type (1=deposit, 2=already paid,
				4=final after down payment, 5=subcontractor, 6=co-contractor, 7=e-reporting).
			`),
			i18n.FR: here.Doc(`
				Code utilisé pour décrire le cadre de facturation de la facture. Le mode de
				facturation indique la nature des biens/services et le contexte de paiement.

				Les préfixes de code indiquent la nature de la facture :
				- "B" : Facture de biens
				- "S" : Facture de services
				- "M" : Facture mixte (biens et services qui ne sont pas accessoires l'un de l'autre)

				Le suffixe numérique indique le type de paiement (1=dépôt, 2=déjà payée,
				4=définitive après acompte, 5=sous-traitant, 6=cotraitant, 7=e-reporting).
			`),
		},
		Values: []*cbc.Definition{
			{Code: BillingModeB1, Name: i18n.String{i18n.EN: "Goods - Deposit invoice", i18n.FR: "Biens - Facture de dépôt"}},
			{Code: BillingModeB2, Name: i18n.String{i18n.EN: "Goods - Already paid invoice", i18n.FR: "Biens - Facture déjà payée"}},
			{Code: BillingModeB4, Name: i18n.String{i18n.EN: "Goods - Final invoice (after down payment)", i18n.FR: "Biens - Facture définitive (après acompte)"}},
			{Code: BillingModeB7, Name: i18n.String{i18n.EN: "Goods - E-reporting (VAT already collected)", i18n.FR: "Biens - E-reporting (TVA déjà collectée)"}},
			{Code: BillingModeS1, Name: i18n.String{i18n.EN: "Services - Deposit invoice", i18n.FR: "Services - Facture de dépôt"}},
			{Code: BillingModeS2, Name: i18n.String{i18n.EN: "Services - Already paid invoice", i18n.FR: "Services - Facture déjà payée"}},
			{Code: BillingModeS4, Name: i18n.String{i18n.EN: "Services - Final invoice (after down payment)", i18n.FR: "Services - Facture définitive (après acompte)"}},
			{Code: BillingModeS5, Name: i18n.String{i18n.EN: "Services - Subcontractor invoice", i18n.FR: "Services - Facture de sous-traitance"}},
			{Code: BillingModeS6, Name: i18n.String{i18n.EN: "Services - Co-contractor invoice", i18n.FR: "Services - Facture de cotraitance"}},
			{Code: BillingModeS7, Name: i18n.String{i18n.EN: "Services - E-reporting (VAT already collected)", i18n.FR: "Services - E-reporting (TVA déjà collectée)"}},
			{Code: BillingModeM1, Name: i18n.String{i18n.EN: "Mixed - Deposit invoice", i18n.FR: "Mixte - Facture de dépôt"}},
			{Code: BillingModeM2, Name: i18n.String{i18n.EN: "Mixed - Already paid invoice", i18n.FR: "Mixte - Facture déjà payée"}},
			{Code: BillingModeM4, Name: i18n.String{i18n.EN: "Mixed - Final invoice (after down payment)", i18n.FR: "Mixte - Facture définitive (après acompte)"}},
		},
	},
	{
		Key: ExtKeyB2CCategory,
		Name: i18n.String{
			i18n.EN: "B2C Transaction Category",
			i18n.FR: "Catégorie de transaction B2C",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Classifies a B2C transaction for French e-reporting to the PPF
				(G1.68). Required on Flow 10 B2C invoices.

				- TLB1: Goods deliveries subject to VAT.
				- TPS1: Services subject to VAT.
				- TNT1: Goods / services not subject to French VAT, including
				  intra-EU distance sales per CGI articles 258 A and 259 B.
				- TMA1: Operations under the VAT-on-margin regime
				  (CGI articles 266-1-e, 268, 297 A).
			`),
			i18n.FR: here.Doc(`
				Catégorie de transaction pour le e-reporting au PPF (G1.68).
				Obligatoire sur les factures B2C en Flux 10.

				- TLB1 : Livraisons de biens soumises à la TVA.
				- TPS1 : Prestations de services soumises à la TVA.
				- TNT1 : Livraisons et prestations non soumises à la TVA en
				  France, dont les ventes à distance intracommunautaires
				  (CGI art. 258 A et 259 B).
				- TMA1 : Opérations relevant du régime de TVA sur la marge
				  (CGI art. 266-1-e, 268, 297 A).
			`),
		},
		Values: []*cbc.Definition{
			{Code: B2CCategoryGoods, Name: i18n.String{i18n.EN: "Goods subject to VAT", i18n.FR: "Livraisons de biens soumises à la TVA"}},
			{Code: B2CCategoryServices, Name: i18n.String{i18n.EN: "Services subject to VAT", i18n.FR: "Prestations de services soumises à la TVA"}},
			{Code: B2CCategoryNotTaxable, Name: i18n.String{i18n.EN: "Not subject to French VAT", i18n.FR: "Non soumis à la TVA en France"}},
			{Code: B2CCategoryMargin, Name: i18n.String{i18n.EN: "VAT-on-margin regime", i18n.FR: "Régime de TVA sur la marge"}},
		},
	},
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

// reasonCodeDefinitions builds the value list for the fr-ctc-reason-code
// extension from the authoritative reasonTable — avoids drift between
// the helper table and the extension's accepted value set.
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
// ProcessConditionCode.
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

// statusCodeDefinitions builds the value list for fr-ctc-status-code
// from the authoritative processTable.
func statusCodeDefinitions() []*cbc.Definition {
	out := make([]*cbc.Definition, 0, len(processTable))
	for _, e := range processTable {
		out = append(out, &cbc.Definition{
			Code: cbc.Code(e.Code),
			Name: i18n.String{i18n.EN: processCodeLabels[e.Code]},
		})
	}
	return out
}
