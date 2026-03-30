package ctc

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// French CTC extension keys for B2B e-invoicing
const (
	// ExtKeyBillingMode defines the billing framework mode (B1-B8, S1-S8, M1-M8)
	ExtKeyBillingMode cbc.Key = "fr-ctc-billing-mode"
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

				The value is currently trusted as provided and is not normalised from other invoice
				fields. Automatic inference from document type and payment context may be added in
				a future release.
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

				La valeur est actuellement utilisée telle quelle et n'est pas normalisée à partir
				des autres champs de la facture. Une inférence automatique à partir du type de
				document et du contexte de paiement pourra être ajoutée dans une version future.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: BillingModeB1,
				Name: i18n.String{
					i18n.EN: "Goods - Deposit invoice",
					i18n.FR: "Biens - Facture de dépôt",
				},
			},
			{
				Code: BillingModeB2,
				Name: i18n.String{
					i18n.EN: "Goods - Already paid invoice",
					i18n.FR: "Biens - Facture déjà payée",
				},
			},
			{
				Code: BillingModeB4,
				Name: i18n.String{
					i18n.EN: "Goods - Final invoice (after down payment)",
					i18n.FR: "Biens - Facture définitive (après acompte)",
				},
			},
			{
				Code: BillingModeB7,
				Name: i18n.String{
					i18n.EN: "Goods - E-reporting (VAT already collected)",
					i18n.FR: "Biens - E-reporting (TVA déjà collectée)",
				},
			},
			{
				Code: BillingModeS1,
				Name: i18n.String{
					i18n.EN: "Services - Deposit invoice",
					i18n.FR: "Services - Facture de dépôt",
				},
			},
			{
				Code: BillingModeS2,
				Name: i18n.String{
					i18n.EN: "Services - Already paid invoice",
					i18n.FR: "Services - Facture déjà payée",
				},
			},
			{
				Code: BillingModeS4,
				Name: i18n.String{
					i18n.EN: "Services - Final invoice (after down payment)",
					i18n.FR: "Services - Facture définitive (après acompte)",
				},
			},
			{
				Code: BillingModeS5,
				Name: i18n.String{
					i18n.EN: "Services - Subcontractor invoice",
					i18n.FR: "Services - Facture de sous-traitance",
				},
			},
			{
				Code: BillingModeS6,
				Name: i18n.String{
					i18n.EN: "Services - Co-contractor invoice",
					i18n.FR: "Services - Facture de cotraitance",
				},
			},
			{
				Code: BillingModeS7,
				Name: i18n.String{
					i18n.EN: "Services - E-reporting (VAT already collected)",
					i18n.FR: "Services - E-reporting (TVA déjà collectée)",
				},
			},
			{
				Code: BillingModeM1,
				Name: i18n.String{
					i18n.EN: "Mixed - Deposit invoice",
					i18n.FR: "Mixte - Facture de dépôt",
				},
			},
			{
				Code: BillingModeM2,
				Name: i18n.String{
					i18n.EN: "Mixed - Already paid invoice",
					i18n.FR: "Mixte - Facture déjà payée",
				},
			},
			{
				Code: BillingModeM4,
				Name: i18n.String{
					i18n.EN: "Mixed - Final invoice (after down payment)",
					i18n.FR: "Mixte - Facture définitive (après acompte)",
				},
			},
		},
	},
}
