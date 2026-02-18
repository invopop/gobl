// Package ctc handles the extensions and validation rules for the French
// CTC (Cycle de Traitement de la Commande) Flow 2 B2B e-invoicing requirements.
package ctc

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// Flow2V1 is the key for the French CTC addon
	Flow2V1 cbc.Key = "fr-ctc-flow2-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: Flow2V1,
		Name: i18n.String{
			i18n.EN: "France CTC Flow 2",
			i18n.FR: "France CTC Flux 2",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French CTC (Continuous Transaction Controls) Flow 2 B2B
				e-invoicing mandate from the French electronic invoicing reform.

				This addon extends the EN 16931 European standard with French-specific rules for
				regulated B2B invoices — that is, invoices exchanged between two parties registered
				for VAT in France. It should not be used for invoices that are subject only to
				e-reporting (for example, B2C or cross-border transactions).

				## Billing Mode

				Every invoice must carry a **billing mode** extension (` + "`" + `fr-ctc-billing-mode` + "`" + `) that
				describes the nature of the supply and the payment context. The code is made up of
				a letter prefix and a numeric suffix:

				| Prefix | Meaning |
				|--------|---------|
				| B | Goods (Biens) |
				| S | Services |
				| M | Mixed (goods and services that are not accessory to each other) |

				| Suffix | Meaning |
				|--------|---------|
				| 1 | Standard deposit invoice |
				| 2 | Already-paid invoice (payable amount = 0) |
				| 4 | Final invoice after one or more advance payments |
				| 5 | Subcontractor invoice (services only) |
				| 6 | Co-contractor invoice (services only) |
				| 7 | E-reporting invoice (VAT already collected) |

				## Identities

				A **SIREN** identity (9 digits) must be present on both the supplier and the customer
				for B2B invoices. The **SIRET** (14 digits) is optional; if provided without a SIREN,
				the addon will automatically derive the SIREN from the first nine digits of the SIRET.

				The SIREN is automatically assigned the ` + "`" + `legal` + "`" + ` scope when no other identity on
				the party already carries that scope.

				Private identifiers can be included using the ` + "`" + `private-id` + "`" + ` key; the addon will
				assign ISO scheme ID ` + "`" + `0224` + "`" + ` to these automatically.

				## Electronic Addresses (Inboxes)

				Both the supplier and customer must have an electronic address for B2B invoices.
				Addresses using SIREN scheme ` + "`" + `0225` + "`" + ` must contain only alphanumeric characters
				and the symbols ` + "`" + `+` + "`" + `, ` + "`" + `-` + "`" + `, ` + "`" + `_` + "`" + `, ` + "`" + `/` + "`" + `. If the party has a SIREN identity and no inbox
				carries a ` + "`" + `peppol` + "`" + ` key, the addon will assign the ` + "`" + `peppol` + "`" + ` key to the SIREN inbox
				automatically.

				## Required Notes

				Three invoice-level notes with specific UNTDID text-subject codes are mandatory:

				| Code | Purpose |
				|------|---------|
				| ` + "`" + `PMT` + "`" + ` | Payment terms |
				| ` + "`" + `PMD` + "`" + ` | Payment means |
				| ` + "`" + `AAB` + "`" + ` | Buyer accounting reference (BAR) |

				The BAR note must contain one of the following treatment values: ` + "`" + `B2B` + "`" + `, ` + "`" + `B2BINT` + "`" + `,
				` + "`" + `B2C` + "`" + `, ` + "`" + `OUTOFSCOPE` + "`" + `, or ` + "`" + `ARCHIVEONLY` + "`" + `.

				## Invoice Code

				The invoice series and code must each be at most 35 characters and may only contain
				alphanumeric characters and the symbols ` + "`" + `-` + "`" + `, ` + "`" + `+` + "`" + `, ` + "`" + `_` + "`" + `, ` + "`" + `/` + "`" + `.

				## Currency Conversion

				When the invoice is issued in a currency other than EUR, the gobl.ubl conversion
				library will automatically add EUR equivalents for the tax totals, satisfying the
				BR-FR-CO-12 requirement without any extra input from the user.
			`),
			i18n.FR: here.Doc(`
				Support pour le mandat de facturation électronique CTC (Contrôle Continu des
				Transactions) français Flux 2 B2B, issu de la réforme française de la facturation
				électronique.

				Cet addon étend le standard européen EN 16931 avec des règles spécifiques à la
				France pour les factures B2B réglementées — c'est-à-dire les factures échangées
				entre deux parties assujetties à la TVA en France. Il ne doit pas être utilisé pour
				les factures soumises uniquement à l'e-reporting (par exemple, B2C ou transactions
				transfrontalières).

				## Cadre de facturation

				Chaque facture doit comporter une extension **cadre de facturation**
				(` + "`" + `fr-ctc-billing-mode` + "`" + `) décrivant la nature de la prestation et le contexte de
				paiement. Le code est composé d'un préfixe lettré et d'un suffixe numérique :

				| Préfixe | Signification |
				|---------|---------------|
				| B | Biens |
				| S | Services |
				| M | Mixte (biens et services non accessoires les uns aux autres) |

				| Suffixe | Signification |
				|---------|---------------|
				| 1 | Facture de dépôt standard |
				| 2 | Facture déjà payée (montant exigible = 0) |
				| 4 | Facture définitive après un ou plusieurs acomptes |
				| 5 | Facture de sous-traitance (services uniquement) |
				| 6 | Facture de cotraitance (services uniquement) |
				| 7 | Facture e-reporting (TVA déjà collectée) |

				## Identifiants

				Un identifiant **SIREN** (9 chiffres) doit être présent chez le fournisseur et le
				client pour les factures B2B. Le **SIRET** (14 chiffres) est facultatif ; s'il est
				fourni sans SIREN, l'addon dérive automatiquement le SIREN des neuf premiers chiffres
				du SIRET.

				Le SIREN reçoit automatiquement la portée ` + "`" + `legal` + "`" + ` lorsqu'aucun autre identifiant de
				la partie ne porte déjà cette portée.

				Les identifiants privés peuvent être inclus avec la clé ` + "`" + `private-id` + "`" + ` ; l'addon leur
				attribue automatiquement l'identifiant de schéma ISO ` + "`" + `0224` + "`" + `.

				## Adresses électroniques (boîtes de réception)

				Le fournisseur et le client doivent tous deux disposer d'une adresse électronique
				pour les factures B2B. Les adresses utilisant le schéma SIREN ` + "`" + `0225` + "`" + ` ne doivent
				contenir que des caractères alphanumériques et les symboles ` + "`" + `+` + "`" + `, ` + "`" + `-` + "`" + `, ` + "`" + `_` + "`" + `, ` + "`" + `/` + "`" + `. Si la
				partie possède un identifiant SIREN et qu'aucune boîte de réception ne porte la clé
				` + "`" + `peppol` + "`" + `, l'addon assigne automatiquement cette clé à la boîte SIREN.

				## Notes obligatoires

				Trois notes au niveau de la facture avec des codes objet de texte UNTDID spécifiques
				sont obligatoires :

				| Code | Objet |
				|------|-------|
				| ` + "`" + `PMT` + "`" + ` | Conditions de paiement |
				| ` + "`" + `PMD` + "`" + ` | Moyens de paiement |
				| ` + "`" + `AAB` + "`" + ` | Référence comptable acheteur (BAR) |

				La note BAR doit contenir l'une des valeurs de traitement suivantes : ` + "`" + `B2B` + "`" + `,
				` + "`" + `B2BINT` + "`" + `, ` + "`" + `B2C` + "`" + `, ` + "`" + `OUTOFSCOPE` + "`" + ` ou ` + "`" + `ARCHIVEONLY` + "`" + `.

				## Code de facture

				La série et le code de la facture doivent chacun comporter au maximum 35 caractères
				et ne peuvent contenir que des caractères alphanumériques et les symboles ` + "`" + `-` + "`" + `, ` + "`" + `+` + "`" + `,
				` + "`" + `_` + "`" + `, ` + "`" + `/` + "`" + `.

				## Conversion de devises

				Lorsque la facture est émise dans une devise autre que l'EUR, la bibliothèque de
				conversion gobl.ubl ajoute automatiquement les équivalents en EUR pour les totaux
				de TVA, satisfaisant l'exigence BR-FR-CO-12 sans intervention supplémentaire de
				l'utilisateur.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "External Specifications",
					i18n.FR: "Spécifications Externes",
				},
				URL: "https://www.impots.gouv.fr/specifications-externes-b2b",
			},
		},
		Extensions: extensions,
		Normalizer: normalize,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *org.Party:
		normalizeParty(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *org.Party:
		return validateParty(obj)
	case *org.Identity:
		return validateIdentity(obj)
	case *org.Inbox:
		return validateInbox(obj)
	case *org.Item:
		return validateItem(obj)
	}
	return nil
}
