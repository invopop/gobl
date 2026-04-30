package nfe

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// Addon-specific Payment Means Keys
const (
	MeansKeyDebit cbc.Key = "debit"
)

var paymentMeansKeyMap = map[cbc.Key]cbc.Code{
	pay.MeansKeyCash:   "01", // Dinheiro
	pay.MeansKeyCheque: "02", // Cheque
	pay.MeansKeyCard:   "03", // Cartão de Crédito
	pay.MeansKeyDebitTransfer.With(MeansKeyDebit): "04", // Cartão de Débito
	pay.MeansKeyCreditTransfer:                    "18", // Transferência bancária
	pay.MeansKeyOnline:                            "18", // Carteira Digital
	pay.MeansKeyOther:                             "99", // Outros
}

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if instr.Ext.Has(ExtKeyPaymentMeans) && instr.Key == pay.MeansKeyOther {
		// `other` key does not override the extension
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
	if adv.Ext.Has(ExtKeyPaymentMeans) && adv.Key == pay.MeansKeyOther {
		// `other` key does not override the extension already set
		return
	}
	if code := paymentMeansKeyMap[adv.Key]; code != "" {
		adv.Ext = adv.Ext.Merge(tax.ExtensionsOf(tax.ExtMap{
			ExtKeyPaymentMeans: code,
		}))
	}
}

func payInstructionsRules() *rules.Set {
	return rules.For(new(pay.Instructions),
		rules.Field("ext",
			rules.Assert("01", fmt.Sprintf("payment instructions require '%s' extension", ExtKeyPaymentMeans),
				tax.ExtensionsRequire(ExtKeyPaymentMeans),
			),
		),
	)
}

func payAdvanceRules() *rules.Set {
	return rules.For(new(pay.Record),
		rules.Field("ext",
			rules.Assert("01", fmt.Sprintf("payment advance requires '%s' extension", ExtKeyPaymentMeans),
				tax.ExtensionsRequire(ExtKeyPaymentMeans),
			),
		),
	)
}
