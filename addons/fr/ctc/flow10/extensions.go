package flow10

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// French CTC extension keys for e-reporting
const (
	// ExtKeyBillingMode defines the billing framework mode (B1-B8, S1-S8, M1-M8).
	// Shared conceptually with Flow 2: both flows use the same underlying key
	// "fr-ctc-billing-mode" and identical value set; each addon declares the
	// definition independently so consumers can opt into either flow in
	// isolation.
	ExtKeyBillingMode cbc.Key = "fr-ctc-billing-mode"

	// ExtKeyB2CCategory classifies the B2C transaction for PPF reporting
	// (G1.68). Required on B2C invoices and B2C payments.
	ExtKeyB2CCategory cbc.Key = "fr-ctc-b2c-category"
)

// B2C transaction category codes (G1.68)
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

// Billing mode codes (Cadre de Facturation)
// The prefix indicates the invoice nature:
//   - B: Goods invoice (Biens)
//   - S: Services invoice
//   - M: Mixed/dual invoice (goods and services that are not accessory to each other)
const (
	// BillingModeB1: Deposit of a goods invoice
	BillingModeB1 cbc.Code = "B1"
	// BillingModeB2: Deposit of an already paid goods invoice
	BillingModeB2 cbc.Code = "B2"
	// BillingModeB4: Deposit of a final invoice (after down payment) for goods
	BillingModeB4 cbc.Code = "B4"
	// BillingModeB7: Deposit of a goods invoice subject to e-reporting (VAT already collected)
	BillingModeB7 cbc.Code = "B7"
	// BillingModeS1: Deposit of a service invoice
	BillingModeS1 cbc.Code = "S1"
	// BillingModeS2: Deposit of an already paid service invoice
	BillingModeS2 cbc.Code = "S2"
	// BillingModeS4: Deposit of a final invoice (after down payment) for services
	BillingModeS4 cbc.Code = "S4"
	// BillingModeS5: Deposit by a subcontractor of a service invoice
	BillingModeS5 cbc.Code = "S5"
	// BillingModeS6: Deposit by a co-contractor of a service invoice
	BillingModeS6 cbc.Code = "S6"
	// BillingModeS7: Deposit of a service invoice subject to e-reporting (VAT already collected)
	BillingModeS7 cbc.Code = "S7"
	// BillingModeM1: Deposit of a dual invoice (goods and services)
	BillingModeM1 cbc.Code = "M1"
	// BillingModeM2: Deposit of an already paid dual invoice
	BillingModeM2 cbc.Code = "M2"
	// BillingModeM4: Deposit of a final invoice (after down payment) - dual
	BillingModeM4 cbc.Code = "M4"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyB2CCategory,
		Name: i18n.String{
			i18n.EN: "B2C Transaction Category",
			i18n.FR: "Catégorie de transaction B2C",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Classifies a B2C transaction for French Flow 10 reporting to the PPF
				(G1.68). Required on B2C invoices and B2C payments.

				- TLB1: Goods deliveries subject to VAT.
				- TPS1: Services subject to VAT.
				- TNT1: Goods / services not subject to French VAT, including
				  intra-EU distance sales per CGI articles 258 A and 259 B.
				- TMA1: Operations under the VAT-on-margin regime
				  (CGI articles 266-1-e, 268, 297 A).
			`),
			i18n.FR: here.Doc(`
				Catégorie de transaction pour le reporting Flux 10 au PPF (G1.68).
				Obligatoire sur les factures et paiements B2C.

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
			{
				Code: B2CCategoryGoods,
				Name: i18n.String{
					i18n.EN: "Goods subject to VAT",
					i18n.FR: "Livraisons de biens soumises à la TVA",
				},
			},
			{
				Code: B2CCategoryServices,
				Name: i18n.String{
					i18n.EN: "Services subject to VAT",
					i18n.FR: "Prestations de services soumises à la TVA",
				},
			},
			{
				Code: B2CCategoryNotTaxable,
				Name: i18n.String{
					i18n.EN: "Not subject to French VAT",
					i18n.FR: "Non soumis à la TVA en France",
				},
			},
			{
				Code: B2CCategoryMargin,
				Name: i18n.String{
					i18n.EN: "VAT-on-margin regime",
					i18n.FR: "Régime de TVA sur la marge",
				},
			},
		},
	},
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
}
