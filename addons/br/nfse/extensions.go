package nfse

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Brazilian extension keys required to issue NFS-e documents. In an initial
// assessment, these extensions do not seem to apply to documents other than
// NFS-e. However, if when implementing other Fiscal Notes it is found that some
// of these extensions are common, they can be moved to the regime or to a
// shared addon.
const (
	ExtKeyCNAE            = "br-nfse-cnae"
	ExtKeyFiscalIncentive = "br-nfse-fiscal-incentive"
	ExtKeyISSLiability    = "br-nfse-iss-liability"
	ExtKeyMunicipality    = "br-nfse-municipality"
	ExtKeyService         = "br-nfse-service"
	ExtKeySimples         = "br-nfse-simples"
	ExtKeySpecialRegime   = "br-nfse-special-regime"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyCNAE,
		Name: i18n.String{
			i18n.EN: "CNAE code",
			i18n.PT: "Código CNAE",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The CNAE (National Classification of Economic Activities) code for a service.

				List of codes from the IBGE (Brazilian Institute of Geography and Statistics):

				* https://www.ibge.gov.br/en/statistics/technical-documents/statistical-lists-and-classifications/17245-national-classification-of-economic-activities.html
			`),
		},
		Pattern: `^\d{2}[\s\.\-\/]?\d{2}[\s\.\-\/]?\d[\s\.\-\/]?\d{2}$`,
	},
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
		Key: ExtKeyISSLiability,
		Name: i18n.String{
			i18n.EN: "ISS Liability",
			i18n.PT: "Exigibilidade ISS",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Liable",
					i18n.PT: "Exigível",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Not subject",
					i18n.PT: "Não incidência",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PT: "Isenção",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.PT: "Exportação",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Immune",
					i18n.PT: "Imunidade",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Suspended Judicially",
					i18n.PT: "Suspensa por Decisão Judicial",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Suspended Administratively",
					i18n.PT: "Suspensa por Processo Administrativo",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates the ISS liability status, i.e., whether the ISS tax is due or not and why.

				List of codes from the national NFSe ABRASF (v2.04) model:

				* https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download
				(Section 10.2, Field B-38)
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
		Key: ExtKeyService,
		Name: i18n.String{
			i18n.EN: "Service Code",
			i18n.PT: "Código Item Lista Serviço",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The service code as defined by the municipality. Typically, one of the codes listed
				in the Lei Complementar 116/2003, but municipalities can make their own changes.

				For further details on the list of possible codes, see:

				* https://www.planalto.gov.br/ccivil_03/leis/lcp/lcp116.htm
			`),
		},
	},
	{
		Key: ExtKeySimples,
		Name: i18n.String{
			i18n.EN: "Opting for “Simples Nacional” regime",
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
				Indicates whether a party is opting for the “Simples Nacional” (Regime Especial
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
}
