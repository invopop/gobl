package pa

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Panama tax categories
const (
	TaxCategoryISC cbc.Code = "ISC"
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
		Keys: tax.GlobalVATKeys(),
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "DGI - SFEP Technical Documentation",
				},
				URL: "https://dgi.mef.gob.pa/_7facturaelectronica/source/F-T%C3%A9cnica%20de%20Factura%20Electr%C3%B3nica%20para%20los%20Proveedores%20de%20Autorizaci%C3%B3n%20Calificados%20V1.00-Abril2025.pdf",
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
	},
}
