package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyWallet          cbc.Key = "wallet"
	MeansKeyGroceryVouchers cbc.Key = "grocery-vouchers"
	MeansKeyInKind          cbc.Key = "in-kind"
	MeansKeySubrogation     cbc.Key = "subrogation"
	MeansKeyConsignment     cbc.Key = "consignment"
	MeansKeyDebtRelief      cbc.Key = "debt-relief"
	MeansKeyNovation        cbc.Key = "novation"
	MeansKeyMerger          cbc.Key = "merger"
	MeansKeyRemission       cbc.Key = "remission"
	MeansKeyExpiration      cbc.Key = "expiration"
	MeansKeySatisfyCreditor cbc.Key = "satisfy-creditor"
	MeansKeyDebit           cbc.Key = "debit"
	MeansKeyServices        cbc.Key = "services"
	MeansKeyAdvance         cbc.Key = "advance"
	MeansKeyIntermediary    cbc.Key = "intermediary"
)

var paymentMeansKeyMap = map[cbc.Key]tax.ExtValue{
	pay.MeansKeyCash:                                "01",
	pay.MeansKeyCheque:                              "02",
	pay.MeansKeyCreditTransfer:                      "03",
	pay.MeansKeyCard:                                "04",
	pay.MeansKeyOnline.With(MeansKeyWallet):         "05",
	pay.MeansKeyOnline:                              "06",
	pay.MeansKeyOther.With(MeansKeyGroceryVouchers): "08",
	pay.MeansKeyOther.With(MeansKeyInKind):          "12",
	pay.MeansKeyOther.With(MeansKeySubrogation):     "13",
	pay.MeansKeyOther.With(MeansKeyConsignment):     "14",
	pay.MeansKeyOther.With(MeansKeyDebtRelief):      "15",
	pay.MeansKeyNetting:                             "17",
	pay.MeansKeyOther.With(MeansKeyNovation):        "23",
	pay.MeansKeyOther.With(MeansKeyMerger):          "24",
	pay.MeansKeyOther.With(MeansKeyRemission):       "25",
	pay.MeansKeyOther.With(MeansKeyExpiration):      "26",
	pay.MeansKeyOther.With(MeansKeySatisfyCreditor): "27",
	pay.MeansKeyOther.With(MeansKeyDebit):           "28",
	pay.MeansKeyOther.With(MeansKeyServices):        "29",
	pay.MeansKeyOther.With(MeansKeyAdvance):         "30",
	pay.MeansKeyOther.With(MeansKeyIntermediary):    "31",
}

func normalizeInvoicePaymentInstructions(inv *bill.Invoice) {
	if inv.Payment == nil || inv.Payment.Instructions == nil {
		return
	}
	instr := inv.Payment.Instructions
	extVal := paymentMeansKeyMap[instr.Key]
	if extVal != "" {
		if instr.Ext == nil {
			instr.Ext = make(tax.Extensions)
		}
		instr.Ext[ExtKeyPaymentMeans] = extVal
	}
}

func normalizeInvoicePaymentAdvances(inv *bill.Invoice) {
	if inv.Payment == nil || len(inv.Payment.Advances) == 0 {
		return
	}

	for _, adv := range inv.Payment.Advances {
		extVal := paymentMeansKeyMap[adv.Key]
		if extVal != "" {
			if adv.Ext == nil {
				adv.Ext = make(tax.Extensions)
			}
			adv.Ext[ExtKeyPaymentMeans] = extVal
		}
	}
}
