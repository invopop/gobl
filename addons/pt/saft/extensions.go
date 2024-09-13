package saft

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// SAF-T Extension Keys
const (
	ExtKeyExemption   = "pt-saft-exemption" // note: avoid redundant prefixes like `-code`
	ExtKeyTaxRate     = "pt-saft-tax-rate"
	ExtKeyInvoiceType = "pt-saft-invoice-type"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyInvoiceType,
		Name: i18n.String{
			i18n.EN: "Invoice Type",
			i18n.PT: "Tipo de Fatura",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "FT",
				Name: i18n.String{
					i18n.EN: "Standard Invoice",
					i18n.PT: "Fatura",
				},
			},
			{
				Value: "FS",
				Name: i18n.String{
					i18n.EN: "Simplified Invoice",
					i18n.PT: "Fatura Simplificada",
				},
			},
			{
				Value: "FR",
				Name: i18n.String{
					i18n.EN: "Invoice-Receipt",
					i18n.PT: "Fatura-Recibo",
				},
			},
			{
				Value: "ND",
				Name: i18n.String{
					i18n.EN: "Debit Note",
					i18n.PT: "Nota de Débito",
				},
			},
			{
				Value: "NC",
				Name: i18n.String{
					i18n.EN: "Credit Note",
					i18n.PT: "Nota de Crédito",
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
		Values: []*cbc.ValueDefinition{
			{
				Value: "RED",
				Name: i18n.String{
					i18n.EN: "Reduced",
					i18n.PT: "Redução",
				},
			},
			{
				Value: "INT",
				Name: i18n.String{
					i18n.EN: "Intermediate",
					i18n.PT: "Intermédio",
				},
			},
			{
				Value: "NOR",
				Name: i18n.String{
					i18n.EN: "Normal",
					i18n.PT: "Normal",
				},
			},
			{
				Value: "ISE",
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PT: "Isento",
				},
			},
			{
				Value: "OUT",
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
		Values: []*cbc.ValueDefinition{
			{
				Value: "M01",
				Name: i18n.String{
					i18n.EN: "Article 16, No. 6 of the VAT code",
					i18n.PT: "Artigo 16.°, n.° 6 do CIVA",
				},
			},
			{
				Value: "M02",
				Name: i18n.String{
					i18n.EN: "Article 6 of the Decree-Law 198/90 of 19th June",
					i18n.PT: "Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho",
				},
			},
			{
				Value: "M04",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 13 of the VAT code",
					i18n.PT: "Isento artigo 13.° do CIVA",
				},
			},
			{
				Value: "M05",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 14 of the VAT code",
					i18n.PT: "Isento artigo 14.° do CIVA",
				},
			},
			{
				Value: "M06",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 15 of the VAT code",
					i18n.PT: "Isento artigo 15.° do CIVA",
				},
			},
			{
				Value: "M07",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 9 of the VAT code",
					i18n.PT: "Isento artigo 9.° do CIVA",
				},
			},
			{
				Value: "M09",
				Name: i18n.String{
					i18n.EN: "VAT - does not confer right to deduct / Article 62 paragraph b) of the VAT code",
					i18n.PT: "IVA - não confere direito a dedução / Artigo 62.° alínea b) do CIVA",
				},
			},
			{
				Value: "M10",
				Name: i18n.String{
					i18n.EN: "VAT - exemption scheme / Article 57 of the VAT code",
					i18n.PT: "IVA - regime de isenção / Artigo 57.° do CIVA",
				},
			},
			{
				Value: "M11",
				Name: i18n.String{
					i18n.EN: "Special scheme for tobacco / Decree-Law No. 346/85 of 23rd August",
					i18n.PT: "Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto",
				},
			},
			{
				Value: "M12",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Travel agencies / Decree-Law No. 221/85 of 3rd July",
					i18n.PT: "Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho",
				},
			},
			{
				Value: "M13",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Second-hand goods / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
			},
			{
				Value: "M14",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Works of art / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
			},
			{
				Value: "M15",
				Name: i18n.String{
					i18n.EN: "Margin scheme - Collector’s items and antiques / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
			},
			{
				Value: "M16",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 14 of the RITI",
					i18n.PT: "Isento artigo 14.° do RITI",
				},
			},
			{
				Value: "M19",
				Name: i18n.String{
					i18n.EN: "Other exemptions - Temporary exemptions determined by specific legislation",
					i18n.PT: "Outras isenções - Isenções temporárias determinadas em diploma próprio",
				},
			},
			{
				Value: "M20",
				Name: i18n.String{
					i18n.EN: "VAT - flat-rate scheme / Article 59-D No. 2 of the VAT code",
					i18n.PT: "IVA - regime forfetário / Artigo 59.°-D n.°2 do CIVA",
				},
			},
			{
				Value: "M21",
				Name: i18n.String{
					i18n.EN: "VAT - does not confer right to deduct (or similar) - Article 72 No. 4 of the VAT code",
					i18n.PT: "IVA - não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA",
				},
			},
			{
				Value: "M25",
				Name: i18n.String{
					i18n.EN: "Consignment goods - Article 38 No. 1 paragraph a) of the VAT code",
					i18n.PT: "Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA",
				},
			},
			{
				Value: "M30",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph i) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea i) do CIVA",
				},
			},
			{
				Value: "M31",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph j) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea j) do CIVA",
				},
			},
			{
				Value: "M32",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph l) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea I) do CIVA",
				},
			},
			{
				Value: "M33",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph m) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea m) do CIVA",
				},
			},
			{
				Value: "M40",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 6 No. 6 paragraph a) of the VAT code, to the contrary",
					i18n.PT: "IVA - autoliquidação / Artigo 6.° n.° 6 alínea a) do CIVA, a contrário",
				},
			},
			{
				Value: "M41",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 8 No. 3 of the RITI",
					i18n.PT: "IVA - autoliquidação / Artigo 8.° n.° 3 do RITI",
				},
			},
			{
				Value: "M42",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Decree-Law No. 21/2007 of 29 January",
					i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 21/2007, de 29 de janeiro",
				},
			},
			{
				Value: "M43",
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Decree-Law No. 362/99 of 16th September",
					i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 362/99, de 16 de setembro",
				},
			},
			{
				Value: "M99",
				Name: i18n.String{
					i18n.EN: "Not subject to tax or not taxed",
					i18n.PT: "Não sujeito ou não tributado",
				},
			},
		},
	},
}
