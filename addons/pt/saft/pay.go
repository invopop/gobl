package saft

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// PaymentMeansExtensions returns the mapping of payment means to their
// extension values used by SAF-T PT.
func PaymentMeansExtensions() tax.Extensions {
	return paymentMeansMap
}

var paymentMeansMap = tax.Extensions{
	pay.MeansKeyCard:           "CC",
	pay.MeansKeyCreditTransfer: "TB",
	pay.MeansKeyDebitTransfer:  "TB",
	pay.MeansKeyCash:           "NU",
	pay.MeansKeyPromissoryNote: "LC",
	pay.MeansKeyNetting:        "CS",
	pay.MeansKeyCheque:         "CH",
	pay.MeansKeyDirectDebit:    "TB",
	pay.MeansKeyOnline:         "DE",
	pay.MeansKeyOther:          "OU",
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	instr.Ext = mergePaymentMeans(instr.Key, instr.Ext)
}

func normalizePayAdvance(adv *pay.Advance) {
	if adv == nil {
		return
	}
	adv.Ext = mergePaymentMeans(adv.Key, adv.Ext)
}

func mergePaymentMeans(key cbc.Key, ext tax.Extensions) tax.Extensions {
	if key == pay.MeansKeyOther && ext[ExtKeyPaymentMeans] != "" {
		return ext // `other` won't override the extension if already set
	}
	if extVal, ok := paymentMeansMap[key]; ok {
		return ext.Merge(
			tax.Extensions{ExtKeyPaymentMeans: extVal},
		)
	}
	return ext
}
