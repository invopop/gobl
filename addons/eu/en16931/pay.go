package en16931

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var paymentMeansMap = map[cbc.Key]cbc.Code{
	pay.MeansKeyAny:  "1",
	pay.MeansKeyCard: "48",
	pay.MeansKeyCard.With(pay.MeansKeyCredit):         "48",
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
	if val, ok := paymentMeansMap[instr.Key]; ok {
		instr.Ext = instr.Ext.Merge(
			tax.ExtensionsOf(tax.ExtMap{untdid.ExtKeyPaymentMeans: val}),
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

func payTermsRules() *rules.Set {
	return rules.For(new(pay.Terms),
		rules.Assert("01", "either due_dates or notes must be provided (BR-CO-25)",
			is.Func("has due dates or notes", payTermsHasDueDatesOrNotes),
		),
	)
}

func payTermsHasDueDatesOrNotes(val any) bool {
	t, ok := val.(*pay.Terms)
	return !ok || t == nil || len(t.DueDates) > 0 || t.Notes != ""
}
