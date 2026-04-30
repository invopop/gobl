package cfdi

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
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
	MeansKeyServices        cbc.Key = "services"
	MeansKeyAdvance         cbc.Key = "advance"
	MeansKeyIntermediary    cbc.Key = "intermediary"
)

// PaymentMeansExtensions returns the mapping of payment means to their
// extension values used by CFDI.
func PaymentMeansExtensions() tax.Extensions {
	return tax.ExtensionsOf(tax.ExtMap(paymentMeansKeyMap))
}

var paymentMeansKeyMap = cbc.CodeMap{
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
	pay.MeansKeyCard.With(pay.MeansKeyDebit):        "28",
	pay.MeansKeyOther.With(pay.MeansKeyDebit):       "28", // deprecated
	pay.MeansKeyOther.With(MeansKeyServices):        "29",
	pay.MeansKeyOther.With(MeansKeyAdvance):         "30",
	pay.MeansKeyOther.With(MeansKeyIntermediary):    "31",
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if code := paymentMeansKeyMap.Lookup(instr.Key); code != "" {
		instr.Ext = instr.Ext.Merge(tax.ExtensionsOf(tax.ExtMap{
			ExtKeyPaymentMeans: code,
		}))
	}
}

func normalizePayAdvance(adv *pay.Advance) {
	if adv == nil {
		return
	}
	if code := paymentMeansKeyMap.Lookup(adv.Key); code != "" {
		adv.Ext = adv.Ext.Merge(tax.ExtensionsOf(tax.ExtMap{
			ExtKeyPaymentMeans: code,
		}))
	}
}

func payInstructionsRules() *rules.Set {
	return rules.For(new(pay.Instructions),
		rules.Field("ext",
			rules.Assert("01",
				fmt.Sprintf("payment instructions require '%s' extension", ExtKeyPaymentMeans),
				tax.ExtensionsRequire(ExtKeyPaymentMeans),
			),
		),
	)
}

func payAdvanceRules() *rules.Set {
	return rules.For(new(pay.Advance),
		rules.Field("ext",
			rules.Assert("01",
				fmt.Sprintf("payment advance requires '%s' extension", ExtKeyPaymentMeans),
				tax.ExtensionsRequire(ExtKeyPaymentMeans),
			),
		),
	)
}

func payTermsRules() *rules.Set {
	return rules.For(new(pay.Terms),
		rules.Field("notes",
			rules.Assert("01", "notes length must be no more than 1000", is.Length(0, 1000)),
		),
	)
}
