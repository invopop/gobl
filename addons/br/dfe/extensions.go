package dfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Brazilian extension keys for fiscal documents
const (
	ExtKeyFiscalIncentive = "br-dfe-fiscal-incentive"
	ExtKeyMunicipality    = "br-dfe-municipality"
	ExtKeySimples         = "br-dfe-simples"
	ExtKeySpecialRegime   = "br-dfe-special-regime"
	ExtKeyCFOP            = "br-dfe-cfop"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyFiscalIncentive,
		Name: i18n.String{
			i18n.EN: "Fiscal Incentive",
			i18n.PT: "Incentivo Fiscal",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Has incentive",
					i18n.PT: "Possui incentivo",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Does not have incentive",
					i18n.PT: "Não possui incentivo",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates whether a party benefits from a fiscal incentive.

				List of codes from the national NFSe ABRASF (v2.04) model:

				* https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download
				(Section 10.2, Field B-68)
			`),
		},
	},
	{
		Key: ExtKeyMunicipality,
		Name: i18n.String{
			i18n.EN: "IGBE Municipality Code",
			i18n.PT: "Código do Município do IBGE",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The municipality code as defined by the IGBE (Brazilian Institute of Geography and
				Statistics).

				List of codes from the IGBE:

				* https://www.ibge.gov.br/explica/codigos-dos-municipios.php
			`),
		},
		Pattern: `^\d{7}$`,
	},
	{
		Key: ExtKeySimples,
		Name: i18n.String{
			i18n.EN: "Opting for \"Simples Nacional\" regime",
			i18n.PT: "Optante pelo Simples Nacional",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Opt-in",
					i18n.PT: "Optante",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Opt-out",
					i18n.PT: "Não optante",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates whether a party is opting for the "Simples Nacional" (Regime Especial
				Unificado de Arrecadação de Tributos e Contribuições devidos pelas Microempresas e
				Empresas de Pequeno Porte) tax regime

				List of codes from the national NFSe ABRASF (v2.04) model:

				* https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download
				(Section 10.2, Field B-67)
			`),
		},
	},
	{
		Key: ExtKeySpecialRegime,
		Name: i18n.String{
			i18n.EN: "Special Tax Regime",
			i18n.PT: "Regime Especial de Tributação",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Municipal micro-enterprise",
					i18n.PT: "Microempresa municipal",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Estimated",
					i18n.PT: "Estimativa",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Professional Society",
					i18n.PT: "Sociedade de profissionais",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Cooperative",
					i18n.PT: "Cooperativa",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Single micro-entrepreneur (MEI)",
					i18n.PT: "Microempreendedor individual (MEI)",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Micro-enterprise or Small Business (ME EPP)",
					i18n.PT: "Microempresa ou Empresa de Pequeno Porte (ME EPP).",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates a special tax regime that a party is subject to.

				List of codes from the national NFSe ABRASF (v2.04) model:

				* https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download
				(Section 10.2, Field B-66)
			`),
		},
	},
	{
		Key: ExtKeyCFOP,
		Name: i18n.String{
			i18n.EN: "CFOP (Fiscal Operations and Services Code)",
			i18n.PT: "CFOP (Código Fiscal de Operações e Prestações)",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "CFOP - Fiscal Operations and Services Codes (SEFAZ-PE)",
					i18n.PT: "CFOP - Código Fiscal de Operações e Prestações (SEFAZ-PE)",
				},
				URL:         "https://www.sefaz.pe.gov.br/legislacao/tributaria/documents/legislacao/tabelas/cfop.htm",
				ContentType: "text/html",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Four-digit code that classifies the nature of goods movements and service
				provisions for ICMS purposes in Brazil. The first digit indicates the
				operation origin/destination (1–3 for entries; 5–7 for exits), and the
				remaining digits identify the specific type of operation.
			`),
			i18n.PT: here.Doc(`
				Código de quatro dígitos que classifica a natureza das operações de
				circulação de mercadorias e das prestações de serviços para fins de ICMS.
				O primeiro dígito indica a origem/destino da operação (1–3 para entradas;
				5–7 para saídas) e os demais identificam o tipo específico de operação.
			`),
		},
		Pattern: `^[1-7]\d{3}$`,
	},
}
