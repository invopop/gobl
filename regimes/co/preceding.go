package co

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Preceding document correction method constants.
const (
	CorrectionMethodKeyPartial         cbc.Key = "partial"
	CorrectionMethodKeyRevoked         cbc.Key = "revoked"
	CorrectionMethodKeyDiscount        cbc.Key = "discount"
	CorrectionMethodKeyPriceAdjustment cbc.Key = "price-adjustment"
	CorrectionMethodKeyOther           cbc.Key = "other"
)

var correctionMethodList = []*tax.KeyDefinition{
	{
		Key: CorrectionMethodKeyPartial,
		Desc: i18n.String{
			i18n.EN: "Partial refund of part of the goods or services.",
			i18n.ES: "Devolución de parte de los bienes; no aceptación de partes del servicio.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "1",
		},
	},
	{
		Key: CorrectionMethodKeyRevoked,
		Desc: i18n.String{
			i18n.EN: "Previous document has been completely cancelled.",
			i18n.ES: "Anulación de la factura anterior.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "2",
		},
	},
	{
		Key: CorrectionMethodKeyDiscount,
		Desc: i18n.String{
			i18n.EN: "Partial or total discount.",
			i18n.ES: "Rebaja o descuento parcial o total.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "3",
		},
	},
	{
		Key: CorrectionMethodKeyPriceAdjustment,
		Desc: i18n.String{
			i18n.EN: "Ajuste de precio.",
			i18n.ES: "Price adjustment.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "4",
		},
	},
	{
		Key: CorrectionMethodKeyOther,
		Desc: i18n.String{
			i18n.EN: "Otros.",
			i18n.ES: "Other.",
		},
		Map: cbc.CodeMap{
			KeyDIAN: "5",
		},
	},
}

func correctionMethodKeys() []interface{} {
	keys := make([]interface{}, len(correctionMethodList))
	for i, v := range correctionMethodList {
		keys[i] = v.Key
	}
	return keys
}

var isValidCorrectionMethodKey = validation.In(correctionMethodKeys()...)
