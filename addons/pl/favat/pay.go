package favat

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyVoucher cbc.Key = "voucher"
)

var paymentMeansKeyMap = cbc.CodeMap{
	pay.MeansKeyCash:                           "1", // Cash / Gotówka
	pay.MeansKeyCard:                           "2", // Card / Karta
	pay.MeansKeyOther.With(MeansKeyVoucher):    "3", // Voucher / Bon
	pay.MeansKeyCheque:                         "4", // Cheque / Czek
	pay.MeansKeyOther.With(pay.MeansKeyCredit): "5", // Credit / Kredyt
	pay.MeansKeyCreditTransfer:                 "6", // Credit Transfer / Przelew
	pay.MeansKeyOnline:                         "7", // Online / Mobilna
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if code := paymentMeansKeyMap.Lookup(instr.Key); code != "" {
		instr.Ext = instr.Ext.Merge(tax.ExtensionsOf(cbc.CodeMap{
			ExtKeyPaymentMeans: code,
		}))
	}
}

func normalizePayAdvance(adv *pay.Advance) {
	if adv == nil {
		return
	}
	if code := paymentMeansKeyMap.Lookup(adv.Key); code != "" {
		adv.Ext = adv.Ext.Merge(tax.ExtensionsOf(cbc.CodeMap{
			ExtKeyPaymentMeans: code,
		}))
	}
}

func payAdvanceRules() *rules.Set {
	return rules.For(new(pay.Advance),
		rules.Field("date",
			rules.Assert("01", "advance payment date is required", is.Present),
		),
	)
}
