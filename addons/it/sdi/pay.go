package sdi

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyTreasury       cbc.Key = "treasury"
	MeansKeyPaymentSlip    cbc.Key = "payment-slip"
	MeansKeyRID            cbc.Key = "rid"
	MeansKeyRIDUtility     cbc.Key = "rid-utility"
	MeansKeyRIDFast        cbc.Key = "rid-fast"
	MeansKeyRIBA           cbc.Key = "riba"
	MeansKeyTaxReceipt     cbc.Key = "tax-receipt"
	MeansKeySpecialAccount cbc.Key = "special-account"
	MeansKeyPostOffice     cbc.Key = "post-office"
	MeansKeySEPA           cbc.Key = "sepa"
	MeansKeySEPACore       cbc.Key = "sepa-core"
	MeansKeySEPAB2B        cbc.Key = "sepa-b2b"
	MeansKeyPagoPA         cbc.Key = "pagopa"
)

// PaymentMeansExtensions returns the mapping of payment means to their
// extension values used by SDI.
func PaymentMeansExtensions() tax.Extensions {
	return paymentMeansKeyMap
}

var paymentMeansKeyMap = tax.Extensions{
	pay.MeansKeyCash:                                 "MP01",
	pay.MeansKeyCheque:                               "MP02",
	pay.MeansKeyBankDraft:                            "MP03",
	pay.MeansKeyCash.With(MeansKeyTreasury):          "MP04",
	pay.MeansKeyCreditTransfer:                       "MP05",
	pay.MeansKeyPromissoryNote:                       "MP06",
	pay.MeansKeyOther.With(MeansKeyPaymentSlip):      "MP07",
	pay.MeansKeyCard:                                 "MP08",
	pay.MeansKeyOnline:                               "MP08",
	pay.MeansKeyDirectDebit.With(MeansKeyRID):        "MP09",
	pay.MeansKeyDirectDebit.With(MeansKeyRIDUtility): "MP10",
	pay.MeansKeyDirectDebit.With(MeansKeyRIDFast):    "MP11",
	pay.MeansKeyDirectDebit.With(MeansKeyRIBA):       "MP12",
	pay.MeansKeyDebitTransfer:                        "MP13",
	pay.MeansKeyOther.With(MeansKeyTaxReceipt):       "MP14",
	pay.MeansKeyOther.With(MeansKeySpecialAccount):   "MP15",
	pay.MeansKeyDirectDebit:                          "MP16",
	pay.MeansKeyDirectDebit.With(MeansKeyPostOffice): "MP17",
	pay.MeansKeyCheque.With(MeansKeyPostOffice):      "MP18",
	pay.MeansKeyDirectDebit.With(MeansKeySEPA):       "MP19",
	pay.MeansKeyDirectDebit.With(MeansKeySEPACore):   "MP20",
	pay.MeansKeyDirectDebit.With(MeansKeySEPAB2B):    "MP21",
	pay.MeansKeyNetting:                              "MP22",
	pay.MeansKeyOnline.With(MeansKeyPagoPA):          "MP23",
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
