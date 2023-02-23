package co

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Preceding document correction method constants.
const (
	CorrectionMethodKeyPartial cbc.Key = "partial"
	CorrectionMethodKeyRevoked cbc.Key = "revoked"
)

var correctionMethodList = []*tax.KeyDefinition{
	{
		Key:  CorrectionMethodKeyPartial,
		Code: "1",
		Desc: i18n.String{
			i18n.EN: "Partial refund of part of the goods or services.",
			i18n.ES: "Devolución de parte de los bienes; no aceptación de partes del servicio.",
		},
	},
	{
		Key:  CorrectionMethodKeyRevoked,
		Code: "2",
		Desc: i18n.String{
			i18n.EN: "Previous document has been completely cancelled.",
			i18n.ES: "Anulación de la factura anterior.",
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
