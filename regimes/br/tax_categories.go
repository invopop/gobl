package br

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Tax categories specific for Brazil.
const (
	TaxCategoryISS    cbc.Code = "ISS"
	TaxCategoryICMS   cbc.Code = "ICMS"
	TaxCategoryIPI    cbc.Code = "IPI"
	TaxCategoryPIS    cbc.Code = "PIS"
	TaxCategoryCOFINS cbc.Code = "COFINS"
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
		Retained: false,
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
		Retained: false,
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
		Retained: true,
	},
}
