package pt

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.PT: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.PT: "Imposto sobre o Valor Acrescentado",
		},
		Retained: false,
		Extensions: []cbc.Key{
			ExtKeyRegion,
			ExtKeySAFTTaxRate,
			ExtKeyExemptionCode,
		},
		Validation: func(c *tax.Combo) error {
			return validation.ValidateStruct(c,
				validation.Field(&c.Ext,
					// NOTE! We know that some tax rate is required in portugal, but
					// we don't know what it should be for foreign countries.
					// Until this is known, we're removing the validation for the
					// country tax rate.
					validation.When(
						c.Country == "",
						tax.ExtensionsRequires(ExtKeySAFTTaxRate),
					),
					validation.When(
						c.Percent == nil,
						tax.ExtensionsRequires(ExtKeyExemptionCode),
					),
					validation.Skip,
				),
			)
		},
		Rates: []*tax.Rate{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.PT: "Tipo Geral",
				},
				Ext: tax.Extensions{
					ExtKeySAFTTaxRate: "NOR",
				},
				Values: []*tax.RateValue{
					{
						Ext: tax.Extensions{
							ExtKeyRegion: "PT-AC",
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(160, 3),
					},
					{
						Ext: tax.Extensions{
							ExtKeyRegion: "PT-MA",
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(220, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(230, 3),
					},
				},
			},
			{
				Key: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.PT: "Taxa Intermédia", //nolint:misspell
				},
				Ext: tax.Extensions{
					ExtKeySAFTTaxRate: "INT",
				},
				Values: []*tax.RateValue{
					{
						Ext: tax.Extensions{
							ExtKeyRegion: "PT-AC",
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(90, 3),
					},
					{
						Ext: tax.Extensions{
							ExtKeyRegion: "PT-MA",
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(130, 3),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.PT: "Taxa Reduzida",
				},
				Ext: tax.Extensions{
					ExtKeySAFTTaxRate: "RED",
				},
				Values: []*tax.RateValue{
					{
						Ext: tax.Extensions{
							ExtKeyRegion: "PT-AC",
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(40, 3),
					},
					{
						Ext: tax.Extensions{
							ExtKeyRegion: "PT-MA",
						},
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(50, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(60, 3),
					},
				},
			},
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PT: "Isento",
				},
				Exempt: true,
				Ext: tax.Extensions{
					ExtKeySAFTTaxRate: "ISE",
				},
			},
		},
	},
}
