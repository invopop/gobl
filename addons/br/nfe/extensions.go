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
	ExtKeyICMSCST         = "br-nfe-icms-cst"
	ExtKeyICMSCSOSN       = "br-nfe-icms-csosn"
	ExtKeyICMSOrigin      = "br-nfe-icms-origin"
	ExtKeyPISCST          = "br-nfe-pis-cst"
	ExtKeyCOFINSCST       = "br-nfe-cofins-cst"
	ExtKeyPurpose         = "br-nfe-purpose"
	ExtKeyOperationType   = "br-nfe-operation-type"
	ExtKeyCreditNoteType  = "br-nfe-credit-note-type"
	ExtKeyDebitNoteType   = "br-nfe-debit-note-type"
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

// Purpose Codes
const (
	PurposeNormal        cbc.Code = "1"
	PurposeComplementary cbc.Code = "2"
	PurposeAdjustment    cbc.Code = "3"
	PurposeGoodsReturn   cbc.Code = "4"
	PurposeCreditNote    cbc.Code = "5"
	PurposeDebitNote     cbc.Code = "6"
)

// Operation Type Codes
const (
	OperationInbound  cbc.Code = "0"
	OperationOutbound cbc.Code = "1"
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
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "MEI - Individual Micro-entrepreneur",
					i18n.PT: "MEI - Microempreendedor Individual",
				},
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates the tax regime that a party is subject to. Defaults to ~3~
				(normal regime) during normalization when not provided.
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
			{
				Title: i18n.String{
					i18n.EN: "NF-e/NFC-e Technical Note 2024.001 - CRT=4 MEI",
					i18n.PT: "Nota Técnica NF-e/NFC-e 2024.001 - CRT=4 MEI",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=kIiniiSkpKc=",
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
	{
		Key: ExtKeyICMSCST,
		Name: i18n.String{
			i18n.EN: "ICMS Tax Status Code (CST)",
			i18n.PT: "Código de Situação Tributária do ICMS (CST)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				ICMS tax status code (CST) for the line item, used by issuers under the normal
				regime (~br-nfe-regime~ ~3~). Simples Nacional issuers must use the CSOSN code
				(~br-nfe-icms-csosn~) instead. Defaults to ~00~ (taxed in full) during
				normalization when not provided.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "00",
				Name: i18n.String{
					i18n.EN: "Taxed in full",
					i18n.PT: "Tributada integralmente",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Monophasic taxation (fuels)",
					i18n.PT: "Tributação monofásica própria sobre combustíveis",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Taxed with ICMS charged by tax substitution",
					i18n.PT: "Tributada e com cobrança do ICMS por substituição tributária",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Fuel operation with ICMS retention",
					i18n.PT: "Operação com combustível e retenção do ICMS",
				},
			},
			{
				Code: "20",
				Name: i18n.String{
					i18n.EN: "Taxed with tax base reduction",
					i18n.PT: "Com redução da base de cálculo",
				},
			},
			{
				Code: "30",
				Name: i18n.String{
					i18n.EN: "Exempt/non-taxed with ICMS charged by tax substitution",
					i18n.PT: "Isenta/não tributada e com cobrança do ICMS por substituição tributária",
				},
			},
			{
				Code: "40",
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PT: "Isenta",
				},
			},
			{
				Code: "41",
				Name: i18n.String{
					i18n.EN: "Non-taxed",
					i18n.PT: "Não tributada",
				},
			},
			{
				Code: "50",
				Name: i18n.String{
					i18n.EN: "Suspended",
					i18n.PT: "Com suspensão",
				},
			},
			{
				Code: "51",
				Name: i18n.String{
					i18n.EN: "Deferred",
					i18n.PT: "Com diferimento",
				},
			},
			{
				Code: "53",
				Name: i18n.String{
					i18n.EN: "Fuel operation with deferral",
					i18n.PT: "Operação com combustível e diferimento",
				},
			},
			{
				Code: "60",
				Name: i18n.String{
					i18n.EN: "ICMS charged previously by tax substitution",
					i18n.PT: "ICMS cobrado anteriormente por substituição tributária",
				},
			},
			{
				Code: "61",
				Name: i18n.String{
					i18n.EN: "Fuel operation with ICMS withheld",
					i18n.PT: "Operação com combustível e ICMS retido anteriormente",
				},
			},
			{
				Code: "70",
				Name: i18n.String{
					i18n.EN: "Taxed with base reduction and ICMS charged by tax substitution",
					i18n.PT: "Com redução da base de cálculo e cobrança do ICMS por substituição tributária",
				},
			},
			{
				Code: "90",
				Name: i18n.String{
					i18n.EN: "Others",
					i18n.PT: "Outras",
				},
			},
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
		Key: ExtKeyICMSCSOSN,
		Name: i18n.String{
			i18n.EN: "ICMS Simples Nacional Status Code (CSOSN)",
			i18n.PT: "Código de Situação da Operação no Simples Nacional (CSOSN)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				ICMS status code (CSOSN) for the line item, used by issuers under the Simples
				Nacional regime (~br-nfe-regime~ ~1~, ~2~ or ~4~). Normal-regime issuers must
				use the CST code (~br-nfe-icms-cst~) instead.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "101",
				Name: i18n.String{
					i18n.EN: "Taxed under Simples Nacional with credit permission",
					i18n.PT: "Tributada pelo Simples Nacional com permissão de crédito",
				},
			},
			{
				Code: "102",
				Name: i18n.String{
					i18n.EN: "Taxed under Simples Nacional without credit permission",
					i18n.PT: "Tributada pelo Simples Nacional sem permissão de crédito",
				},
			},
			{
				Code: "103",
				Name: i18n.String{
					i18n.EN: "ICMS exemption under Simples Nacional for gross revenue range",
					i18n.PT: "Isenção do ICMS no Simples Nacional para faixa de receita bruta",
				},
			},
			{
				Code: "201",
				Name: i18n.String{
					i18n.EN: "Taxed under Simples Nacional with credit permission and ICMS by tax substitution",
					i18n.PT: "Tributada pelo Simples Nacional com permissão de crédito e com cobrança do ICMS por substituição tributária",
				},
			},
			{
				Code: "202",
				Name: i18n.String{
					i18n.EN: "Taxed under Simples Nacional without credit permission and ICMS by tax substitution",
					i18n.PT: "Tributada pelo Simples Nacional sem permissão de crédito e com cobrança do ICMS por substituição tributária",
				},
			},
			{
				Code: "203",
				Name: i18n.String{
					i18n.EN: "ICMS exemption under Simples Nacional for revenue range and ICMS by tax substitution",
					i18n.PT: "Isenção do ICMS no Simples Nacional para faixa de receita bruta e com cobrança do ICMS por substituição tributária",
				},
			},
			{
				Code: "300",
				Name: i18n.String{
					i18n.EN: "Immune",
					i18n.PT: "Imune",
				},
			},
			{
				Code: "400",
				Name: i18n.String{
					i18n.EN: "Non-taxed under Simples Nacional",
					i18n.PT: "Não tributada pelo Simples Nacional",
				},
			},
			{
				Code: "500",
				Name: i18n.String{
					i18n.EN: "ICMS charged previously by tax substitution or anticipation",
					i18n.PT: "ICMS cobrado anteriormente por substituição tributária ou por antecipação",
				},
			},
			{
				Code: "900",
				Name: i18n.String{
					i18n.EN: "Others",
					i18n.PT: "Outros",
				},
			},
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
		Key: ExtKeyICMSOrigin,
		Name: i18n.String{
			i18n.EN: "ICMS Goods Origin",
			i18n.PT: "Origem da Mercadoria (ICMS)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Origin of the goods for ICMS purposes. Defaults to ~0~ (national) during
				normalization when not provided.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "0",
				Name: i18n.String{
					i18n.EN: "National (except codes 3, 4, 5 and 8)",
					i18n.PT: "Nacional, exceto as indicadas nos códigos 3, 4, 5 e 8",
				},
			},
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Foreign, direct import (except code 6)",
					i18n.PT: "Estrangeira - Importação direta, exceto a indicada no código 6",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Foreign, acquired in the domestic market (except code 7)",
					i18n.PT: "Estrangeira - Adquirida no mercado interno, exceto a indicada no código 7",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "National, import content above 40% and up to 70%",
					i18n.PT: "Nacional, mercadoria ou bem com Conteúdo de Importação superior a 40% e inferior ou igual a 70%",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "National, produced under the basic productive processes",
					i18n.PT: "Nacional, cuja produção tenha sido feita em conformidade com os processos produtivos básicos de que tratam as legislações citadas nos Ajustes",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "National, import content up to 40%",
					i18n.PT: "Nacional, mercadoria ou bem com Conteúdo de Importação inferior ou igual a 40%",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Foreign, direct import, no domestic equivalent (CAMEX list / natural gas)",
					i18n.PT: "Estrangeira - Importação direta, sem similar nacional, constante em lista da CAMEX e gás natural",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Foreign, acquired domestically, no domestic equivalent (CAMEX list / natural gas)",
					i18n.PT: "Estrangeira - Adquirida no mercado interno, sem similar nacional, constante lista CAMEX e gás natural",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "National, import content above 70%",
					i18n.PT: "Nacional, mercadoria ou bem com Conteúdo de Importação superior a 70%",
				},
			},
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
		Key: ExtKeyPISCST,
		Name: i18n.String{
			i18n.EN: "PIS Tax Status Code (CST)",
			i18n.PT: "Código de Situação Tributária do PIS (CST)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				PIS tax status code for the line item. Simples Nacional issuers typically use
				~49~ or ~99~, as PIS is settled within the unified DAS collection. Defaults to
				~01~ (standard-rate taxable operation) during normalization when not provided.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "Taxable operation, standard rate",
					i18n.PT: "Operação Tributável com Alíquota Básica",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Taxable operation, differentiated rate",
					i18n.PT: "Operação Tributável com Alíquota Diferenciada",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "Taxable operation, rate per unit of measure",
					i18n.PT: "Operação Tributável com Alíquota por Unidade de Medida de Produto",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "Taxable monophasic operation, zero-rate resale",
					i18n.PT: "Operação Tributável Monofásica - Revenda a Alíquota Zero",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "Taxable operation by tax substitution",
					i18n.PT: "Operação Tributável por Substituição Tributária",
				},
			},
			{
				Code: "06",
				Name: i18n.String{
					i18n.EN: "Taxable operation, zero rate",
					i18n.PT: "Operação Tributável a Alíquota Zero",
				},
			},
			{
				Code: "07",
				Name: i18n.String{
					i18n.EN: "Operation exempt from the contribution",
					i18n.PT: "Operação Isenta da Contribuição",
				},
			},
			{
				Code: "08",
				Name: i18n.String{
					i18n.EN: "Operation without incidence of the contribution",
					i18n.PT: "Operação Sem Incidência da Contribuição",
				},
			},
			{
				Code: "09",
				Name: i18n.String{
					i18n.EN: "Operation with suspension of the contribution",
					i18n.PT: "Operação com Suspensão da Contribuição",
				},
			},
			{
				Code: "49",
				Name: i18n.String{
					i18n.EN: "Other outbound operations",
					i18n.PT: "Outras Operações de Saída",
				},
			},
			{
				Code: "50",
				Name: i18n.String{
					i18n.EN: "Credit operation, exclusively taxed domestic revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Tributada no Mercado Interno",
				},
			},
			{
				Code: "51",
				Name: i18n.String{
					i18n.EN: "Credit operation, exclusively non-taxed domestic revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Não Tributada no Mercado Interno",
				},
			},
			{
				Code: "52",
				Name: i18n.String{
					i18n.EN: "Credit operation, exclusively export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada Exclusivamente a Receita de Exportação",
				},
			},
			{
				Code: "53",
				Name: i18n.String{
					i18n.EN: "Credit operation, taxed and non-taxed domestic revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno",
				},
			},
			{
				Code: "54",
				Name: i18n.String{
					i18n.EN: "Credit operation, taxed domestic and export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "55",
				Name: i18n.String{
					i18n.EN: "Credit operation, non-taxed domestic and export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "56",
				Name: i18n.String{
					i18n.EN: "Credit operation, taxed/non-taxed domestic and export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "60",
				Name: i18n.String{
					i18n.EN: "Presumed credit, exclusively taxed domestic revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Tributada no Mercado Interno",
				},
			},
			{
				Code: "61",
				Name: i18n.String{
					i18n.EN: "Presumed credit, exclusively non-taxed domestic revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Não-Tributada no Mercado Interno",
				},
			},
			{
				Code: "62",
				Name: i18n.String{
					i18n.EN: "Presumed credit, exclusively export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita de Exportação",
				},
			},
			{
				Code: "63",
				Name: i18n.String{
					i18n.EN: "Presumed credit, taxed and non-taxed domestic revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno",
				},
			},
			{
				Code: "64",
				Name: i18n.String{
					i18n.EN: "Presumed credit, taxed domestic and export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "65",
				Name: i18n.String{
					i18n.EN: "Presumed credit, non-taxed domestic and export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "66",
				Name: i18n.String{
					i18n.EN: "Presumed credit, taxed/non-taxed domestic and export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "67",
				Name: i18n.String{
					i18n.EN: "Presumed credit, other operations",
					i18n.PT: "Crédito Presumido - Outras Operações",
				},
			},
			{
				Code: "70",
				Name: i18n.String{
					i18n.EN: "Acquisition operation without credit right",
					i18n.PT: "Operação de Aquisição sem Direito a Crédito",
				},
			},
			{
				Code: "71",
				Name: i18n.String{
					i18n.EN: "Acquisition operation with exemption",
					i18n.PT: "Operação de Aquisição com Isenção",
				},
			},
			{
				Code: "72",
				Name: i18n.String{
					i18n.EN: "Acquisition operation with suspension",
					i18n.PT: "Operação de Aquisição com Suspensão",
				},
			},
			{
				Code: "73",
				Name: i18n.String{
					i18n.EN: "Acquisition operation at zero rate",
					i18n.PT: "Operação de Aquisição a Alíquota Zero",
				},
			},
			{
				Code: "74",
				Name: i18n.String{
					i18n.EN: "Acquisition operation without incidence of the contribution",
					i18n.PT: "Operação de Aquisição sem Incidência da Contribuição",
				},
			},
			{
				Code: "75",
				Name: i18n.String{
					i18n.EN: "Acquisition operation by tax substitution",
					i18n.PT: "Operação de Aquisição por Substituição Tributária",
				},
			},
			{
				Code: "98",
				Name: i18n.String{
					i18n.EN: "Other inbound operations",
					i18n.PT: "Outras Operações de Entrada",
				},
			},
			{
				Code: "99",
				Name: i18n.String{
					i18n.EN: "Other operations",
					i18n.PT: "Outras Operações",
				},
			},
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
		Key: ExtKeyCOFINSCST,
		Name: i18n.String{
			i18n.EN: "COFINS Tax Status Code (CST)",
			i18n.PT: "Código de Situação Tributária do COFINS (CST)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				COFINS tax status code for the line item. Simples Nacional issuers typically
				use ~49~ or ~99~, as COFINS is settled within the unified DAS collection.
				Defaults to ~01~ (standard-rate taxable operation) during normalization when
				not provided.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "Taxable operation, standard rate",
					i18n.PT: "Operação Tributável com Alíquota Básica",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Taxable operation, differentiated rate",
					i18n.PT: "Operação Tributável com Alíquota Diferenciada",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "Taxable operation, rate per unit of measure",
					i18n.PT: "Operação Tributável com Alíquota por Unidade de Medida de Produto",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "Taxable monophasic operation, zero-rate resale",
					i18n.PT: "Operação Tributável Monofásica - Revenda a Alíquota Zero",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "Taxable operation by tax substitution",
					i18n.PT: "Operação Tributável por Substituição Tributária",
				},
			},
			{
				Code: "06",
				Name: i18n.String{
					i18n.EN: "Taxable operation, zero rate",
					i18n.PT: "Operação Tributável a Alíquota Zero",
				},
			},
			{
				Code: "07",
				Name: i18n.String{
					i18n.EN: "Operation exempt from the contribution",
					i18n.PT: "Operação Isenta da Contribuição",
				},
			},
			{
				Code: "08",
				Name: i18n.String{
					i18n.EN: "Operation without incidence of the contribution",
					i18n.PT: "Operação Sem Incidência da Contribuição",
				},
			},
			{
				Code: "09",
				Name: i18n.String{
					i18n.EN: "Operation with suspension of the contribution",
					i18n.PT: "Operação com Suspensão da Contribuição",
				},
			},
			{
				Code: "49",
				Name: i18n.String{
					i18n.EN: "Other outbound operations",
					i18n.PT: "Outras Operações de Saída",
				},
			},
			{
				Code: "50",
				Name: i18n.String{
					i18n.EN: "Credit operation, exclusively taxed domestic revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Tributada no Mercado Interno",
				},
			},
			{
				Code: "51",
				Name: i18n.String{
					i18n.EN: "Credit operation, exclusively non-taxed domestic revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Não Tributada no Mercado Interno",
				},
			},
			{
				Code: "52",
				Name: i18n.String{
					i18n.EN: "Credit operation, exclusively export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada Exclusivamente a Receita de Exportação",
				},
			},
			{
				Code: "53",
				Name: i18n.String{
					i18n.EN: "Credit operation, taxed and non-taxed domestic revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno",
				},
			},
			{
				Code: "54",
				Name: i18n.String{
					i18n.EN: "Credit operation, taxed domestic and export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "55",
				Name: i18n.String{
					i18n.EN: "Credit operation, non-taxed domestic and export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "56",
				Name: i18n.String{
					i18n.EN: "Credit operation, taxed/non-taxed domestic and export revenue",
					i18n.PT: "Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "60",
				Name: i18n.String{
					i18n.EN: "Presumed credit, exclusively taxed domestic revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Tributada no Mercado Interno",
				},
			},
			{
				Code: "61",
				Name: i18n.String{
					i18n.EN: "Presumed credit, exclusively non-taxed domestic revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Não-Tributada no Mercado Interno",
				},
			},
			{
				Code: "62",
				Name: i18n.String{
					i18n.EN: "Presumed credit, exclusively export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita de Exportação",
				},
			},
			{
				Code: "63",
				Name: i18n.String{
					i18n.EN: "Presumed credit, taxed and non-taxed domestic revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno",
				},
			},
			{
				Code: "64",
				Name: i18n.String{
					i18n.EN: "Presumed credit, taxed domestic and export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "65",
				Name: i18n.String{
					i18n.EN: "Presumed credit, non-taxed domestic and export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "66",
				Name: i18n.String{
					i18n.EN: "Presumed credit, taxed/non-taxed domestic and export revenue",
					i18n.PT: "Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno e de Exportação",
				},
			},
			{
				Code: "67",
				Name: i18n.String{
					i18n.EN: "Presumed credit, other operations",
					i18n.PT: "Crédito Presumido - Outras Operações",
				},
			},
			{
				Code: "70",
				Name: i18n.String{
					i18n.EN: "Acquisition operation without credit right",
					i18n.PT: "Operação de Aquisição sem Direito a Crédito",
				},
			},
			{
				Code: "71",
				Name: i18n.String{
					i18n.EN: "Acquisition operation with exemption",
					i18n.PT: "Operação de Aquisição com Isenção",
				},
			},
			{
				Code: "72",
				Name: i18n.String{
					i18n.EN: "Acquisition operation with suspension",
					i18n.PT: "Operação de Aquisição com Suspensão",
				},
			},
			{
				Code: "73",
				Name: i18n.String{
					i18n.EN: "Acquisition operation at zero rate",
					i18n.PT: "Operação de Aquisição a Alíquota Zero",
				},
			},
			{
				Code: "74",
				Name: i18n.String{
					i18n.EN: "Acquisition operation without incidence of the contribution",
					i18n.PT: "Operação de Aquisição sem Incidência da Contribuição",
				},
			},
			{
				Code: "75",
				Name: i18n.String{
					i18n.EN: "Acquisition operation by tax substitution",
					i18n.PT: "Operação de Aquisição por Substituição Tributária",
				},
			},
			{
				Code: "98",
				Name: i18n.String{
					i18n.EN: "Other inbound operations",
					i18n.PT: "Outras Operações de Entrada",
				},
			},
			{
				Code: "99",
				Name: i18n.String{
					i18n.EN: "Other operations",
					i18n.PT: "Outras Operações",
				},
			},
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
		Key: ExtKeyPurpose,
		Name: i18n.String{
			i18n.EN: "Purpose Code",
			i18n.PT: "Código de Finalidade",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code that identifies the purpose of the fiscal document (CEFAZ field ~finNFe~,
				B25). Standard invoices are set to ~1~ (normal) via a tax scenario.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: PurposeNormal,
				Name: i18n.String{
					i18n.EN: "Normal",
					i18n.PT: "Normal",
				},
			},
			{
				Code: PurposeComplementary,
				Name: i18n.String{
					i18n.EN: "Complementary",
					i18n.PT: "Complementar",
				},
			},
			{
				Code: PurposeAdjustment,
				Name: i18n.String{
					i18n.EN: "Adjustment",
					i18n.PT: "Ajuste",
				},
			},
			{
				Code: PurposeGoodsReturn,
				Name: i18n.String{
					i18n.EN: "Goods Return",
					i18n.PT: "Devolução/Retorno de mercadoria",
				},
			},
			{
				Code: PurposeCreditNote,
				Name: i18n.String{
					i18n.EN: "Credit Note",
					i18n.PT: "Nota de Crédito",
				},
			},
			{
				Code: PurposeDebitNote,
				Name: i18n.String{
					i18n.EN: "Debit Note",
					i18n.PT: "Nota de Débito",
				},
			},
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
			{
				Title: i18n.String{
					i18n.EN: "NF-e/NFC-e Technical Note 2025.002 - RTC",
					i18n.PT: "Nota Técnica NF-e/NFC-e 2025.002 - RTC",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=pD4YrecPV6s=",
				ContentType: "application/pdf",
			},
		},
	},
	{
		Key: ExtKeyOperationType,
		Name: i18n.String{
			i18n.EN: "Operation Type Code",
			i18n.PT: "Código do Tipo de Operação",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code that identifies the type of operation, indicating whether it is inbound
				or outbound (CEFAZ field ~tpNF~, B11). Standard invoices are set to ~1~
				(outbound) via a tax scenario.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: OperationInbound,
				Name: i18n.String{
					i18n.EN: "Inbound",
					i18n.PT: "Entrada",
				},
			},
			{
				Code: OperationOutbound,
				Name: i18n.String{
					i18n.EN: "Outbound",
					i18n.PT: "Saída",
				},
			},
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
		Key: ExtKeyCreditNoteType,
		Name: i18n.String{
			i18n.EN: "Credit Note Type Code",
			i18n.PT: "Código do Tipo de Nota de Crédito",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code that identifies the type of credit note (CEFAZ field ~tpNFCredito~,
				B25.2) according to the RTM (applies to IBS/CBS only).
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "Penalty and interest",
					i18n.PT: "Multa e juros",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Appropriation of presumed IBS credit in the ZFM",
					i18n.PT: "Apropriação de crédito presumido de IBS na ZFM",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "Return",
					i18n.PT: "Retorno",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "Reduction of values",
					i18n.PT: "Redução de valores",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "Credit transfer on succession",
					i18n.PT: "Transferência de crédito na sucessão",
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "NF-e/NFC-e Technical Note 2025.002 - RTC",
					i18n.PT: "Nota Técnica NF-e/NFC-e 2025.002 - RTC",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=pD4YrecPV6s=",
				ContentType: "application/pdf",
			},
		},
	},
	{
		Key: ExtKeyDebitNoteType,
		Name: i18n.String{
			i18n.EN: "Debit Note Type Code",
			i18n.PT: "Código do Tipo de Nota de Débito",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code that identifies the type of debit note (CEFAZ field ~tpNFDebito~, B25.1)
				according to the RTM (applies to IBS/CBS only).
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "Transfer of credits from Cooperatives",
					i18n.PT: "Transferência de créditos de Cooperativas",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Credit annulment",
					i18n.PT: "Anulação de crédito",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "Debits not processed in the regular assessment",
					i18n.PT: "Débitos não processados na apuração regular",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "Penalty and interest",
					i18n.PT: "Multa e juros",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "Credit transfer on succession",
					i18n.PT: "Transferência de crédito na sucessão",
				},
			},
			{
				Code: "06",
				Name: i18n.String{
					i18n.EN: "Advance payment",
					i18n.PT: "Pagamento antecipado",
				},
			},
			{
				Code: "07",
				Name: i18n.String{
					i18n.EN: "Inventory loss",
					i18n.PT: "Perda em estoque",
				},
			},
			{
				Code: "08",
				Name: i18n.String{
					i18n.EN: "Exit from Simples Nacional",
					i18n.PT: "Desenquadramento do Simples Nacional",
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "NF-e/NFC-e Technical Note 2025.002 - RTC",
					i18n.PT: "Nota Técnica NF-e/NFC-e 2025.002 - RTC",
				},
				URL:         "https://www.nfe.fazenda.gov.br/portal/exibirArquivo.aspx?conteudo=pD4YrecPV6s=",
				ContentType: "application/pdf",
			},
		},
	},
}
