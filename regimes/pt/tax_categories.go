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
	TaxRateExempt        cbc.Key = "exempt"
	TaxRateOutlay        cbc.Key = "outlay"
	TaxRateIntrastate    cbc.Key = "intrastate-export"
	TaxRateImports       cbc.Key = "imports"
	TaxRateExports       cbc.Key = "exports"
	TaxRateSuspension    cbc.Key = "suspension-scheme"
	TaxRateInternalOps   cbc.Key = "internal-operations"
	TaxRateSmallRetail   cbc.Key = "small-retail-scheme"
	TaxRateExemptScheme  cbc.Key = "exempt-scheme"
	TaxRateTobacco       cbc.Key = "tobacco-scheme"
	TaxRateMargin        cbc.Key = "margin-scheme"
	TaxRateTravel        cbc.Key = "travel"
	TaxRateSecondHand    cbc.Key = "second-hand"
	TaxRateArt           cbc.Key = "art"
	TaxRateAntiques      cbc.Key = "antiques"
	TaxRateTransmission  cbc.Key = "goods-transmission"
	TaxRateOther         cbc.Key = "other"
	TaxRateFlatRate      cbc.Key = "flat-rate-scheme"
	TaxRateNonDeductible cbc.Key = "non-deductible"
	TaxRateConsignment   cbc.Key = "consignment-goods"
	TaxRateReverseCharge cbc.Key = "reverse-charge"
	TaxRateWaste         cbc.Key = "waste"
	TaxRateCivilEng      cbc.Key = "civil-eng"
	TaxRateGreenhouse    cbc.Key = "greenhouse"
	TaxRateWoods         cbc.Key = "woods"
	TaxRateB2B           cbc.Key = "b2b"
	TaxRateIntraEU       cbc.Key = "intraeu"
	TaxRateRealEstate    cbc.Key = "real-estate"
	TaxRateGold          cbc.Key = "gold"
	TaxRateNonTaxable    cbc.Key = "non-taxable"
)

// AT Tax Codes
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
		Retained: false,
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
						Percent: num.MakePercentage(16, 2),
					},
					{
						Zones:   []l10n.Code{ZoneMadeira},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(22, 2),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(23, 2),
					},
				},
				Codes: cbc.CodeSet{
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
						Percent: num.MakePercentage(9, 2),
					},
					{
						Zones:   []l10n.Code{ZoneMadeira},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(12, 2),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(13, 2),
					},
				},
				Codes: cbc.CodeSet{
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
						Percent: num.MakePercentage(4, 2),
					},
					{
						Zones:   []l10n.Code{ZoneMadeira},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(5, 2),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(6, 2),
					},
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode: TaxCodeReduced,
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateOutlay),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Article 16, No. 6 of the VAT code",
					i18n.PT: "Artigo 16.°, n.° 6 do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M01",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateIntrastate),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Article 6 of the Decree-Law 198/90 of 19th June",
					i18n.PT: "Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M02",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateImports),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 13 of the VAT code",
					i18n.PT: "Isento artigo 13.° do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M04",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateExports),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 14 of the VAT code",
					i18n.PT: "Isento artigo 14.° do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M05",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateSuspension),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 15 of the VAT code",
					i18n.PT: "Isento artigo 15.° do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M06",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateInternalOps),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to article 9 of the VAT code",
					i18n.PT: "Isento artigo 9.° do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M07",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateSmallRetail),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - does not confer right to deduct / Article 62 paragraph b) of the VAT code",
					i18n.PT: "IVA - não confere direito a dedução / Artigo 62.° alínea b) do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M09",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateExemptScheme),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - exemption scheme / Article 57 of the VAT code",
					i18n.PT: "IVA - regime de isenção / Artigo 57.° do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M10",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateTobacco),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Special scheme for tobacco / Decree-Law No. 346/85 of 23rd August",
					i18n.PT: "Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M11",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateMargin).With(TaxRateTravel),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Margin scheme - Travel agencies / Decree-Law No. 221/85 of 3rd July",
					i18n.PT: "Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M12",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateMargin).With(TaxRateSecondHand),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Margin scheme - Second-hand goods / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M13",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateMargin).With(TaxRateArt),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Margin scheme - Works of art / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M14",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateMargin).With(TaxRateAntiques),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Margin scheme - Collector’s items and antiques / Decree-Law No. 199/96 of 18th October",
					i18n.PT: "Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M15",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateTransmission),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 14 of the RITI",
					i18n.PT: "Isento artigo 14.° do RITI",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M16",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateOther),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Other exemptions - Temporary exemptions determined by specific legislation",
					i18n.PT: "Outras isenções - Isenções temporárias determinadas em diploma próprio",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M19",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateFlatRate),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - flat-rate scheme / Article 59-D No. 2 of the VAT code",
					i18n.PT: "IVA - regime forfetário / Artigo 59.°-D n.°2 do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M20",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateNonDeductible),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - does not confer right to deduct (or similar) - Article 72 No. 4 of the VAT code",
					i18n.PT: "IVA - não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M21",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateConsignment),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Consignment goods - Article 38 No. 1 paragraph a) of the VAT code",
					i18n.PT: "Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M25",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWaste),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph i) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea i) do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M30",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateCivilEng),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph j) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea j) do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M31",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGreenhouse),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph l) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea I) do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M32",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWoods),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 2 No. 1 paragraph m) of the VAT code",
					i18n.PT: "IVA - autoliquidação / Artigo 2.° n.° 1 alínea m) do CIVA",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M33",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateB2B),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 6 No. 6 paragraph a) of the VAT code, to the contrary",
					i18n.PT: "IVA - autoliquidação / Artigo 6.° n.° 6 alínea a) do CIVA, a contrário",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M40",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateIntraEU),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Article 8 No. 3 of the RITI",
					i18n.PT: "IVA - autoliquidação / Artigo 8.° n.° 3 do RITI",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M41",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateRealEstate),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Decree-Law No. 21/2007 of 29 January",
					i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 21/2007, de 29 de janeiro",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M42",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGold),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "VAT - reverse charge / Decree-Law No. 362/99 of 16th September",
					i18n.PT: "IVA - autoliquidação / Decreto-Lei n.° 362/99, de 16 de setembro",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M43",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateNonTaxable),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Not subject to tax or not taxed",
					i18n.PT: "Não sujeito ou não tributado",
				},
				Codes: cbc.CodeSet{
					KeyATTaxCode:          TaxCodeExempt,
					KeyATTaxExemptionCode: "M99",
				},
			},
		},
	},
}
