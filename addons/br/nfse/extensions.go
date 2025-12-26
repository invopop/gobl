package nfse

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Brazilian extension keys required to issue NFS-e documents
const (
	ExtKeyCNAE            = "br-nfse-cnae"
	ExtKeyFiscalIncentive = "br-nfse-fiscal-incentive"
	ExtKeyISSLiability    = "br-nfse-iss-liability"
	ExtKeyService         = "br-nfse-service"
	ExtKeySimples         = "br-nfse-simples"
	ExtKeySpecialRegime   = "br-nfse-special-regime"
	ExtKeyOperation       = "br-nfse-operation"
	ExtKeyTaxStatus       = "br-nfse-tax-status"
	ExtKeyTaxClass        = "br-nfse-tax-class"
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
	{
		Key: ExtKeyOperation,
		Name: i18n.String{
			i18n.EN: "Operation Indicator",
			i18n.PT: "Indicador da operação",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates the operation type for the determination of the IBS and CBS taxes by the
				tax authorities.

				Maps to the ~cIndOp~ field in the NFS-e national layout.

				List of possible values:

				* https://www.gov.br/nfse/pt-br/biblioteca/documentacao-tecnica/rtc/anexovii-indop_ibscbs_v1-00-00.xlsx
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical Note SE/CGNFS-e nº 004",
					i18n.PT: "Nota Técnica SE/CGNFS-e nº 004",
				},
				URL: "https://www.gov.br/nfse/pt-br/biblioteca/documentacao-tecnica/rtc-producao-restrita-piloto/nt-004-se-cgnfse-novo-layout-rtc-v2-00-20251210.pdf",
			},
			{
				Title: i18n.String{
					i18n.EN: "Annex VI - Operation Indicators Table IBS/CBS v1.00",
					i18n.PT: "Anexo VI - Tabela de Indicadores de Operação IBS/CBS v1.00",
				},
				URL: "https://www.gov.br/nfse/pt-br/biblioteca/documentacao-tecnica/rtc/anexovii-indop_ibscbs_v1-00-00.xlsx",
			},
		},
		Pattern: `^\d{6}$`,
	},
	{
		Key: ExtKeyTaxStatus,
		Name: i18n.String{
			i18n.EN: "Tax Status Code (CST)",
			i18n.PT: "Código de situação tributária (CST)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates the tax status of the operation for the determination of the IBS and CBS
				taxes by the tax authorities.

				Maps to the ~CST~ field in the NFS-e national layout.

				List of possible values:

				* https://dfe-portal.svrs.rs.gov.br/DFE/ClassificacaoTributaria
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical Report RT 2025.002 - CST and cClassTrib Tables IBS/CBS",
					i18n.PT: "Informe Técnico RT 2025.002 - Tabelas CST e cClassTrib IBS/CBS",
				},
				URL: "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=gya58CS0dHU=",
			},
		},
		Pattern: `^\d{3}$`,
	},
	{
		Key: ExtKeyTaxClass,
		Name: i18n.String{
			i18n.EN: "Tax Classification Code",
			i18n.PT: "Código de classificação tributária",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates the tax classification code for the determination of the IBS and CBS taxes
				by the tax authorities.

				Maps to the ~cClassTrib~ field in the NFS-e national layout.

				List of possible values:

				* https://dfe-portal.svrs.rs.gov.br/DFE/ClassificacaoTributaria
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical Report RT 2025.002 - CST and cClassTrib Tables IBS/CBS",
					i18n.PT: "Informe Técnico RT 2025.002 - Tabelas CST e cClassTrib IBS/CBS",
				},
				URL: "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=gya58CS0dHU=",
			},
		},
		Pattern: `^\d{6}$`,
	},
}
