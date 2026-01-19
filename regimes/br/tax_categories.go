package br

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Tax categories specific for Brazil.
const (
	TaxCategoryISS       cbc.Code = "ISS"
	TaxCategoryICMS      cbc.Code = "ICMS"
	TaxCategoryIPI       cbc.Code = "IPI"
	TaxCategoryPIS       cbc.Code = "PIS"
	TaxCategoryPISRet    cbc.Code = "PISRet"
	TaxCategoryCOFINS    cbc.Code = "COFINS"
	TaxCategoryCOFINSRet cbc.Code = "COFINSRet"
	TaxCategoryCSLL      cbc.Code = "CSLL"
	TaxCategoryINSS      cbc.Code = "INSS"
	TaxCategoryIRRF      cbc.Code = "IRRF"
	TaxCategoryIBS       cbc.Code = "IBS"
	TaxCategoryCBS       cbc.Code = "CBS"
)

var taxCategories = []*tax.CategoryDef{
	//
	// Municipal Service Tax (ISS)
	//
	{
		Code: TaxCategoryISS,
		Name: i18n.String{
			i18n.EN: "ISS",
			i18n.PT: "ISS",
		},
		Title: i18n.String{
			i18n.EN: "Municipal Service Tax",
			i18n.PT: "Imposto Sobre Serviços",
		},
		Retained:    false,
		Informative: true,
	},
	//
	// State value-added tax (ICMS)
	//
	{
		Code: TaxCategoryICMS,
		Name: i18n.String{
			i18n.EN: "ICMS",
			i18n.PT: "ICMS",
		},
		Title: i18n.String{
			i18n.EN: "State value-added tax",
			i18n.PT: "Imposto sobre Circulação de Mercadorias e Serviços",
		},
		Retained:    false,
		Informative: true,
	},
	//
	// Federal value-added Tax (IPI)
	//
	{
		Code: TaxCategoryIPI,
		Name: i18n.String{
			i18n.EN: "IPI",
			i18n.PT: "IPI",
		},
		Title: i18n.String{
			i18n.EN: "Federal value-added Tax",
			i18n.PT: "Imposto sobre Produtos Industrializados",
		},
		Retained: false,
	},
	//
	// Social Integration Program (PIS)
	//
	{
		Code: TaxCategoryPIS,
		Name: i18n.String{
			i18n.EN: "PIS",
			i18n.PT: "PIS",
		},
		Title: i18n.String{
			i18n.EN: "Social Integration Program",
			i18n.PT: "Programa de Integração Social",
		},
		Retained:    false,
		Informative: true,
	},
	{
		Code: TaxCategoryPISRet,
		Name: i18n.String{
			i18n.EN: "PIS (Retained)",
			i18n.PT: "PIS (Retido)",
		},
		Title: i18n.String{
			i18n.EN: "Social Integration Program Retained",
			i18n.PT: "Programa de Integração Social Retido",
		},
		Retained: true,
	},
	//
	// Contribution for the Financing of Social Security (COFINS)
	//
	{
		Code: TaxCategoryCOFINS,
		Name: i18n.String{
			i18n.EN: "COFINS",
			i18n.PT: "COFINS",
		},
		Title: i18n.String{
			i18n.EN: "Contribution for the Financing of Social Security",
			i18n.PT: "Contribuição para o Financiamento da Seguridade Social",
		},
		Informative: true,
	},
	{
		Code: TaxCategoryCOFINSRet,
		Name: i18n.String{
			i18n.EN: "COFINS (Retained)",
			i18n.PT: "COFINS (Retido)",
		},
		Title: i18n.String{
			i18n.EN: "Contribution for the Financing of Social Security Retained",
			i18n.PT: "Contribuição para o Financiamento da Seguridade Social Retido",
		},
		Retained: true,
	},
	//
	// Social Contribution on Net Profit (CSLL)
	//
	{
		Code: TaxCategoryCSLL,
		Name: i18n.String{
			i18n.EN: "CSLL",
			i18n.PT: "CSLL",
		},
		Title: i18n.String{
			i18n.EN: "Social Contribution on Net Profit",
			i18n.PT: "Contribuição Social sobre o Lucro Líquido",
		},
		Retained: true,
	},
	//
	// National Social Security Institute (INSS)
	//
	{
		Code: TaxCategoryINSS,
		Name: i18n.String{
			i18n.EN: "INSS",
			i18n.PT: "INSS",
		},
		Title: i18n.String{
			i18n.EN: "National Social Security Institute Withholding",
			i18n.PT: "Retenção do Instituto Nacional do Seguro Social",
		},
		Retained: true,
	},
	//
	// Income Tax Withheld at Source (IRRF)
	//
	{
		Code: TaxCategoryIRRF,
		Name: i18n.String{
			i18n.EN: "IRRF",
			i18n.PT: "IRRF",
		},
		Title: i18n.String{
			i18n.EN: "Income Tax Withheld at Source",
			i18n.PT: "Imposto de Renda Retido na Fonte",
		},
		Retained: true,
	},
	//
	// Goods and Services Tax (IBS)
	//
	{
		Code: TaxCategoryIBS,
		Name: i18n.String{
			i18n.EN: "IBS",
			i18n.PT: "IBS",
		},
		Title: i18n.String{
			i18n.EN: "Goods and Services Tax",
			i18n.PT: "Imposto sobre Bens e Serviços",
		},
		// IBS is not really informative, but since during 2026 reporting is
		// mandatory but collection is not, we make it informative to avoid
		// increasing the tax totals.
		Informative: true,
	},
	//
	// Federal Goods and Services Contribution (CBS)
	//
	{
		Code: TaxCategoryCBS,
		Name: i18n.String{
			i18n.EN: "CBS",
			i18n.PT: "CBS",
		},
		Title: i18n.String{
			i18n.EN: "Federal Goods and Services Contribution",
			i18n.PT: "Contribuição sobre Bens e Serviços",
		},
		// CBS is not really informative, but since during 2026 reporting is
		// mandatory but collection is not, we make it informative to avoid
		// increasing the tax totals.
		Informative: true,
	},
}
