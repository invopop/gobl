package en16931

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var paymentMeansMap = tax.Extensions{
	pay.MeansKeyAny:                                   "1",
	pay.MeansKeyCard:                                  "48",
	pay.MeansKeyCreditTransfer:                        "30",
	pay.MeansKeyDebitTransfer:                         "31",
	pay.MeansKeyCash:                                  "10",
	pay.MeansKeyCheque:                                "20",
	pay.MeansKeyBankDraft:                             "21",
	pay.MeansKeyDirectDebit:                           "49",
	pay.MeansKeyOnline:                                "68",
	pay.MeansKeyPromissoryNote:                        "60",
	pay.MeansKeyNetting:                               "97",
	pay.MeansKeyCreditTransfer.With(pay.MeansKeySEPA): "58",
	pay.MeansKeyDirectDebit.With(pay.MeansKeySEPA):    "59",
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if val, ok := paymentMeansMap[instr.Key]; ok {
		instr.Ext = instr.Ext.Merge(
			tax.Extensions{untdid.ExtKeyPaymentMeans: val},
		)
	}
}

func validatePayInstructions(instr *pay.Instructions) error {
	return validation.ValidateStruct(instr,
		validation.Field(&instr.Ext,
			// BR-49
			tax.ExtensionsRequire(untdid.ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}
