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
	ExtKeyProductType  cbc.Key = "pt-saft-product-type"
	ExtKeyPaymentMeans cbc.Key = "pt-saft-payment-means"
	ExtKeySource       cbc.Key = "pt-saft-source"
	ExtKeySourceRef    cbc.Key = "pt-saft-source-ref"

	// Document types extensions
	ExtKeyInvoiceType  cbc.Key = "pt-saft-invoice-type"
	ExtKeyWorkType     cbc.Key = "pt-saft-work-type"
	ExtKeyMovementType cbc.Key = "pt-saft-movement-type"
	ExtKeyPaymentType  cbc.Key = "pt-saft-payment-type"
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

// Document types
const (
	InvoiceTypeStandard       cbc.Code = "FT"
	InvoiceTypeSimplified     cbc.Code = "FS"
	InvoiceTypeInvoiceReceipt cbc.Code = "FR"
	InvoiceTypeDebitNote      cbc.Code = "ND"
	InvoiceTypeCreditNote     cbc.Code = "NC"

	MovementTypeDeliveryNote cbc.Code = "GR"
	MovementTypeWaybill      cbc.Code = "GT"
	MovementTypeFixedAssets  cbc.Code = "GA"
	MovementTypeConsignment  cbc.Code = "GC"
	MovementTypeReturn       cbc.Code = "GD"

	WorkTypeTableQueries      cbc.Code = "CM"
	WorkTypeConsignmentCredit cbc.Code = "CC"
	WorkTypeConsignmentInv    cbc.Code = "FC"
	WorkTypeWorksheets        cbc.Code = "FO"
	WorkTypePurchaseOrder     cbc.Code = "NE"
	WorkTypeOther             cbc.Code = "OU"
	WorkTypeBudgets           cbc.Code = "OR"
	WorkTypeProforma          cbc.Code = "PF"
	WorkTypeDocuments         cbc.Code = "DC"
	WorkTypePremium           cbc.Code = "RP"
	WorkTypeChargeback        cbc.Code = "RE"
	WorkTypeCoInsurers        cbc.Code = "CS"
	WorkTypeLeadCoInsurer     cbc.Code = "LD"
	WorkTypeReinsurance       cbc.Code = "RA"

	PaymentTypeCash  cbc.Code = "RC"
	PaymentTypeOther cbc.Code = "RG"
)

// SourceBilling values
const (
	SourceBillingProduced   cbc.Code = "P"
	SourceBillingIntegrated cbc.Code = "I"
	SourceBillingManual     cbc.Code = "M"
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

				The table below shows how this mapping is done:

				| Code | Name                | GOBL Type     | GOBL Tax Tag    |
				| ---- | ------------------- | ------------- | --------------- |
				| ~FT~ | Standard Invoice    | ~standard~    |                 |
				| ~FS~ | Simplified Invoice  | ~standard~    | ~simplified~    |
				| ~FR~ | Invoice-Receipt     | ~standard~    | ~invoice-receipt~ |
				| ~ND~ | Debit Note          | ~debit-note~  |                 |
				| ~NC~ | Credit Note         | ~credit-note~ |                 |

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
				The SAF-T's ~TaxExemptionCode~ (Código do motivo de isenção de imposto) is a code that
				specifies the reason the VAT tax is exempt in a Portuguese invoice. When the ~exempt~ tag
				is used in a tax combo, the ~ext~ map's ~pt-saft-exemption~ property is required.

				The SAF-T's ~TaxExemptionReason~ (Motivo da isenção de imposto) is a text that justifies
				the exemption referencing the relevant legislation. In GOBL, this is provided with a
				special line-level ~note~.

				By default, if no note is provided, GOBL will automatically assign a default one
				consistent with the exemption code. However, note that default texts are generic and not
				always sufficiently precise to comply with the regulations. The invoice issuer should be
				given the option to provide a custom note with the appropriate descriptive text.

				For example, you could define an invoice line exempt of tax with the exemption code and
				the reason as follows:

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
											"pt-saft-exemption": "M40"
										}
								}
							]
						},
						"notes": [
							{
								"key": "legal",
								"code": "M40",
								"src": "pt-saft-exemption",
								"text": "Artigo 6.º n.º 6 alínea a) do CIVA, a contrário"
							}
						]
					]
				}
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "AT Tax Exemption Codes",
					i18n.PT: "Códigos de motivo de isenção",
				},
				URL:         "https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Fatcorews/Documents/Tabela_Codigos_Motivo_Isencao.pdf",
				ContentType: "application/pdf",
			},
			{
				Title: i18n.String{
					i18n.EN: "Art. 2.2.14 of Despacho nº8632/2014",
					i18n.PT: "Art. 2.2.14 do Despacho nº8632/2014",
				},
				URL:         "https://files.diariodarepublica.pt/2s/2014/07/126000000/1725517261.pdf",
				ContentType: "application/pdf",
			},
			{
				Title: i18n.String{
					i18n.EN: "Field 4.4.19.7 of Portaria nº302/2016",
					i18n.PT: "Campo 4.4.19.7 da Portaria nº302/2016",
				},
				URL:         "https://files.diariodarepublica.pt/1s/2016/12/23100/0427304379.pdf",
				ContentType: "application/pdf",
			},
		},
		Values: []*cbc.Definition{
			// Names adapted to meet the requirements of art. 2.2.14 of
			// Despacho nº8632/2014 (for printed invoices) and campo 4.4.19.7
			// of Portaria nº302/2016 (for SAF-T files).
			{
				Code: "M01",
				Name: i18n.String{
					i18n.EN: "Article 16, No. 6, paragraphs a) to d) of the VAT code",
					i18n.PT: "Artigo 16.º, n.º 6, alíneas a) a d) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Amounts invoiced as being received on behalf of the purchaser
						— e.g. agency sales or third-party collections; supplier records
						funds in “accounts of third parties”, no VAT is added by the
						supplier.
					`),
				},
			},
			{
				Code: "M02",
				Name: i18n.String{
					i18n.EN: "Article 6 of Decree-Law No. 198/90 of 19th June",
					i18n.PT: "Artigo 6.º do Decreto-Lei n.º 198/90, de 19 de junho",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						B2B deliveries above €1,000 from a domestic supplier to a
						taxable person in Portugal who exports in the same state — VAT
						zero-rated when export is immediate.
					`),
				},
			},
			{
				Code: "M04",
				Name: i18n.String{
					i18n.EN: "Article 13 of the VAT code",
					i18n.PT: "Artigo 13.º do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Certain types of importation or re-importation expressly exempt
						under Portuguese VAT rules.
					`),
				},
			},
			{
				Code: "M05",
				Name: i18n.String{
					i18n.EN: "Article 14 of the VAT code",
					i18n.PT: "Artigo 14.º do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Exports, deemed exports or international transport services
						provided from Portugal — zero-rated for VAT.
					`),
				},
			},
			{
				Code: "M06",
				Name: i18n.String{
					i18n.EN: "Article 15 of the VAT code",
					i18n.PT: "Artigo 15.º do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Goods or services under suspension regimes (e.g. transhipment,
						processing), or special transitional exemption schemes.
					`),
				},
			},
			{
				Code: "M07",
				Name: i18n.String{
					i18n.EN: "Article 9 of the VAT code",
					i18n.PT: "Artigo 9.º do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Activities such as health, education, culture, insurance,
						rentals, or authorized lotteries — exempt goods and services, but
						without VAT recovery rights.
					`),
				},
			},
			{
				Code: "M09",
				Name: i18n.String{
					i18n.EN: "Article 62 paragraph b) of the VAT code",
					i18n.PT: "Artigo 62.º alínea b) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Retailers under the “small trader regime” (volume below threshold)
						who cannot recover VAT on inputs. The VAT is not levied, but the
						taxpayer has no deduction rights.
					`),
				},
			},
			{
				Code: "M10",
				Name: i18n.String{
					i18n.EN: "Article 57 of the VAT code",
					i18n.PT: "Artigo 57.º do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Flat-rate VAT scheme for small taxpayers (typically individuals
						with turnover below €15,000 in the preceding year). Exemption
						applies, albeit with a limited VAT recovery regime.
					`),
				},
			},
			{
				Code: "M11",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 346/85 of 23rd August",
					i18n.PT: "Decreto-Lei n.º 346/85, de 23 de agosto",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Suppliers or resellers of manufactured tobacco or matches taxed
						under the special tobacco tax regime — VAT is either not due or
						treated under excise rules.
					`),
				},
			},
			{
				Code: "M12",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 221/85 of 3rd July",
					i18n.PT: "Decreto-Lei n.º 221/85, de 3 de julho",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Travel agencies selling tours in their own name (even when
						booking services from others) are taxed only on their margin, not
						on the full amount.
					`),
				},
			},
			{
				Code: "M13",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Decreto-Lei n.º 199/96, de 18 de outubro",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Resale of used goods under the margin scheme — VAT applies solely
						to the profit margin.
					`),
				},
			},
			{
				Code: "M14",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Decreto-Lei n.º 199/96, de 18 de outubro",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Art pieces sold under the margin taxation rule — VAT calculated
						only on the dealer's margin.
					`),
				},
			},
			{
				Code: "M15",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Decreto-Lei n.º 199/96, de 18 de outubro",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Antiques or collectibles treated under a margin-based VAT regime,
						meaning VAT on margin only.
					`),
				},
			},
			{
				Code: "M16",
				Name: i18n.String{
					i18n.EN: "Article 14 of the RITI",
					i18n.PT: "Artigo 14.º do RITI",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Standard rule for intra-EU goods transfers (Normal
						intra-Community supplies) — zero-VAT, no import VAT regime;
						recipient accounts for VAT via VAT return.
					`),
				},
			},
			{
				Code: "M19",
				Name: i18n.String{
					i18n.EN: "Other exemptions",
					i18n.PT: "Outras isenções",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Temporary, sector-specific or exceptional exemptions established
						by specific legislation (e.g. emergency relief).
					`),
				},
			},
			{
				Code: "M20",
				Name: i18n.String{
					i18n.EN: "Article 59-D No. 2 of the VAT code",
					i18n.PT: "Artigo 59.º-D n.º 2 do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Flat-rate VAT exemption for certain agricultural activities under
						the “forfetário” scheme (simplified turnover-based regime).
					`),
				},
			},
			{
				Code: "M21",
				Name: i18n.String{
					i18n.EN: "Article 72 No. 4 of the VAT code",
					i18n.PT: "Artigo 72.º n.º 4 do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Fuel resellers operating on behalf of distributors cannot deduct
						VAT they incur — vendor issues invoice at 0%, no VAT recovery by
						reseller.
					`),
				},
			},
			{
				Code: "M25",
				Name: i18n.String{
					i18n.EN: "Article 38 No. 1 paragraph a) of the VAT code",
					i18n.PT: "Artigo 38.º n.º 1 alínea a) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Consignment stock sales — consignor issues invoice to consignee,
						but goods remain under consignment; VAT is not yet due until
						final sale.
					`),
				},
			},
			{
				Code: "M26",
				Name: i18n.String{
					i18n.EN: "Law No. 17/2023 of 14th April",
					i18n.PT: "Lei n.º 17/2023, de 14 de abril",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						“IVA Zero” food-basket scheme allowed VAT exemption with right to
						deduction. Active 18 Apr 2023 to 4 Jan 2024, now expired.
					`),
				},
			},
			{
				Code: "M30",
				Name: i18n.String{
					i18n.EN: "Article 2 No. 1 paragraph i) of the VAT code",
					i18n.PT: "Artigo 2.º n.º 1 alínea i) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge on purchase of waste, scrap or recyclable materials
						listed in Annex E of VAT Code — purchaser self-accounts VAT.
					`),
				},
			},
			{
				Code: "M31",
				Name: i18n.String{
					i18n.EN: "Article 2 No. 1 paragraph j) of the VAT code",
					i18n.PT: "Artigo 2.º n.º 1 alínea j) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge for construction or civil-works services under
						subcontracting; client self-accounts VAT.
					`),
				},
			},
			{
				Code: "M32",
				Name: i18n.String{
					i18n.EN: "Article 2 No. 1 paragraph l) of the VAT code",
					i18n.PT: "Artigo 2.º n.º 1 alínea l) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge on services involving emission rights, certified
						reductions or reduction units (environmental credits).
					`),
				},
			},
			{
				Code: "M33",
				Name: i18n.String{
					i18n.EN: "Article 2 No. 1 paragraph m) of the VAT code",
					i18n.PT: "Artigo 2.º n.º 1 alínea m) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge on raw cork, wood, pine-cones with shell —
						purchaser accounts for VAT.
					`),
				},
			},
			{
				Code: "M34",
				Name: i18n.String{
					i18n.EN: "Article 2 No. 1 paragraph n) of the VAT code",
					i18n.PT: "Artigo 2.º n.º 1 alínea n) do CIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge on electricity sold for self-consumption (≤1MW
						capacity); the supplier issues invoice at 0% VAT.
					`),
				},
			},
			{
				Code: "M40",
				Name: i18n.String{
					i18n.EN: "Article 6 No. 6 paragraph a) of the VAT code, to the contrary",
					i18n.PT: "Artigo 6.º n.º 6 alínea a) do CIVA, a contrário",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge for cross-border B2B services from suppliers
						without a fixed establishment in Portugal — Portuguese client
						self-accounts.
					`),
				},
			},
			{
				Code: "M41",
				Name: i18n.String{
					i18n.EN: "Article 8 No. 3 of the RITI",
					i18n.PT: "Artigo 8.º n.º 3 do RITI",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge for intracommunity triangular transactions (RITI
						art. 8.3) — recipient (third party) must self-liquidate VAT.
					`),
				},
			},
			{
				Code: "M42",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 21/2007 of 29th January",
					i18n.PT: "Decreto-Lei n.º 21/2007, de 29 de janeiro",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When supplier renounces to VAT exemption on a real-estate sale,
						the buyer must self-account VAT.
					`),
				},
			},
			{
				Code: "M43",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 362/99 of 16th September",
					i18n.PT: "Decreto-Lei n.º 362/99, de 16 de setembro",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Reverse-charge on investment gold transactions that are not
						VAT-exempt under special gold regulations; buyer reports VAT.
					`),
				},
			},
			{
				Code: "M44",
				Name: i18n.String{
					i18n.EN: "Article 6.º of the CIVA – Specific rules",
					i18n.PT: "Artigo 6.º do CIVA – Regras específicas",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						To be used in services not provided in Portugal under the terms of
						paragraphs 7 and following of article 6.º of the VAT Code.
					`),
				},
			},
			{
				Code: "M45",
				Name: i18n.String{
					i18n.EN: "Art 58º-A of the CIVA (IVA - cross-border exemption regime)",
					i18n.PT: "Art 58º-A do CIVA (IVA - regime transfronteiriço de isenção)",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						To be used in operations located in another Member State of the European
						Union that are exempt from VAT there, by virtue of the supplier of goods
						or provider of services having adhered to the Cross-Border Exemption
						Regime for operations carried out in that Member State.
					`),
				},
			},
			{
				Code: "M46",
				Name: i18n.String{
					i18n.EN: "Decree-Law No. 19/2017, of February 14",
					i18n.PT: "Decreto-lei n.º 19/2017, de 14 de fevereiro",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						To be used in the transfer of goods in the personal luggage of travelers
						without domicile in the European Union (a.k.a. tax free), under the
						terms of the mentioned decree-law.
					`),
				},
			},
			{
				Code: "M99",
				Name: i18n.String{
					i18n.EN: "Not subject to tax or not taxed",
					i18n.PT: "Não sujeito ou não tributado",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Other cases where VAT is legally not due — e.g. not subject due
						to CIVA articles 2, 3, or 4 (outside scope of VAT).
					`),
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
		Key: ExtKeyWorkType,
		Name: i18n.String{
			i18n.EN: "Document Type",
			i18n.PT: "Tipo de documento",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				SAF-T's ~WorkType~ (Tipo de documento de conferência) specifies the type of a working
				document. In GOBL, this type can be set using the ~pt-saft-work-type~ extension in either
				~bill.Invoice~ or ~bill.Order~ documents. GOBL will set the extension for you based on the
				document type in some cases.

				The table below shows the supported work type codes and their compatibility with GOBL objects:

				| Code   | Name                            | GOBL Doc | GOBL Type  |
				| ------ | ------------------------------- | -------- | ---------- |
				| ~PF~   | Pró-forma                       | Invoice  | ~proforma~ |
				| ~FC~   | Fatura de consignação           | Invoice  |            |
				| ~CC~   | Credito de consignação          | Invoice  |            |
				| ~CM~   | Consultas de mesa               | Order    |            |
				| ~FO~   | Folhas de obra                  | Order    |            |
				| ~NE~   | Nota de Encomenda               | Order    | ~purchase~ |
				| ~OU~   | Outros                          | Order    |            |
				| ~OR~   | Orçamentos                      | Order    | ~quote~    |
				| ~DC~   | Documentos de conferência       | Order    |            |
				| ~RP~   | Prémio ou recibo de prémio      | Order    |            |
				| ~RE~   | Estorno ou recibo de estorno    | Order    |            |
				| ~CS~   | Imputação a co-seguradoras      | Order    |            |
				| ~LD~   | Imputação a co-seguradora líder | Order    |            |
				| ~RA~   | Resseguro aceite                | Order    |            |

				Example for a proforma invoice:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"type": "proforma",
					// ...
					"tax": {
						"ext": {
							"pt-saft-work-type": "PF"
						}
					},
					// ...
				~~~

				Example for a purchase order:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/order",
					"type": "purchase",
					// ...
					"tax": {
						"ext": {
							"pt-saft-work-type": "NE"
						}
					},
					// ...
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: WorkTypeTableQueries,
				Name: i18n.String{
					i18n.EN: "Table orders",
					i18n.PT: "Consultas de mesa",
				},
			},
			{
				Code: WorkTypeConsignmentCredit,
				Name: i18n.String{
					i18n.EN: "Consignment credit note",
					i18n.PT: "Credito de consignação",
				},
			},
			{
				Code: WorkTypeConsignmentInv,
				Name: i18n.String{
					i18n.EN: "VAT-compliant consignment invoice (Article 38)",
					i18n.PT: "Fatura de consignação nos termos do art.º 38º do código do IVA",
				},
			},
			{
				Code: WorkTypeWorksheets,
				Name: i18n.String{
					i18n.EN: "Work orders",
					i18n.PT: "Folhas de obra",
				},
			},
			{
				Code: WorkTypePurchaseOrder,
				Name: i18n.String{
					i18n.EN: "Purchase order",
					i18n.PT: "Nota de Encomenda",
				},
			},
			{
				Code: WorkTypeOther,
				Name: i18n.String{
					i18n.EN: "Other documents",
					i18n.PT: "Outros",
				},
			},
			{
				Code: WorkTypeBudgets,
				Name: i18n.String{
					i18n.EN: "Quotations",
					i18n.PT: "Orçamentos",
				},
			},
			{
				Code: WorkTypeProforma,
				Name: i18n.String{
					i18n.EN: "Pro forma invoice",
					i18n.PT: "Pró-forma",
				},
			},
			{
				Code: WorkTypeDocuments,
				Name: i18n.String{
					i18n.EN: "Delivery verification documents",
					i18n.PT: "Documentos emitidos que sejam suscetíveis de apresentação ao cliente para conferência de mercadorias ou de prestação de serviços",
				},
				Desc: i18n.String{
					i18n.EN: "For data up to 2017-06-30",
					i18n.PT: "Para dados até 2017-06-30",
				},
			},
			{
				Code: WorkTypePremium,
				Name: i18n.String{
					i18n.EN: "Premium Receipt",
					i18n.PT: "Prémio ou recibo de prémio",
				},
			},
			{
				Code: WorkTypeChargeback,
				Name: i18n.String{
					i18n.EN: "Chargeback Receipt",
					i18n.PT: "Estorno ou recibo de estorno",
				},
			},
			{
				Code: WorkTypeCoInsurers,
				Name: i18n.String{
					i18n.EN: "Co-insurers Allocation",
					i18n.PT: "Imputação a co-seguradoras",
				},
			},
			{
				Code: WorkTypeLeadCoInsurer,
				Name: i18n.String{
					i18n.EN: "Lead Co-insurer Allocation",
					i18n.PT: "Imputação a co-seguradora líder",
				},
			},
			{
				Code: WorkTypeReinsurance,
				Name: i18n.String{
					i18n.EN: "Accepted Reinsurance",
					i18n.PT: "Resseguro aceite",
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
				value in your ~bill.Invoice~ advances or in you ~bill.Payment~ method. GOBL maps certain
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
				~~~
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
	{
		Key: ExtKeySource,
		Name: i18n.String{
			i18n.EN: "Document Source",
			i18n.PT: "Origem do documento",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				SAF-T's ~SourceBilling~ (Origem do documento) field specifies the source of the document.
				GOBL provides the ~pt-saft-source~ extension to set this value in your documents.
				By default, GOBL will set this extension to "P" (Produced).

				The table below shows the supported source billing codes:

				| Code | Name       | Description                                         |
				| ---- | ---------- | --------------------------------------------------- |
				| ~P~  | Produced   | Document produced by the application                |
				| ~I~  | Integrated | Integrated document produced by another application |
				| ~M~  | Manual     | Document from recovery or issued manually           |

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"ext": {
						"pt-saft-source": "P"
					},
					// ...
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: SourceBillingProduced,
				Name: i18n.String{
					i18n.EN: "Produced",
					i18n.PT: "Produzido",
				},
				Desc: i18n.String{
					i18n.EN: "Document produced by the application",
					i18n.PT: "Documento produzido na aplicação",
				},
			},
			{
				Code: SourceBillingIntegrated,
				Name: i18n.String{
					i18n.EN: "Integrated",
					i18n.PT: "Integrado",
				},
				Desc: i18n.String{
					i18n.EN: "Integrated document produced by another application",
					i18n.PT: "Documento integrado e produzido noutra aplicação",
				},
			},
			{
				Code: SourceBillingManual,
				Name: i18n.String{
					i18n.EN: "Manual",
					i18n.PT: "Manual",
				},
				Desc: i18n.String{
					i18n.EN: "Document from recovery or issued manually",
					i18n.PT: "Documento proveniente de recuperação ou de emissão manual",
				},
			},
		},
	},
	{
		Key: ExtKeySourceRef,
		Name: i18n.String{
			i18n.EN: "Source Document Reference",
			i18n.PT: "Referência do documento de origem",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				GOBL provides the ~pt-saft-source-ref~ extension to provide the full reference to a document
				integrated from another system, recovered or issued manually.

				This extension is required when the document source (~pt-saft-source~ extension) is
				"M" (manual) or "I" (integrated). It must contain the complete document reference to be appended
				to the SAF-T's ~HashControl~ field as stipulated by Despacho n.o 8632/2014.

				Example with a manual document:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"ext": {
						"pt-saft-source": "M",
						"pt-saft-source-ref": "FTM abc/00001"
					},
					// ...
				}
				~~~

				Example with a recovered document:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"ext": {
						"pt-saft-source": "M",
						"pt-saft-source-ref": "FTD FT SERIESA/123"
					},
					// ...
				}
				~~~
			`),
		},
	},
	{
		Key: ExtKeyMovementType,
		Name: i18n.String{
			i18n.EN: "Movement Type",
			i18n.PT: "Tipo de documento",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				SAF-T's ~MovementType~ (Tipo de documento de movimentação de mercadorias) specifies the type of
				a delivery document. In GOBL,this type can be set using the ~pt-saft-movement-type~ extension.
				If not provided explicitly, GOBL will set the extension for you based on the type of your delivery
				document.

				The table below shows how this mapping is done:

				| Code | Name                | GOBL Type     |
				| ---- | ------------------- | ------------- |
				| ~GR~ | Delivery note       | ~note~        |
				| ~GT~ | Waybill             | ~waybill~     |

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/delivery",
					// ...
					"type": "note",
					// ...
					"ext": {
						"pt-saft-movement-type": "GR"
					},
					// ...
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: MovementTypeDeliveryNote,
				Name: i18n.String{
					i18n.EN: "Delivery note",
					i18n.PT: "Guia de remessa",
				},
			},
			{
				Code: MovementTypeWaybill,
				Name: i18n.String{
					i18n.EN: "Waybill",
					i18n.PT: "Guia de transporte",
				},
				Desc: i18n.String{
					i18n.EN: "Include global waybills here",
					i18n.PT: "Incluir aqui as guias globais",
				},
			},
			{
				Code: MovementTypeFixedAssets,
				Name: i18n.String{
					i18n.EN: "Guide to the movement own fixed assets",
					i18n.PT: "Guia de movimentação de ativos fixos próprios",
				},
			},
			{
				Code: MovementTypeConsignment,
				Name: i18n.String{
					i18n.EN: "Consignment note",
					i18n.PT: "Guia de consignação",
				},
			},
			{
				Code: MovementTypeReturn,
				Name: i18n.String{
					i18n.EN: "Returns slip or note",
					i18n.PT: "Guia ou nota de devolução",
				},
			},
		},
	},
}
