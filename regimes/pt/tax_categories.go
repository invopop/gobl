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
				Map: cbc.CodeMap{
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
				Map: cbc.CodeMap{
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
				Map: cbc.CodeMap{
					KeyATTaxCode: TaxCodeReduced,
				},
			},
			{
				Key:    TaxRateExempt,
				Exempt: true,
				Map: cbc.CodeMap{
					KeyATTaxCode: TaxCodeExempt,
				},
				Extensions: []cbc.Key{ExtKeyExemptionReason},
			},
		},
	},
}
