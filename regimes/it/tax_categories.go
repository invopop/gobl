package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
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

// Keys used for meta data from external sources.
const (
	KeyFatturaPATipoRitenuta      cbc.Key = "fatturapa-tipo-ritenuta"
	KeyFatturaPAModalitaPagamento cbc.Key = "fatturapa-modalita-pagamento"
)

var categories = []*tax.CategoryDef{
	{
		Code:     tax.CategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.IT: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.IT: "Imposta sul Valore Aggiunto",
		},
		Keys: tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Minimum Rate",
					i18n.IT: "Aliquota Minima",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(40, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.IT: "Aliquota Ridotta",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.IT: "Aliquota Intermedia",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.IT: "Aliquota Ordinaria",
				},
				Values: []*tax.RateValueDef{
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
		Title: i18n.String{
			i18n.EN: "Personal Income Tax",
			i18n.IT: "Imposta sul Reddito delle Persone Fisiche",
		},
		Map: cbc.CodeMap{
			KeyFatturaPATipoRitenuta: "RT01",
		},
	},
	{
		Code:     TaxCategoryIRES,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRES",
			i18n.IT: "IRES",
		},
		Title: i18n.String{
			i18n.EN: "Corporate Income Tax",
			i18n.IT: "Imposta sul Reddito delle Società",
		},
		Map: cbc.CodeMap{
			KeyFatturaPATipoRitenuta: "RT02",
		},
	},
	{
		Code:     TaxCategoryINPS,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "INPS Contribution",
			i18n.IT: "Contributo INPS", // nolint:misspell
		},
		Title: i18n.String{
			i18n.EN: "Contribution to the National Social Security Institute",
			i18n.IT: "Contributo Istituto Nazionale della Previdenza Sociale", // nolint:misspell
		},
		Map: cbc.CodeMap{
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
		Title: i18n.String{
			i18n.EN: "Contribution to the National Welfare Board for Sales Agents and Representatives",
			i18n.IT: "Contributo Ente Nazionale Assistenza Agenti e Rappresentanti di Commercio", // nolint:misspell
		},
		Map: cbc.CodeMap{
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
		Title: i18n.String{
			i18n.EN: "Contribution to the National Pension and Welfare Board for Doctors",
			i18n.IT: "Contributo - Ente Nazionale Previdenza e Assistenza Medici", // nolint:misspell
		},
		Map: cbc.CodeMap{
			KeyFatturaPATipoRitenuta: "RT05",
		},
	},
}
