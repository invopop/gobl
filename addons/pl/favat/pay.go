package favat

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyCoupon     cbc.Key = "coupon"
	MeansKeyCheque     cbc.Key = "cheque"
	MeansKeyLoan       cbc.Key = "loan"
	MeansKeyDebtRelief cbc.Key = "credit-transfer"
	MeansKeyMobile     cbc.Key = "mobile"
)

var paymentMeansKeyMap = tax.Extensions{
	pay.MeansKeyCash:                       "1", // Cash / Gotówka
	pay.MeansKeyCard:                       "2", // Card / Karta
	pay.MeansKeyOther.With(MeansKeyCoupon): "3", // Coupon / Bon
	pay.MeansKeyCheque:                     "4", // Cheque / Czek
	pay.MeansKeyOnline.With(MeansKeyLoan):  "5", // Loan / Kredyt
	pay.MeansKeyCreditTransfer:             "6", // Wire Transfer / Przelew
	pay.MeansKeyOther.With(MeansKeyMobile): "7", // Mobile / Mobilna
}

// PaymentMeansExtensions returns the mapping of payment means to their
// extension values used by FA_VAT.
func PaymentMeansExtensions() tax.Extensions {
	return paymentMeansKeyMap
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if code := paymentMeansKeyMap[instr.Key]; code != "" {
		instr.Ext = instr.Ext.Merge(tax.Extensions{
			ExtKeyPaymentMeans: code,
		})
	}
}

func normalizePayAdvance(adv *pay.Advance) {
	if adv == nil {
		return
	}
	if code := paymentMeansKeyMap[adv.Key]; code != "" {
		adv.Ext = adv.Ext.Merge(tax.Extensions{
			ExtKeyPaymentMeans: code,
		})
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
