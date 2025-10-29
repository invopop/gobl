package nfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// NF-e Extension Keys
const (
	ExtKeyModel           = "br-nfe-model"
	ExtKeyPresence        = "br-nfe-presence"
	ExtKeyPaymentMeans    = "br-nfe-payment-means"
	ExtKeyCFOP            = "br-nfe-cfop"
	ExtKeyFiscalIncentive = "br-nfe-fiscal-incentive"
	ExtKeyRegime          = "br-nfe-regime"
	ExtKeySpecialRegime   = "br-nfe-special-regime"
)

// Model Codes
const (
	ModelNFe  cbc.Code = "55"
	ModelNFCe cbc.Code = "65"
)

// Presence Indicator Codes
const (
	PresenceNotApplicable cbc.Code = "0"
	PresenceInPerson      cbc.Code = "1"
	PresenceInternet      cbc.Code = "2"
	PresenceRemote        cbc.Code = "3"
	PresenceDelivery      cbc.Code = "4"
	PresenceOffsite       cbc.Code = "5"
	PresenceOther         cbc.Code = "9"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyModel,
		Name: i18n.String{
			i18n.EN: "Fiscal Document Model Code",
			i18n.PT: "Código do Modelo do Documento Fiscal",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Taxpayer Guidance Manual v7.0 - Annex I – Layout and Validation Rules for NF-e and NFC-e",
					i18n.PT: "Manual de Orientação ao Contribuinte v7.0 - Anexo I – Leiaute e Regras de Validação da NF-e e da NFC-e",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=J%20I%20v4eN00E=",
				ContentType: "application/pdf",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the fiscal document model. It will be
				determined automatically by GOBL during normalization according to
				the scenario definitions.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ModelNFe,
				Name: i18n.String{
					i18n.EN: "NF-e",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For B2B transfer of goods and certain intra-state
						services.
					`),
				},
			},
			{
				Code: ModelNFCe,
				Name: i18n.String{
					i18n.EN: "NFC-e",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For B2C (retail) transfer of goods and certain
						intra-state services.
					`),
				},
			},
		},
	},
	{
		Key: ExtKeyPresence,
		Name: i18n.String{
			i18n.EN: "Buyer Presence Indicator",
			i18n.PT: "Indicador de Presença do Comprador",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Taxpayer Guidance Manual v7.0 - Annex I – Layout and Validation Rules for NF-e and NFC-e",
					i18n.PT: "Manual de Orientação ao Contribuinte v7.0 - Anexo I – Leiaute e Regras de Validação da NF-e e da NFC-e",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=J%20I%20v4eN00E=",
				ContentType: "application/pdf",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicator of the buyer's presence at the commercial establishment
				at the time of the operation. This field is used to classify the
				type of commercial transaction according to Brazilian tax regulations.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: PresenceNotApplicable,
				Name: i18n.String{
					i18n.EN: "Not applicable",
					i18n.PT: "Não se aplica",
				},
			},
			{
				Code: PresenceInPerson,
				Name: i18n.String{
					i18n.EN: "In-person operation",
					i18n.PT: "Operação presencial",
				},
			},
			{
				Code: PresenceInternet,
				Name: i18n.String{
					i18n.EN: "Non-in-person operation, via Internet",
					i18n.PT: "Operação não presencial, pela Internet",
				},
			},
			{
				Code: PresenceRemote,
				Name: i18n.String{
					i18n.EN: "Non-in-person operation, Tele-service",
					i18n.PT: "Operação não presencial, Teleatendimento",
				},
			},
			{
				Code: PresenceDelivery,
				Name: i18n.String{
					i18n.EN: "NFC-e in operation with home delivery",
					i18n.PT: "NFC-e em operação com entrega a domicílio",
				},
			},
			{
				Code: PresenceOffsite,
				Name: i18n.String{
					i18n.EN: "In-person operation, outside establishment",
					i18n.PT: "Operação presencial, fora do estabelecimento",
				},
			},
			{
				Code: PresenceOther,
				Name: i18n.String{
					i18n.EN: "Non-in-person operation, others",
					i18n.PT: "Operação não presencial, outros",
				},
			},
		},
	},
	{
		Key: ExtKeyPaymentMeans,
		Name: i18n.String{
			i18n.EN: "Payment Method",
			i18n.PT: "Meio de Pagamento",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Taxpayer Guidance Manual v7.0 - Annex I – Layout and Validation Rules for NF-e and NFC-e",
					i18n.PT: "Manual de Orientação ao Contribuinte v7.0 - Anexo I – Leiaute e Regras de Validação da NF-e e da NFC-e",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=J%20I%20v4eN00E=",
				ContentType: "application/pdf",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the payment method used for the transaction.
			`),
			i18n.PT: here.Doc(`
				Código utilizado para identificar o meio de pagamento utilizado na transação.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "Cash",
					i18n.PT: "Dinheiro",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Check",
					i18n.PT: "Cheque",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "Credit Card",
					i18n.PT: "Cartão de Crédito",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "Debit Card",
					i18n.PT: "Cartão de Débito",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "Store Credit",
					i18n.PT: "Crédito Loja",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Food Voucher",
					i18n.PT: "Vale Alimentação",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Meal Voucher",
					i18n.PT: "Vale Refeição",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Gift Voucher",
					i18n.PT: "Vale Presente",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Fuel Voucher",
					i18n.PT: "Vale Combustível",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Bank Slip (Boleto)",
					i18n.PT: "Boleto Bancário",
				},
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "Bank Deposit",
					i18n.PT: "Depósito Bancário",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Instant Payment (PIX)",
					i18n.PT: "Pagamento Instantâneo (PIX)",
				},
			},
			{
				Code: "18",
				Name: i18n.String{
					i18n.EN: "Bank Transfer, Digital Wallet",
					i18n.PT: "Transferência bancária, Carteira Digital",
				},
			},
			{
				Code: "19",
				Name: i18n.String{
					i18n.EN: "Loyalty Program, Cashback, Virtual Credit",
					i18n.PT: "Programa de fidelidade, Cashback, Crédito Virtual",
				},
			},
			{
				Code: "90",
				Name: i18n.String{
					i18n.EN: "No Payment",
					i18n.PT: "Sem pagamento",
				},
			},
			{
				Code: "99",
				Name: i18n.String{
					i18n.EN: "Others",
					i18n.PT: "Outros",
				},
			},
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
	{
		Key: ExtKeyFiscalIncentive,
		Name: i18n.String{
			i18n.EN: "Fiscal Incentive Indicator",
			i18n.PT: "Indicador de Incentivo Fiscal",
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
					i18n.EN: "Taxpayer Guidance Manual v7.0 - Annex I – Layout and Validation Rules for NF-e and NFC-e",
					i18n.PT: "Manual de Orientação ao Contribuinte v7.0 - Anexo I – Leiaute e Regras de Validação da NF-e e da NFC-e",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=J%20I%20v4eN00E=",
				ContentType: "application/pdf",
			},
		},
	},
	{
		Key: ExtKeyRegime,
		Name: i18n.String{
			i18n.EN: "Tax Regime Code",
			i18n.PT: "Código de Regime Tributário",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Simples Nacional",
					i18n.PT: "Simples Nacional",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Simples Nacional, Excess",
					i18n.PT: "Simples Nacional, Excesso",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Normal",
					i18n.PT: "Normal",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates the tax regime that a party is subject to.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Taxpayer Guidance Manual v7.0 - Annex I – Layout and Validation Rules for NF-e and NFC-e",
					i18n.PT: "Manual de Orientação ao Contribuinte v7.0 - Anexo I – Leiaute e Regras de Validação da NF-e e da NFC-e",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=J%20I%20v4eN00E=",
				ContentType: "application/pdf",
			},
		},
	},
	{
		Key: ExtKeySpecialRegime,
		Name: i18n.String{
			i18n.EN: "Special Tax Regime Code",
			i18n.PT: "Código do Regime Especial de Tributação",
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
					i18n.EN: "Taxpayer Guidance Manual v7.0 - Annex I – Layout and Validation Rules for NF-e and NFC-e",
					i18n.PT: "Manual de Orientação ao Contribuinte v7.0 - Anexo I – Leiaute e Regras de Validação da NF-e e da NFC-e",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=J%20I%20v4eN00E=",
				ContentType: "application/pdf",
			},
		},
	},
}
