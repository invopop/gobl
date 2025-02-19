package saft

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// SAF-T Extension Keys
const (
	ExtKeyExemption    cbc.Key = "pt-saft-exemption"
	ExtKeyTaxRate      cbc.Key = "pt-saft-tax-rate"
	ExtKeyInvoiceType  cbc.Key = "pt-saft-invoice-type"
	ExtKeyProductType  cbc.Key = "pt-saft-product-type"
	ExtKeyPaymentType  cbc.Key = "pt-saft-payment-type"
	ExtKeyPaymentMeans cbc.Key = "pt-saft-payment-means"
)

// Invoice types
const (
	InvoiceTypeStandard       cbc.Code = "FT"
	InvoiceTypeSimplified     cbc.Code = "FS"
	InvoiceTypeInvoiceReceipt cbc.Code = "FR"
	InvoiceTypeDebitNote      cbc.Code = "ND"
	InvoiceTypeCreditNote     cbc.Code = "NC"
)

// Payment types
const (
	PaymentTypeCash  cbc.Code = "RC"
	PaymentTypeOther cbc.Code = "RG"
)

// Tax rates
const (
	TaxRateReduced      cbc.Code = "RED"
	TaxRateIntermediate cbc.Code = "INT"
	TaxRateNormal       cbc.Code = "NOR"
	TaxRateExempt       cbc.Code = "ISE"
	TaxRateOther        cbc.Code = "OUT"
)

// Product Types
const (
	ProductTypeGoods   cbc.Code = "P"
	ProductTypeService cbc.Code = "S"
	ProductTypeOther   cbc.Code = "O"
	ProductTypeExcise  cbc.Code = "E"
	ProductTypeFee     cbc.Code = "I"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyInvoiceType,
		Name: i18n.String{
			i18n.EN: "Invoice Type",
			i18n.PT: "Tipo de Fatura",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				SAF-T's ~InvoiceType~ (Tipo de documento) specifies the type of a sales invoice. In GOBL,
				this type can be set using the ~pt-saft-invoice-type~ extension in the tax section. GOBL
				will set the extension for you based on the type and the tax tags you set in your invoice.

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$tags": [
						"invoice-receipt"
					],
					// ...
					"type": "standard",
					// ...
					"tax": {
						"ext": {
							"pt-saft-invoice-type": "FR"
						}
					},
					// ...
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: InvoiceTypeStandard,
				Name: i18n.String{
					i18n.EN: "Standard Invoice",
					i18n.PT: "Fatura",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued under article 36 of the VAT code.",
					i18n.PT: "Fatura, emitida nos termos do artigo 36.o do Código do IVA",
				},
			},
			{
				Code: InvoiceTypeSimplified,
				Name: i18n.String{
					i18n.EN: "Simplified Invoice",
					i18n.PT: "Fatura Simplificada",
				},
				Desc: i18n.String{
					i18n.EN: "Simplified invoice issued under article 40 of the VAT code.",
					i18n.PT: "Fatura simplificada, emitida nos termos do artigo 40.o do Código do IVA",
				},
			},
			{
				Code: InvoiceTypeInvoiceReceipt,
				Name: i18n.String{
					i18n.EN: "Invoice-Receipt",
					i18n.PT: "Fatura-Recibo",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued after payment.",
					i18n.PT: "Fatura-recibo",
				},
			},
			{
				Code: InvoiceTypeDebitNote,
				Name: i18n.String{
					i18n.EN: "Debit Note",
					i18n.PT: "Nota de Débito",
				},
			},
			{
				Code: InvoiceTypeCreditNote,
				Name: i18n.String{
					i18n.EN: "Credit Note",
					i18n.PT: "Nota de Crédito",
				},
			},
		},
	},
	{
		Key: ExtKeyPaymentType,
		Name: i18n.String{
			i18n.EN: "Payment Type",
			i18n.PT: "Tipo de Pagamento",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				To report payment receipts to the AT, GOBL provides conversion from ~bill.Payment~
				documents. In a payment, the SAF-T's ~PaymentType~ (Tipo de documento) field specifies its
				type. In GOBL, this type can be set using the ~pt-saft-payment-type~ extension. GOBL will
				set the extension automatically based on the type and the tax tags you set. The table
				below shows how this mapping is done:

				| Code | Name                                       | GOBL Type | GOBL Tax Tag |
				| ---- | ------------------------------------------ | --------- | ------------ |
				| RG   | Outro Recibo                               | ~receipt~ |              |
				| RC   | Recibo no âmbito do regime de IVA de Caixa | ~receipt~ | ~vat-cash~   |

				For example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/payment",
					// ...
					"type": "receipt",
					// ...
					"ext": {
						"pt-saft-receipt-type": "RG"
					},
					// ...
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: PaymentTypeCash,
				Name: i18n.String{
					i18n.EN: "Receipt under the VAT Cash scheme",
					i18n.PT: "Recibo no âmbito do regime de IVA de Caixa",
				},
			},
			{
				Code: PaymentTypeOther,
				Name: i18n.String{
					i18n.EN: "Other Receipt",
					i18n.PT: "Outro Recibo",
				},
			},
		},
	},
	{
		Key: ExtKeyTaxRate,
		Name: i18n.String{
			i18n.EN: "Tax Rate Code",
			i18n.PT: "Código da Taxa de Imposto",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The SAF-T's ~TaxCode~ (Código do imposto) is required for invoice items that apply VAT.
				GOBL provides the ~pt-saft-tax-rate~ extension to set this code at line tax level. It also
				determines it automatically this code using the ~rate~ field (when present). The following
				table lists the supported tax codes and how GOBL will map them:

				| Code   | Name            | GOBL Tax Rate  |
				| ------ | --------------- | -------------- |
				| ~NOR~  | Tipo Geral      | ~standard~     |
				| ~INT~  | Taxa Intermédia | ~intermediate~ |
				| ~RED~  | Taxa Reduzida   | ~reduced~      |
				| ~ISE~  | Isenta          | ~exempt~       |
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced",
					i18n.PT: "Redução",
				},
			},
			{
				Code: TaxRateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate",
					i18n.PT: "Intermédio",
				},
			},
			{
				Code: TaxRateNormal,
				Name: i18n.String{
					i18n.EN: "Normal",
					i18n.PT: "Normal",
				},
			},
			{
				Code: TaxRateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PT: "Isento",
				},
			},
			{
				Code: TaxRateOther,
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.PT: "Outro",
				},
			},
		},
	},
	{
		Key: ExtKeyExemption,
		Name: i18n.String{
			i18n.EN: "Tax exemption reason code",
			i18n.PT: "Código do motivo de isenção de imposto",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				AT's ~TaxExemptionCode~ (Código do motivo de isenção de imposto) is a code that
				specifies the reason the VAT tax is exempt in a Portuguese invoice. When the ~exempt~ tag
				is used in a tax combo, the ~ext~ map's ~pt-exemption-code~ property is required.

				For example, you could define an invoice line exempt of tax as follows:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"lines": [
						{
							// ...
							"item": {
								"name": "Some service exempt of tax",
								"price": "25.00"
							},
							"tax": [
								{
										"cat": "VAT",
										"rate": "exempt",
										"ext": {
											"pt-saft-tax-rate": "ISE",
											"pt-saft-exemption": "M19"
										}
								}
							]
						}
					]
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "M01",
				Name: i18n.String{
					i18n.EN: "Article 16, No. 6 of the VAT code",
					i18n.PT: "Artigo 16.°, n.° 6 do CIVA",
				},
			},
			{
				Code: "M02",
				Name: i18n.String{
					i18n.EN: "Article 6 of the Decree-Law 198/90 of 19th June",
					i18n.PT: "Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho",
				},
			},
			{
				Code: "M04",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 13 of the VAT code",
					i18n.PT: "Isento artigo 13.° do CIVA",
				},
			},
			{
				Code: "M05",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 14 of the VAT code",
					i18n.PT: "Isento artigo 14.° do CIVA",
				},
			},
			{
				Code: "M06",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 15 of the VAT code",
					i18n.PT: "Isento artigo 15.° do CIVA",
				},
			},
			{
				Code: "M07",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 9 of the VAT code",
					i18n.PT: "Isento artigo 9.° do CIVA",
				},
			},
			{
				Code: "M09",
				Name: i18n.String{
					i18n.EN: "VAT - does not confer right to deduct / Article 62 paragraph b) of the VAT code",
					i18n.PT: "IVA - não confere direito a dedução / Artigo 62.° alínea b) do CIVA",
				},
			},
			{
				Code: "M10",
				Name: i18n.String{
					i18n.EN: "VAT - exemption scheme / Article 57 of the VAT code",
					i18n.PT: "IVA - regime de isenção / Artigo 57.° do CIVA",
				},
			},
			{
				Code: "M11",
				Name: i18n.String{
					i18n.EN: "Special scheme for tobacco / Decree-Law No. 346/85 of 23rd August",
					i18n.PT: "Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto",
				},
			},
			{
				Code: "M12",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Travel agencies / Decree-Law No. 221/85 of 3rd July",
					i18n.PT: "Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho",
				},
			},
			{
				Code: "M13",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Second-hand goods / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
			},
			{
				Code: "M14",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Works of art / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
			},
			{
				Code: "M15",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Collector’s items and antiques / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
			},
			{
				Code: "M16",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 14 of the RITI",
					i18n.PT: "Isento artigo 14.° do RITI",
				},
			},
			{
				Code: "M19",
				Name: i18n.String{
					i18n.EN: "Other exemptions - Temporary exemptions determined by specific legislation",
					i18n.PT: "Outras isenções - Isenções temporárias determinadas em diploma próprio",
				},
			},
			{
				Code: "M20",
				Name: i18n.String{
					i18n.EN: "VAT - flat-rate scheme / Article 59-D No. 2 of the VAT code",
					i18n.PT: "IVA - regime forfetário / Artigo 59.°-D n.°2 do CIVA",
				},
			},
			{
				Code: "M21",
				Name: i18n.String{
					i18n.EN: "VAT - does not confer right to deduct (or similar) - Article 72 No. 4 of the VAT code",
					i18n.PT: "IVA - não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA",
				},
			},
			{
				Code: "M25",
				Name: i18n.String{
					i18n.EN: "Consignment goods - Article 38 No. 1 paragraph a) of the VAT code",
					i18n.PT: "Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA",
				},
			},
			{
				Code: "M30",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph i) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea i) do CIVA",
				},
			},
			{
				Code: "M31",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph j) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea j) do CIVA",
				},
			},
			{
				Code: "M32",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph l) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea I) do CIVA",
				},
			},
			{
				Code: "M33",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph m) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea m) do CIVA",
				},
			},
			{
				Code: "M40",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 6 No. 6 paragraph a) of the VAT code, to the contrary",
					i18n.PT: "IVA - autoliquidação / Artigo 6.° n.° 6 alínea a) do CIVA, a contrário",
				},
			},
			{
				Code: "M41",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 8 No. 3 of the RITI",
					i18n.PT: "IVA - autoliquidação / Artigo 8.° n.° 3 do RITI",
				},
			},
			{
				Code: "M42",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Decree-Law No. 21/2007 of 29 January",
					i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 21/2007, de 29 de janeiro",
				},
			},
			{
				Code: "M43",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Decree-Law No. 362/99 of 16th September",
					i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 362/99, de 16 de setembro",
				},
			},
			{
				Code: "M99",
				Name: i18n.String{
					i18n.EN: "Not subject to tax or not taxed",
					i18n.PT: "Não sujeito ou não tributado",
				},
			},
		},
	},
	{
		Key: ExtKeyProductType,
		Name: i18n.String{
			i18n.EN: "Product Type",
			i18n.PT: "Tipo de Produto",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				SAF-T's ~ProductType~ (Indicador de produto ou serviço) indicates the type of each line
				item in an invoice. The ~pt-saft-product-type~ extension used at line item level allows to
				set the product type to one of the allowed values.

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"lines": [
						{
							// ...
							"item": {
								"name": "Some service",
								"price": "25.00",
								"ext": {
									"pt-saft-product-type": "S"
								}
							},
							// ...
						}
					]
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ProductTypeGoods,
				Name: i18n.String{
					i18n.EN: "Goods",
					i18n.PT: "Produtos",
				},
			},
			{
				Code: ProductTypeService,
				Name: i18n.String{
					i18n.EN: "Services",
					i18n.PT: "Serviços",
				},
			},
			{
				Code: ProductTypeOther,
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.PT: "Outros",
				},
				Desc: i18n.String{
					i18n.EN: "Other (e.g., debited postage, advances received or disposal of assets)",
					i18n.PT: "Outros (ex., portes debitados, adiantamentos recebidos ou alienação de ativos)",
				},
			},
			{
				Code: ProductTypeExcise,
				Name: i18n.String{
					i18n.EN: "Excise Duties",
					i18n.PT: "Impostos Especiais",
				},
				Desc: i18n.String{
					i18n.EN: "Excise Duties (e.g., IABA, ISP, IT)",
					i18n.PT: "Impostos Especiais de Consumo (ex., IABA, ISP, IT)",
				},
			},
			{
				Code: ProductTypeFee,
				Name: i18n.String{
					i18n.EN: "Taxes/Fees",
					i18n.PT: "Impostos/Taxas",
				},
				Desc: i18n.String{
					i18n.EN: "Taxes, fees and parafiscal charges (except VAT and IS which should be reflected in table 2.5 - TaxTable and Excise Duties, which should be filled in with code 'E')",
					i18n.PT: "Impostos, taxas e encargos parafiscais – exceto IVA e IS que deverão ser refletidos na tabela 2.5 – Tabela de impostos (TaxTable) e Impostos Especiais de Consumo, que deverão ser preenchidos com o código 'E'.",
				},
			},
		},
	},
	{
		Key: ExtKeyPaymentMeans,
		Name: i18n.String{
			i18n.EN: "Payment Means",
			i18n.PT: "Meio de Pagamento",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The SAF-T's ~PaymentMechanism~ (Meios de pagamento) field specifies the payment means in a
				sales invoice or payment. GOBL provides the ~pt-saft-payment-means~ extension to set this
				value in your ~bill.Invoice~ advances or in you ~bill.Receipt~ method. GOBL maps certain
				payment mean keys automatically to this extension:

				| Code | Name                                               | GOBL Payment Means                                    |
				| ---- | -------------------------------------------------- | ----------------------------------------------------- |
				| ~CC~ | Cartão crédito                                     | ~card~                                                |
				| ~CD~ | Cartão débito                                      | (*)                                                   |
				| ~CH~ | Cheque bancário                                    | ~cheque~                                              |
				| ~CI~ | Letter of credit                                   | (*)                                                   |
				| ~CO~ | Cheque ou cartão oferta                            | (*)                                                   |
				| ~CS~ | Compensação de saldos em conta corrente            | ~netting~                                             |
				| ~DE~ | Dinheiro eletrónico                                | ~online~                                              |
				| ~LC~ | Letra comercial                                    | ~promissory-note~                                     |
				| ~MB~ | Referências de pagamento para Multibanco           | (*)                                                   |
				| ~NU~ | Numerário                                          | ~cash~                                                |
				| ~OU~ | Outro                                              | ~other~                                               |
				| ~PR~ | Permuta de bens                                    | (*)                                                   |
				| ~TB~ | Transferência bancária ou débito direto autorizado | ~credit-transfer~, ~debit-transfer~ or ~direct-debit~ |
				| ~TR~ | Títulos de compensação extrassalarial              | (*)                                                   |

				(*) For codes not mapped from a GOBL Payment Mean, use ~other~ and explicitly set the
				extension.

				For example, in an GOBL invoice:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"payment": {
						"advances": [
							{
								"date": "2023-01-30",
								"key": "credit-transfer",
								"description": "Adiantamento",
								"amount": "100.00",
								"ext": {
									"pt-saft-payment-means": "TB"
								}
							}
						]
					},
					// ...
				}
				~~~

				For example, in a GOBL receipt:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/receipt",
					// ...
					"method": {
						"key": "other",
						"detail": "Compensação extrassalarial",
						"ext": {
							"pt-saft-payment-means": "TR"
						}
					},
					// ...
				}
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "CC",
				Name: i18n.String{
					i18n.EN: "Credit card",
					i18n.PT: "Cartão crédito",
				},
			},
			{
				Code: "CD",
				Name: i18n.String{
					i18n.EN: "Debit card",
					i18n.PT: "Cartão débito",
				},
			},
			{
				Code: "CH",
				Name: i18n.String{
					i18n.EN: "Bank cheque",
					i18n.PT: "Cheque bancário",
				},
			},
			{
				Code: "CI",
				Name: i18n.String{
					i18n.EN: "International documentary credit",
					i18n.PT: "Letter of credit",
				},
			},
			{
				Code: "CO",
				Name: i18n.String{
					i18n.EN: "Gift cheque or card",
					i18n.PT: "Cheque ou cartão oferta",
				},
			},
			{
				Code: "CS",
				Name: i18n.String{
					i18n.EN: "Settlement of balances in current account",
					i18n.PT: "Compensação de saldos em conta corrente",
				},
			},
			{
				Code: "DE",
				Name: i18n.String{
					i18n.EN: "Electronic money",
					i18n.PT: "Dinheiro eletrónico",
				},
			},
			{
				Code: "LC",
				Name: i18n.String{
					i18n.EN: "Commercial bill",
					i18n.PT: "Letra comercial",
				},
			},
			{
				Code: "MB",
				Name: i18n.String{
					i18n.EN: "Multibanco payment references",
					i18n.PT: "Referências de pagamento para Multibanco",
				},
			},
			{
				Code: "NU",
				Name: i18n.String{
					i18n.EN: "Cash",
					i18n.PT: "Numerário",
				},
			},
			{
				Code: "OU",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.PT: "Outro",
				},
			},
			{
				Code: "PR",
				Name: i18n.String{
					i18n.EN: "Barter",
					i18n.PT: "Permuta de bens",
				},
			},
			{
				Code: "TB",
				Name: i18n.String{
					i18n.EN: "Bank transfer or direct debit",
					i18n.PT: "Transferência bancária ou débito direto autorizado",
				},
			},
			{
				Code: "TR",
				Name: i18n.String{
					i18n.EN: "Supplementary compensation",
					i18n.PT: "Títulos de compensação extrassalarial",
				},
			},
		},
	},
}
