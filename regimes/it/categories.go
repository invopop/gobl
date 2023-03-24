package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Local tax category definitions which are not considered standard.
// There is a 6th retained tax type, RT06 "Other contributions", which is
// currently not supported.
const (
	// https://www.agenziaentrate.gov.it/portale/imposta-sul-reddito-delle-persone-fisiche-irpef-/aliquote-e-calcolo-dell-irpef
	TaxCategoryIRPEF    cbc.Code = "IRPEF"
	TaxCategoryIRES     cbc.Code = "IRES"
	TaxCategoryINPS     cbc.Code = "INPS"
	TaxCategoryENASARCO cbc.Code = "ENASARCO"
	TaxCategoryENPAM    cbc.Code = "ENPAM"
)

var categories = []*tax.Category{
	{
		Code:     common.TaxCategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.IT: "IVA",
		},
		Desc: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.IT: "Imposta sul Valore Aggiunto",
		},
		Tags: vatTaxTags,
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.IT: "Aliquota Zero",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: common.TaxRateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Minimum Rate",
					i18n.IT: "Aliquota Minima",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(40, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.IT: "Aliquota Ridotta",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			{
				Key: common.TaxRateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.IT: "Aliquota Intermedia",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Ordinary Rate",
					i18n.IT: "Aliquota Ordinaria",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(220, 3),
					},
				},
			},
		},
	},
	{
		// IT: https://www.agenziaentrate.gov.it/portale/imposta-sul-reddito-delle-persone-fisiche-irpef-/aliquote-e-calcolo-dell-irpef
		// EN: https://www.agenziaentrate.gov.it/portale/web/english/information-for-specific-categories-of-workers
		Code:     TaxCategoryIRPEF,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRPEF",
			i18n.IT: "IRPEF",
		},
		Desc: i18n.String{
			i18n.EN: "Personal Income Tax",
			i18n.IT: "Imposta sul Reddito delle Persone Fisiche",
		},
		Tags: retainedTaxTags,
		Meta: cbc.Meta{
			KeyFatturaPATipoRitenuta: "RT01",
		},
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Ordinary Rate",
					i18n.IT: "Aliquota Ordinaria",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(200, 3),
					},
				},
			},
		},
	},
	{
		Code:     TaxCategoryIRES,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRES",
			i18n.IT: "IRES",
		},
		Desc: i18n.String{
			i18n.EN: "Corporate Income Tax",
			i18n.IT: "Imposta sul Reddito delle Società",
		},
		Tags: retainedTaxTags,
		Meta: cbc.Meta{
			KeyFatturaPATipoRitenuta: "RT02",
		},
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Ordinary Rate",
					i18n.IT: "Aliquota Ordinaria",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(240, 3),
					},
				},
			},
		},
	},
	{
		Code:     TaxCategoryINPS,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "INPS Contribution",
			i18n.IT: "Contributo INPS", // nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Contribution to the National Social Security Institute",
			i18n.IT: "Contributo Istituto Nazionale della Previdenza Sociale", // nolint:misspell
		},
		Tags: retainedTaxTags,
		Meta: cbc.Meta{
			KeyFatturaPATipoRitenuta: "RT03",
		},
	},
	{
		Code:     TaxCategoryENASARCO,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "ENASARCO Contribution",
			i18n.IT: "Contributo ENASARCO", // nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Contribution to the National Welfare Board for Sales Agents and Representatives",
			i18n.IT: "Contributo Ente Nazionale Assistenza Agenti e Rappresentanti di Commercio", // nolint:misspell
		},
		Tags: retainedTaxTags,
		Meta: cbc.Meta{
			KeyFatturaPATipoRitenuta: "RT04",
		},
	},
	{
		Code:     TaxCategoryENPAM,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "ENPAM Contribution",
			i18n.IT: "Contributo ENPAM", // nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Contribution to the National Pension and Welfare Board for Doctors",
			i18n.IT: "Contributo - Ente Nazionale Previdenza e Assistenza Medici", // nolint:misspell
		},
		Tags: retainedTaxTags,
		Meta: cbc.Meta{
			KeyFatturaPATipoRitenuta: "RT05",
		},
	},
}
