package nfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Addon-specific Payment Means Keys
const (
	MeansKeyDebit cbc.Key = "debit"
)

var paymentMeansKeyMap = tax.Extensions{
	pay.MeansKeyCash:   "01", // Dinheiro
	pay.MeansKeyCheque: "02", // Cheque
	pay.MeansKeyCard:   "03", // Cartão de Crédito
	pay.MeansKeyDebitTransfer.With(MeansKeyDebit): "04", // Cartão de Débito
	pay.MeansKeyCreditTransfer:                    "18", // Transferência bancária
	pay.MeansKeyOnline:                            "18", // Carteira Digital
	pay.MeansKeyOther:                             "99", // Outros
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil || (instr.Ext != nil && instr.Key == pay.MeansKeyOther) {
		return
	}
	if code := paymentMeansKeyMap[instr.Key]; code != "" {
		instr.Ext = instr.Ext.Merge(tax.Extensions{
			ExtKeyPaymentMeans: code,
		})
	}
}

func normalizePayAdvance(adv *pay.Advance) {
	if adv == nil || (adv.Ext != nil && adv.Key == pay.MeansKeyOther) {
		return
	}
	if code := paymentMeansKeyMap[adv.Key]; code != "" {
		adv.Ext = adv.Ext.Merge(tax.Extensions{
			ExtKeyPaymentMeans: code,
		})
	}
}

func validatePayInstructions(i *pay.Instructions) error {
	if i == nil {
		return nil
	}
	return validation.ValidateStruct(i,
		validation.Field(&i.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}

func validatePayAdvance(a *pay.Advance) error {
	if a == nil {
		return nil
	}
	return validation.ValidateStruct(a,
		validation.Field(&a.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}
