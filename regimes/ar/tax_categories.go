package ar

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Tax categories specific for Argentina that extend the common
// regime codes.
const (
	TaxCategoryRIVA cbc.Code = "RIVA" // IVA Retenido (Retained VAT)
	TaxCategoryRG   cbc.Code = "RG"   // Ganancias (Income Tax Withholding)
	TaxCategoryIB   cbc.Code = "IB"   // Ingresos Brutos (Gross Income Tax - Provincial)
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		//
		// IVA (Impuesto al Valor Agregado)
		//
		{
			Code:     tax.CategoryVAT,
			Retained: false,
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.ES: "IVA",
			},
			Title: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.ES: "Impuesto al Valor Agregado",
			},
			Description: &i18n.String{
				i18n.EN: here.Doc(`
					Known in Spanish as "Impuesto al Valor Agregado" (IVA), Argentina's VAT is a
					consumption tax applied to the sale of goods and services. The standard rate is 21%,
					with reduced rates of 10.5% for certain essential goods and services, and 0% for
					exports and specific categories.
				`),
			},
			Keys: tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "General Rate",
						i18n.ES: "Alícuota General",
					},
					Values: []*tax.RateValueDef{
						{
							// Standard VAT rate - 21%
							// Reference: AFIP - Law 23.349 and modifications
							// Source: https://www.avalara.com/vatlive/en/country-guides/south-america/argentina/argentina-vat-compliance-and-rates.html
							// Source: https://santandertrade.com/en/portal/establish-overseas/argentina/tax-system
							Since:   cal.NewDate(1995, 4, 1),
							Percent: num.MakePercentage(210, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced Rate",
						i18n.ES: "Alícuota Reducida",
					},
					Values: []*tax.RateValueDef{
						{
							// Reduced VAT rate - 10.5%
							// Applied to: construction work, medicine, public transportation,
							// livestock and agricultural products for food, etc.
							// Reference: AFIP - Law 23.349 Article 28
							// Source: https://www.avalara.com/vatlive/en/country-guides/south-america/argentina/argentina-vat-compliance-and-rates.html
							// Source: https://myargentinepassport.com/en/taxes-in-argentina-for-foreigners/
							Since:   cal.NewDate(1995, 4, 1),
							Percent: num.MakePercentage(105, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateSuperReduced,
					Name: i18n.String{
						i18n.EN: "Super-Reduced Rate",
						i18n.ES: "Alícuota Diferencial",
					},
					Values: []*tax.RateValueDef{
						{
							// Super-reduced VAT rate - 2.5%
							// Applied to specific categories such as capital goods
							// Reference: AFIP - Decree 493/2001
							Since:   cal.NewDate(2001, 5, 1),
							Percent: num.MakePercentage(25, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyZero},
					Rate: tax.RateZero,
					Name: i18n.String{
						i18n.EN: "Zero Rate",
						i18n.ES: "Tasa Cero",
					},
					Values: []*tax.RateValueDef{
						{
							// 0% rate for exports and specific categories
							// Reference: AFIP - Law 23.349 Article 43
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
			},
		},

		//
		// IVA Retenido (Retained VAT)
		//
		{
			Code:     TaxCategoryRIVA,
			Retained: true,
			Name: i18n.String{
				i18n.EN: "Retained VAT",
				i18n.ES: "IVA Retenido",
			},
			Title: i18n.String{
				i18n.EN: "Retained Value Added Tax",
				i18n.ES: "Impuesto al Valor Agregado Retenido",
			},
			Description: &i18n.String{
				i18n.EN: here.Doc(`
					VAT retention scheme where the purchasing party withholds a percentage of VAT
					from the payment and deposits it directly with AFIP. The retention percentage
					varies based on the supplier's tax classification and registration status.
				`),
			},
			// No fixed rates - varies by taxpayer category
			// Reference: AFIP - RG 2854/2010 and modifications
			// Source: https://www.afip.gob.ar/sire/percepciones-retenciones/
			Rates: []*tax.RateDef{},
		},

		//
		// Ganancias (Income Tax Withholding)
		//
		{
			Code:     TaxCategoryRG,
			Retained: true,
			Name: i18n.String{
				i18n.EN: "Income Tax Withholding",
				i18n.ES: "Retención de Ganancias",
			},
			Title: i18n.String{
				i18n.EN: "Withholding of Income Tax",
				i18n.ES: "Retención del Impuesto a las Ganancias",
			},
			Description: &i18n.String{
				i18n.EN: here.Doc(`
					Income tax withholding applied to payments for services and other taxable income.
					The withholding percentage varies based on the type of service and the recipient's
					tax classification. Common rates range from 0.5% to 35%.
				`),
			},
			// Variable rates based on service type and taxpayer category
			// Reference: AFIP - RG 830/2000, RG 4003/2017 and modifications
			// Source: https://www.afip.gob.ar/genericos/guiavirtual/directorio_subcategoria_nivel3.aspx?id_nivel1=563id_nivel2%3D607&id_nivel3=686
			// Source: https://servicioscf.afip.gob.ar/calc-rg830/
			Rates: []*tax.RateDef{},
		},

		//
		// Ingresos Brutos (Gross Income Tax - Provincial)
		//
		{
			Code:     TaxCategoryIB,
			Retained: false,
			Name: i18n.String{
				i18n.EN: "Gross Income Tax",
				i18n.ES: "Ingresos Brutos",
			},
			Title: i18n.String{
				i18n.EN: "Gross Income Tax",
				i18n.ES: "Impuesto sobre los Ingresos Brutos",
			},
			Description: &i18n.String{
				i18n.EN: here.Doc(`
					Provincial tax levied on gross income from economic activities. Each Argentine
					province sets its own rates and regulations. Common rates range from 1% to 5%
					depending on the province and type of activity.
				`),
			},
			// Provincial tax - rates vary by jurisdiction and activity
			// Reference: Provincial tax codes (each province has its own)
			Rates: []*tax.RateDef{},
		},
	}
}
