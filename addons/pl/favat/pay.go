package favat

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyVoucher cbc.Key = "voucher"
	MeansKeyCredit  cbc.Key = "credit"
)

var paymentMeansKeyMap = tax.Extensions{
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

func validatePayAdvance(adv *pay.Advance) error {
	if adv == nil {
		return nil
	}
	return validation.ValidateStruct(adv,
		validation.Field(&adv.Date,
			validation.Required,
			validation.Skip,
		),
	)
}
