package choruspro

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

const (
	// ExtKeyFramework is the key for the general information framework.
	ExtKeyFramework cbc.Key = "fr-choruspro-framework"
	// ExtKeyScheme is the key for the scheme.
	ExtKeyScheme cbc.Key = "fr-choruspro-scheme"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyFramework,
		Name: i18n.String{
			i18n.EN: "General Information Framework",
			i18n.FR: "Informations générales",
		},
		Desc: i18n.String{
			i18n.EN: "Due to the complexity of the values, GOBL will not apply scenarios.This means that during normalization the extension will be set to A1 if not present. This behavior is not deterministic and goes against GOBL.",
		},
		Values: []*cbc.Definition{
			{
				Code: "A1",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier of an invoice",
					i18n.FR: "Dépôt par un fournisseur d'une facture",
				},
			},
			{
				Code: "A2",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier of an invoice already paid",
					i18n.FR: "Dépôt par un fournisseur d'une facture déjà payée",
				},
			},
			{
				Code: "A3",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier of a Memorandum on Justice Costs",
					i18n.FR: "Dépôt par un fournisseur d'un mémoire de frais de justice",
				},
			},
			{
				Code: "A4",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier of a draft monthly statement",
					i18n.FR: "Dépôt par un fournisseur d'un projet de décompte mensuel",
				},
			},
			{
				Code: "A5",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier of an account statement",
					i18n.FR: "Dépôt par un fournisseur d'un état d'acompte",
				},
			},
			{
				Code: "A6",
				Name: i18n.String{
					i18n.EN: "Work invoice document sent to a financial service",
					i18n.FR: "Pièce de facturation de travaux transmise au service financier",
				},
			},
			{
				Code: "A7",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier of a draft final statement",
					i18n.FR: "Dépôt par un fournisseur d'un projet de décompte final",
				},
			},
			{
				Code: "A8",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier of a general and definitive statement",
					i18n.FR: "Dépôt par un fournisseur d'un décompte général et définitif",
				},
			},
			{
				Code: "A9",
				Name: i18n.String{
					i18n.EN: "Submission by a subcontractor of an invoice",
					i18n.FR: "Dépôt par un sous-traitant d'une facture",
				},
			},
			{
				Code: "A10",
				Name: i18n.String{
					i18n.EN: "Submission by a subcontractor of a draft monthly statement",
					i18n.FR: "Dépôt par un sous-traitant d'un projet de décompte mensuel",
				},
			},
			{
				Code: "A12",
				Name: i18n.String{
					i18n.EN: "Submission by a joint contractor of an invoice",
					i18n.FR: "Dépôt par un cotraitant d'une facture",
				},
			},
			{
				Code: "A13",
				Name: i18n.String{
					i18n.EN: "Submission by a joint contractor of a draft monthly statement",
					i18n.FR: "Dépôt par un cotraitant d'un projet de décompte mensuel",
				},
			},
			{
				Code: "A14",
				Name: i18n.String{
					i18n.EN: "Submission by a joint contractor of a draft final statement",
					i18n.FR: "Dépôt par un cotraitant d'un projet de décompte final",
				},
			},
			{
				Code: "A15",
				Name: i18n.String{
					i18n.EN: "Submission by a project manager of an account statement",
					i18n.FR: "Dépôt par une MOE d'un état d'acompte",
				},
			},
			{
				Code: "A16",
				Name: i18n.String{
					i18n.EN: "Submission by a project manager of a validated account statement",
					i18n.FR: "Dépôt par une MOE d'un état d'acompte validé",
				},
			},
			{
				Code: "A17",
				Name: i18n.String{
					i18n.EN: "Submission by a project manager of a draft general statement",
					i18n.FR: "Dépôt par une MOE d'un projet de décompte général",
				},
			},
			{
				Code: "A18",
				Name: i18n.String{
					i18n.EN: "Submission by a project manager of a general statement",
					i18n.FR: "Dépôt par une MOE d'un décompte général",
				},
			},
			{
				Code: "A19",
				Name: i18n.String{
					i18n.EN: "Submission by a contracting authority of a validated account statement",
					i18n.FR: "Dépôt par une MOA d'un état d'acompte validé",
				},
			},
			{
				Code: "A20",
				Name: i18n.String{
					i18n.EN: "Submission by a contracting authority of a general statement",
					i18n.FR: "Dépôt par une MOA d'un décompte général",
				},
			},
			{
				Code: "A21",
				Name: i18n.String{
					i18n.EN: "Submission by a beneficiary of an ICT reimbursement request",
					i18n.FR: "Dépôt par un bénéficiaire d'une demande de remboursement de la TIC",
				},
			},
			{
				Code: "A22",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier or an authorised representative of a draft general statement as part of a tacit procedure",
					i18n.FR: "Dépôt par un fournisseur ou mandataire d'un projet de décompte général dans le cadre d'une procédure tacite",
				},
			},
			{
				Code: "A23",
				Name: i18n.String{
					i18n.EN: "Submission by a supplier or an authorised representative of a tacit general and final statement",
					i18n.FR: "Dépôt par un fournisseur ou mandataire d'un décompte général et définitif tacite",
				},
			},
			{
				Code: "A24",
				Name: i18n.String{
					i18n.EN: "Submission by an authorised representative of a tacit general and final statement",
					i18n.FR: "Dépôt par une MOE d'un décompte général et définitif tacite",
				},
			},
			{
				Code: "A25",
				Name: i18n.String{
					i18n.EN: "Submission by an authorised representative of a general and final statement as part of a tacit procedure",
					i18n.FR: "Dépôt par une MOA d'un décompte général et définitif dans le cadre d'une procédure tacite",
				},
			},
		},
	},

	{
		Key: ExtKeyScheme,
		Name: i18n.String{
			i18n.EN: "Scheme",
			i18n.FR: "Schéma",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Third party with SIRET",
					i18n.FR: "Tiers avec SIRET",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "European structure outside France",
					i18n.FR: "Structure Européenne hors France",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Structure outside the EU",
					i18n.FR: "Structure hors UE",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "RIDET",
					i18n.FR: "RIDET",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Tahiti Number",
					i18n.FR: "Numéro Tahiti",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.FR: "Autre",
				},
			},
		},
	},
}
