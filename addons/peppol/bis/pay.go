package bis

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// directDebitMeansCodes are the UNTDID 4461 codes that identify direct-debit
// payment methods and therefore require a mandate reference (PEPPOL-EN16931-R061).
var directDebitMeansCodes = []string{"49", "59"}

func payInstructionsRules() *rules.Set {
	return rules.For(new(pay.Instructions),
		// PEPPOL-EN16931-R061: mandate reference required for direct debit.
		rules.Assert("R061", "mandate reference is required for direct debit payments (PEPPOL-EN16931-R061)",
			is.Func("direct debit mandate", directDebitMandatePresent),
		),
	)
}

func directDebitMandatePresent(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans).String()
	isDirectDebit := false
	for _, c := range directDebitMeansCodes {
		if c == code {
			isDirectDebit = true
			break
		}
	}
	if !isDirectDebit {
		return true
	}
	return instr.DirectDebit != nil && instr.DirectDebit.Ref != ""
}
