package nfse

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Brazilian extension keys required to issue NFS-e documents.
const (
	ExtKeyFiscalIncentive = "br-nfse-fiscal-incentive"
	ExtKeyMunicipality    = "br-nfse-municipality"
	ExtKeyService         = "br-nfse-service"
	ExtKeySimplesNacional = "br-nfse-simples-nacional"
	ExtKeySpecialRegime   = "br-nfse-special-regime"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyFiscalIncentive,
		Name: i18n.String{
			i18n.EN: "Fiscal Incentive",
			i18n.PT: "Incentivo Fiscal",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
				Name: i18n.String{
					i18n.EN: "Has incentive",
					i18n.PT: "Possui incentivo",
				},
			},
			{
				Value: "2",
				Name: i18n.String{
					i18n.EN: "Does not have incentive",
					i18n.PT: "Não possui incentivo",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates whether a party benefits from a fiscal incentive.

				List of codes taken from the national NFSe standard:
				https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download
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

				For further details on the list of possible codes, see:

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
		Key: ExtKeySimplesNacional,
		Name: i18n.String{
			i18n.EN: "Opting for “Simples Nacional”",
			i18n.PT: "Optante pelo Simples Nacional",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
				Name: i18n.String{
					i18n.EN: "Opt-in",
					i18n.PT: "Optante",
				},
			},
			{
				Value: "2",
				Name: i18n.String{
					i18n.EN: "Opt-out",
					i18n.PT: "Não optante",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates whether a party is opting for the “Simples Nacional” tax regime.

				List of codes taken from the national NFSe standard:
				https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download
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
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
				Name: i18n.String{
					i18n.EN: "Municipal micro-enterprise",
					i18n.PT: "Microempresa municipal",
				},
			},
			{
				Value: "2",
				Name: i18n.String{
					i18n.EN: "Estimated",
					i18n.PT: "Estimativa",
				},
			},
			{
				Value: "3",
				Name: i18n.String{
					i18n.EN: "Professional Society",
					i18n.PT: "Sociedade de profissionais",
				},
			},
			{
				Value: "4",
				Name: i18n.String{
					i18n.EN: "Cooperative",
					i18n.PT: "Cooperativa",
				},
			},
			{
				Value: "5",
				Name: i18n.String{
					i18n.EN: "Single micro-entrepreneur (MEI)",
					i18n.PT: "Microempreendedor individual (MEI)",
				},
			},
			{
				Value: "6",
				Name: i18n.String{
					i18n.EN: "Micro-enterprise or Small Business (ME EPP)",
					i18n.PT: "Microempresa ou Empresa de Pequeno Porte (ME EPP).",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates a special tax regime that the party is subject to.

				List of codes taken from the national NFSe standard:
				https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e-manual-de-orientacao-do-contribuinte-2-04/download
				(Section 10.2, Field B-66)
			`),
		},
	},
}
