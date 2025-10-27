package nfse

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Brazilian extension keys required to issue NFS-e documents
const (
	ExtKeyCNAE         = "br-nfse-cnae"
	ExtKeyISSLiability = "br-nfse-iss-liability"
	ExtKeyService      = "br-nfse-service"
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
}
