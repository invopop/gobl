package co

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Preceding document correction method constants.
const (
	CorrectionKeyPartial         cbc.Key = "partial"
	CorrectionKeyRevoked         cbc.Key = "revoked"
	CorrectionKeyDiscount        cbc.Key = "discount"
	CorrectionKeyPriceAdjustment cbc.Key = "price-adjustment"
	CorrectionKeyOther           cbc.Key = "other"
)

var correctionList = []*tax.KeyDefinition{
	{
		Key: CorrectionKeyPartial,
		Name: i18n.String{
			i18n.EN: "Partial refund",
			i18n.ES: "Devolución parcial",
		},
		Desc: i18n.String{
			i18n.EN: "Partial refund of part of the goods or services.",
			i18n.ES: "Devolución de parte de los bienes; no aceptación de partes del servicio.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "1",
		},
	},
	{
		Key: CorrectionKeyRevoked,
		Name: i18n.String{
			i18n.EN: "Revoked",
			i18n.ES: "Anulación",
		},
		Desc: i18n.String{
			i18n.EN: "Previous document has been completely cancelled.",
			i18n.ES: "Anulación de la factura anterior.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "2",
		},
	},
	{
		Key: CorrectionKeyDiscount,
		Name: i18n.String{
			i18n.EN: "Discount",
			i18n.ES: "Descuento",
		},
		Desc: i18n.String{
			i18n.EN: "Partial or total discount.",
			i18n.ES: "Rebaja o descuento parcial o total.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "3",
		},
	},
	{
		Key: CorrectionKeyPriceAdjustment,
		Name: i18n.String{
			i18n.EN: "Adjustment",
			i18n.ES: "Ajuste",
		},
		Desc: i18n.String{
			i18n.EN: "Price adjustment.",
			i18n.ES: "Ajuste de precio.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "4",
		},
	},
	{
		Key: CorrectionKeyOther,
		Name: i18n.String{
			i18n.EN: "Other",
			i18n.ES: "Otros",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "5",
		},
	},
}

func correctionKeys() []interface{} {
	keys := make([]interface{}, len(correctionList))
	for i, v := range correctionList {
		keys[i] = v.Key
	}
	return keys
}

var isValidCorrectionKey = validation.In(correctionKeys()...)
