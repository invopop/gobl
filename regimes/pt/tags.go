package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Tax exemption tags
const (
	TagExempt cbc.Key = "exempt"
	TagM01    cbc.Key = "m01"
	TagM02    cbc.Key = "m02"
	TagM04    cbc.Key = "m04"
	TagM05    cbc.Key = "m05"
	TagM06    cbc.Key = "m06"
	TagM07    cbc.Key = "m07"
	TagM09    cbc.Key = "m09"
	TagM10    cbc.Key = "m10"
	TagM11    cbc.Key = "m11"
	TagM12    cbc.Key = "m12"
	TagM13    cbc.Key = "m13"
	TagM14    cbc.Key = "m14"
	TagM15    cbc.Key = "m15"
	TagM16    cbc.Key = "m16"
	TagM19    cbc.Key = "m19"
	TagM20    cbc.Key = "m20"
	TagM21    cbc.Key = "m21"
	TagM25    cbc.Key = "m25"
	TagM30    cbc.Key = "m30"
	TagM31    cbc.Key = "m31"
	TagM32    cbc.Key = "m32"
	TagM33    cbc.Key = "m33"
	TagM40    cbc.Key = "m40"
	TagM41    cbc.Key = "m41"
	TagM42    cbc.Key = "m42"
	TagM43    cbc.Key = "m43"
	TagM99    cbc.Key = "m99"
)

var vatTaxTags = []*tax.Tag{
	{
		Key: TagExempt.With(TagM01),
		Name: i18n.String{
			i18n.EN: "Article 16, No. 6 of the VAT code",
			i18n.PT: "Artigo 16.°, n.° 6 do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M01",
		},
	},
	{
		Key: TagExempt.With(TagM02),
		Name: i18n.String{
			i18n.EN: "Article 6 of the Decree-Law 198/90 of 19th June",
			i18n.PT: "Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M02",
		},
	},
	{
		Key: TagExempt.With(TagM04),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to article 13 of the VAT code",
			i18n.PT: "Isento artigo 13.° do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M04",
		},
	},
	{
		Key: TagExempt.With(TagM05),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to article 14 of the VAT code",
			i18n.PT: "Isento artigo 14.° do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M05",
		},
	},
	{
		Key: TagExempt.With(TagM06),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to article 15 of the VAT code",
			i18n.PT: "Isento artigo 15.° do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M06",
		},
	},
	{
		Key: TagExempt.With(TagM07),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to article 9 of the VAT code",
			i18n.PT: "Isento artigo 9.° do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M07",
		},
	},
	{
		Key: TagExempt.With(TagM09),
		Name: i18n.String{
			i18n.EN: "VAT - does not confer right to deduct / Article 62 paragraph b) of the VAT code",
			i18n.PT: "IVA - não confere direito a dedução / Artigo 62.° alínea b) do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M09",
		},
	},
	{
		Key: TagExempt.With(TagM10),
		Name: i18n.String{
			i18n.EN: "VAT - exemption scheme / Article 57 of the VAT code",
			i18n.PT: "IVA - regime de isenção / Artigo 57.° do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M10",
		},
	},
	{
		Key: TagExempt.With(TagM11),
		Name: i18n.String{
			i18n.EN: "Special scheme for tobacco / Decree-Law No. 346/85 of 23rd August",
			i18n.PT: "Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M11",
		},
	},
	{
		Key: TagExempt.With(TagM12),
		Name: i18n.String{
			i18n.EN: "Margin scheme - Travel agencies / Decree-Law No. 221/85 of 3rd July",
			i18n.PT: "Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M12",
		},
	},
	{
		Key: TagExempt.With(TagM13),
		Name: i18n.String{
			i18n.EN: "Margin scheme - Second-hand goods / Decree-Law No. 199/96 of 18th October",
			i18n.PT: "Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M13",
		},
	},
	{
		Key: TagExempt.With(TagM14),
		Name: i18n.String{
			i18n.EN: "Margin scheme - Works of art / Decree-Law No. 199/96 of 18th October",
			i18n.PT: "Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M14",
		},
	},
	{
		Key: TagExempt.With(TagM15),
		Name: i18n.String{
			i18n.EN: "Margin scheme - Collector’s items and antiques / Decree-Law No. 199/96 of 18th October",
			i18n.PT: "Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M15",
		},
	},
	{
		Key: TagExempt.With(TagM16),
		Name: i18n.String{
			i18n.EN: "Exempt pursuant to Article 14 of the RITI",
			i18n.PT: "Isento artigo 14.° do RITI",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M16",
		},
	},
	{
		Key: TagExempt.With(TagM19),
		Name: i18n.String{
			i18n.EN: "Other exemptions - Temporary exemptions determined by specific legislation",
			i18n.PT: "Outras isenções - Isenções temporárias determinadas em diploma próprio",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M19",
		},
	},
	{
		Key: TagExempt.With(TagM20),
		Name: i18n.String{
			i18n.EN: "VAT - flat-rate scheme / Article 59-D No. 2 of the VAT code",
			i18n.PT: "IVA - regime forfetário / Artigo 59.°-D n.°2 do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M20",
		},
	},
	{
		Key: TagExempt.With(TagM21),
		Name: i18n.String{
			i18n.EN: "VAT - does not confer right to deduct (or similar) - Article 72 No. 4 of the VAT code",
			i18n.PT: "IVA - não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M21",
		},
	},
	{
		Key: TagExempt.With(TagM25),
		Name: i18n.String{
			i18n.EN: "Consignment goods - Article 38 No. 1 paragraph a) of the VAT code",
			i18n.PT: "Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M25",
		},
	},
	{
		Key: TagExempt.With(TagM30),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph i) of the VAT code",
			i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea i) do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M30",
		},
	},
	{
		Key: TagExempt.With(TagM31),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph j) of the VAT code",
			i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea j) do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M31",
		},
	},
	{
		Key: TagExempt.With(TagM32),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph l) of the VAT code",
			i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea I) do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M32",
		},
	},
	{
		Key: TagExempt.With(TagM33),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph m) of the VAT code",
			i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea m) do CIVA",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M33",
		},
	},
	{
		Key: TagExempt.With(TagM40),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Article 6 No. 6 paragraph a) of the VAT code, to the contrary",
			i18n.PT: "IVA - autoliquidação / Artigo 6.° n.° 6 alínea a) do CIVA, a contrário",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M40",
		},
	},
	{
		Key: TagExempt.With(TagM41),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Article 8 No. 3 of the RITI",
			i18n.PT: "IVA - autoliquidação / Artigo 8.° n.° 3 do RITI",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M41",
		},
	},
	{
		Key: TagExempt.With(TagM42),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Decree-Law No. 21/2007 of 29 January",
			i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 21/2007, de 29 de janeiro",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M42",
		},
	},
	{
		Key: TagExempt.With(TagM43),
		Name: i18n.String{
			i18n.EN: "VAT - reverse charge / Decree-Law No. 362/99 of 16th September",
			i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 362/99, de 16 de setembro",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M43",
		},
	},
	{
		Key: TagExempt.With(TagM99),
		Name: i18n.String{
			i18n.EN: "Not subject to tax or not taxed",
			i18n.PT: "Não sujeito ou não tributado",
		},
		Meta: cbc.Meta{
			KeyTaxExemptionCode: "M99",
		},
	},
}
