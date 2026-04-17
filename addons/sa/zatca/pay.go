package zatca

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if val, ok := en16931.PaymentMeansMap[instr.Key]; ok {
		instr.Ext = instr.Ext.Merge(
			tax.Extensions{untdid.ExtKeyPaymentMeans: val},
		)
	}
}

func payInstructionsRules() *rules.Set {
	return rules.For(new(pay.Instructions),
		rules.Field("ext",
			rules.Assert("01", "payment means extension is required (BR-49)",
				tax.ExtensionsRequire(untdid.ExtKeyPaymentMeans),
			),
		),
	)
}
