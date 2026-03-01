package pa

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

const (
	TaxCategoryISC cbc.Code = "ISC"
)

const (
	RateIncreased cbc.Key = "increased"
)

var taxCategories = []*tax.CategoryDef{
	//
	// ITBMS
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "ITBMS",
			i18n.ES: "ITBMS",
		},
		Title: i18n.String{
			i18n.EN: "Transfer of Tangible Movable Goods and Services Tax",
			i18n.ES: "Impuesto de Transferencia de Bienes Muebles y Servicios",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Tasa General",
				},
			Values: []*tax.RateValueDef{
				{
					Percent: num.MakePercentage(70, 3),
					Since:   cal.NewDate(2010, 3, 15),
				},
				{
					Percent: num.MakePercentage(50, 3),
					Since:   cal.NewDate(1977, 3, 1),
				},
			},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: RateIncreased,
				Name: i18n.String{
					i18n.EN: "Increased Rate (Hotels, Alcoholic Beverages)",
					i18n.ES: "Tasa Incrementada (Hoteles, Bebidas Alcohólicas)",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(100, 3),
						Since:   cal.NewDate(2010, 3, 15),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Special Rate (Tobacco)",
					i18n.ES: "Tasa Especial (Tabaco)",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(150, 3),
						Since:   cal.NewDate(2010, 3, 15),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "DGI - SFEP Technical Documentation",
				},
				URL: "https://dgi.mef.gob.pa/FacturaElectronica/Documentacion.html",
			},
		},
	},

	//
	// ISC (Impuesto Selectivo al Consumo)
	//
	{
		Code: TaxCategoryISC,
		Name: i18n.String{
			i18n.EN: "ISC",
			i18n.ES: "ISC",
		},
		Title: i18n.String{
			i18n.EN: "Selective Consumption Tax",
			i18n.ES: "Impuesto Selectivo al Consumo",
		},
		Retained: false,
		Rates:    []*tax.RateDef{},
	},
}
