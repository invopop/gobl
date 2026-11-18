package en16931

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

var paymentMeansMap = cbc.CodeMap{
	pay.MeansKeyAny:                                   "1",
	pay.MeansKeyCard:                                  "48",
	pay.MeansKeyCard.With(pay.MeansKeyDebit):          "55",
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
	if val := paymentMeansMap.Lookup(instr.Key); val != "" {
		instr.Ext = instr.Ext.Merge(
			tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: val}),
		)
	}
}

func payInstructionsRules() *rules.Set {
	return rules.For(new(pay.Instructions),
		rules.Field("ext",
			rules.Assert("01", "payment means extension is required (BR-49)",
				tax.ExtensionsRequire(untdid.ExtKeyPaymentMeans),
			),
		),
	)
}
