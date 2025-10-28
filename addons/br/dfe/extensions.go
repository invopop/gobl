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
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "NFS-e ABRASF Taxpayer Guidance Manual (v2.04)",
					i18n.PT: "NFS-e ABRASF Manual de Orientação do Contribuinte (v2.04)",
				},
				URL: "https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download",
			},
		},
	},
	{
		Key: ExtKeyMunicipality,
		Name: i18n.String{
			i18n.EN: "IBGE Municipality Code",
			i18n.PT: "Código do Município do IBGE",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The municipality code as defined by the IBGE (Brazilian Institute of Geography and
				Statistics).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "IBGE - Municipalities Codes",
					i18n.PT: "IBGE - Códigos dos Municípios",
				},
				URL: "https://www.ibge.gov.br/explica/codigos-dos-municipios.php",
			},
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
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "NFS-e ABRASF Taxpayer Guidance Manual (v2.04)",
					i18n.PT: "NFS-e ABRASF Manual de Orientação do Contribuinte (v2.04)",
				},
				URL: "https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download",
			},
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
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "NFS-e ABRASF Taxpayer Guidance Manual (v2.04)",
					i18n.PT: "NFS-e ABRASF Manual de Orientação do Contribuinte (v2.04)",
				},
				URL: "https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download",
			},
		},
	},
}
