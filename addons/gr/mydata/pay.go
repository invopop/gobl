package mydata

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
	MeansKeyForeign cbc.Key = "foreign"
)

// PaymentMeansExtensions returns the mapping of payment means to their
// extension values used by myDATA.
func PaymentMeansExtensions() tax.Extensions {
	return tax.ExtensionsOf(paymentMeansMap)
}

var paymentMeansMap = map[cbc.Key]cbc.Code{
	pay.MeansKeyCreditTransfer:                       "1",
	pay.MeansKeyCreditTransfer.With(MeansKeyForeign): "2",
	pay.MeansKeyCash:                                 "3",
	pay.MeansKeyCheque:                               "4",
	pay.MeansKeyPromissoryNote:                       "5",
	pay.MeansKeyOnline:                               "6",
	pay.MeansKeyCard:                                 "7",
}

func normalizePayInstructions(i *pay.Instructions) {
	if i == nil {
		return
	}
	extVal := paymentMeansMap[i.Key]
	if extVal != "" {
		if i.Ext.IsZero() {
			i.Ext = tax.MakeExtensions()
		}
		i.Ext = i.Ext.Set(ExtKeyPaymentMeans, extVal)
	}
}

func normalizePayAdvance(a *pay.Advance) {
	if a == nil {
		return
	}
	extVal := paymentMeansMap[a.Key]
	if extVal != "" {
		if a.Ext.IsZero() {
			a.Ext = tax.MakeExtensions()
		}
		a.Ext = a.Ext.Set(ExtKeyPaymentMeans, extVal)
	}
}

func payInstructionsRules() *rules.Set {
	return rules.For(new(pay.Instructions),
		rules.Field("key",
			rules.Assert("01", "payment instructions key is required", is.Present),
		),
		rules.Field("ext",
			rules.Assert("02",
				fmt.Sprintf("payment instructions require '%s' extension", ExtKeyPaymentMeans),
				tax.ExtensionsRequire(ExtKeyPaymentMeans),
			),
		),
	)
}

func payAdvanceRules() *rules.Set {
	return rules.For(new(pay.Advance),
		rules.Field("key",
			rules.Assert("01", "payment advance key is required", is.Present),
		),
		rules.Field("ext",
			rules.Assert("02",
				fmt.Sprintf("payment advance requires '%s' extension", ExtKeyPaymentMeans),
				tax.ExtensionsRequire(ExtKeyPaymentMeans),
			),
		),
	)
}
