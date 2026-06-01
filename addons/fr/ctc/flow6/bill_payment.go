package flow6

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// normalizePayment surfaces the CDAR ProcessConditionCode for the
// payment on the fr-ctc-flow6-status extension and defaults the roles
// on the payment's parties — mirrors what normalizeStatus does for
// bill.Status. Advice payments are issued by the payer (BY → SE);
// receipt payments by the payee (SE → BY).
func normalizePayment(pmt *bill.Payment) {
	if pmt == nil {
		return
	}
	switch pmt.Type {
	case bill.PaymentTypeAdvice:
		pmt.Ext = pmt.Ext.Set(ExtKeyStatus, "211")
		setPartyRoleDefault(pmt.Customer, RoleSeller)
		setPartyRoleDefault(pmt.Supplier, RoleBuyer)
		// Default characteristic for an advice (211 Paiement transmis)
		// is "amount paid" (MPA). Callers can override to MEN or RAP.
		pmt.Ext = pmt.Ext.SetOneOf(ExtKeyCondition,
			ConditionAmountPaid, ConditionAmountReceived, ConditionAmountRemaining,
		)
	case bill.PaymentTypeReceipt:
		pmt.Ext = pmt.Ext.Set(ExtKeyStatus, "212")
		setPartyRoleDefault(pmt.Supplier, RoleSeller)
		setPartyRoleDefault(pmt.Customer, RoleBuyer)
		// Default characteristic for a receipt (212 Encaissée) is
		// "amount received" (MEN). Callers can override to MPA or RAP.
		pmt.Ext = pmt.Ext.SetOneOf(ExtKeyCondition,
			ConditionAmountReceived, ConditionAmountPaid, ConditionAmountRemaining,
		)
	}
}

// billPaymentRules validates the integrity of the addon's own extensions
// and the supported payment shape. French CTC format/business rules are
// the converter's responsibility — see the package doc.
func billPaymentRules() *rules.Set {
	return rules.For(new(bill.Payment),
		rules.Field("type",
			rules.Assert("01", "payment type must be 'advice' (CDAR 211) or 'receipt' (CDAR 212); 'request' is not a Flow 6 CDV event",
				is.In(bill.PaymentTypeAdvice, bill.PaymentTypeReceipt),
			),
		),
		rules.Field("ext",
			rules.Assert("02", "payment ext fr-ctc-flow6-status must be a Payment-applicable ProcessConditionCode (211 advice or 212 receipt); codes 200-210, 213 belong on bill.Status",
				tax.ExtensionsHasCodes(ExtKeyStatus, paymentProcessCodes...),
			),
			rules.Assert("03", "payment ext fr-ctc-flow6-condition must be a Payment-applicable CharacteristicTypeCode (MEN, MPA, RAP); status-only codes belong on a bill.Reason under bill.Status",
				tax.ExtensionsHasCodes(ExtKeyCondition, paymentConditionCodes...),
			),
		),
	)
}
