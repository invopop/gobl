package pt

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Tax rate exemption tags
const (
	TaxRateExempt cbc.Key = "exempt"
)

// AT Tax Map
const (
	TaxCodeStandard     cbc.Code = "NOR"
	TaxCodeIntermediate cbc.Code = "INT"
	TaxCodeReduced      cbc.Code = "RED"
	TaxCodeExempt       cbc.Code = "ISE"
	TaxCodeOther        cbc.Code = "OUT"
)

var taxCategories = []*tax.Category{
	// VAT
	{
		Code: common.TaxCategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.PT: "IVA",
		},
		Desc: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.PT: "Imposto sobre o Valor Acrescentado",
		},
		Retained:     false,
		RateRequired: true,
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.PT: "Tipo Geral",
				},
				Values: []*tax.RateValue{
					{
						Zones:   []l10n.Code{ZoneAzores},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(160, 3),
					},
					{
						Zones:   []l10n.Code{ZoneMadeira},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(220, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(230, 3),
					},
				},
				Map: cbc.CodeSet{
					KeyATTaxCode: TaxCodeStandard,
				},
			},
			{
				Key: common.TaxRateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.PT: "Taxa Intermédia", //nolint:misspell
				},
				Values: []*tax.RateValue{
					{
						Zones:   []l10n.Code{ZoneAzores},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(90, 3),
					},
					{
						Zones:   []l10n.Code{ZoneMadeira},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(130, 3),
					},
				},
				Map: cbc.CodeSet{
					KeyATTaxCode: TaxCodeIntermediate,
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.PT: "Taxa Reduzida",
				},
				Values: []*tax.RateValue{
					{
						Zones:   []l10n.Code{ZoneAzores},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(40, 3),
					},
					{
						Zones:   []l10n.Code{ZoneMadeira},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(50, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(60, 3),
					},
				},
				Map: cbc.CodeSet{
					KeyATTaxCode: TaxCodeReduced,
				},
			},
			{
				Key:    TaxRateExempt,
				Exempt: true,
				Map: cbc.CodeSet{
					KeyATTaxCode: TaxCodeExempt,
				},
				Codes: []*tax.CodeDefinition{
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
		},
	},
}
