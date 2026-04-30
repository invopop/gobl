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
	MeansKeyCredit  cbc.Key = "credit"
)

var paymentMeansKeyMap = map[cbc.Key]cbc.Code{
	pay.MeansKeyCash:                        "1", // Cash / Gotówka
	pay.MeansKeyCard:                        "2", // Card / Karta
	pay.MeansKeyOther.With(MeansKeyVoucher): "3", // Voucher / Bon
	pay.MeansKeyCheque:                      "4", // Cheque / Czek
	pay.MeansKeyOther.With(MeansKeyCredit):  "5", // Credit / Kredyt
	pay.MeansKeyCreditTransfer:              "6", // Credit Transfer / Przelew
	pay.MeansKeyOnline:                      "7", // Online / Mobilna
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if code := paymentMeansKeyMap[instr.Key]; code != "" {
		instr.Ext = instr.Ext.Merge(tax.ExtensionsOf(tax.ExtMap{
			ExtKeyPaymentMeans: code,
		}))
	}
}

func normalizePayRecord(adv *pay.Record) {
	if adv == nil {
		return
	}
	if code := paymentMeansKeyMap[adv.Key]; code != "" {
		adv.Ext = adv.Ext.Merge(tax.ExtensionsOf(tax.ExtMap{
			ExtKeyPaymentMeans: code,
		}))
	}
}

func payAdvanceRules() *rules.Set {
	return rules.For(new(pay.Record),
		rules.Field("date",
			rules.Assert("01", "advance payment date is required", is.Present),
		),
	)
}
