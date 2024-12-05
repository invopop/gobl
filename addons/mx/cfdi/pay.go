package cfdi

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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

// PaymentMeansExtensions returns the mapping of payment means to their
// extension values used by CFDI.
func PaymentMeansExtensions() tax.Extensions {
	return paymentMeansKeyMap
}

var paymentMeansKeyMap = tax.Extensions{
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

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	extVal := paymentMeansKeyMap[instr.Key]
	if extVal != "" {
		if instr.Ext == nil {
			instr.Ext = make(tax.Extensions)
		}
		instr.Ext[ExtKeyPaymentMeans] = extVal
	}
}

func normalizePayAdvance(adv *pay.Advance) {
	if adv == nil {
		return
	}
	extVal := paymentMeansKeyMap[adv.Key]
	if extVal != "" {
		if adv.Ext == nil {
			adv.Ext = make(tax.Extensions)
		}
		adv.Ext[ExtKeyPaymentMeans] = extVal
	}
}

func validatePayAdvance(a *pay.Advance) error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}

func validatePayInstructions(i *pay.Instructions) error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}

func validatePayTerms(terms *pay.Terms) error {
	if terms == nil {
		return nil
	}
	return validation.ValidateStruct(terms,
		validation.Field(&terms.Notes,
			validation.Length(0, 1000),
			validation.Skip,
		),
	)
}
