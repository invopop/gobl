package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// billPaymentRules validates only the supported payment shape: Flow 10
// e-reporting carries payment receipts. The e-reporting business rules
// are the converter's responsibility — see the package doc.
func billPaymentRules() *rules.Set {
	return rules.For(new(bill.Payment),
		rules.Field("type",
			rules.Assert("01", "payment type must be 'receipt' for Flow 10 reporting",
				is.In(bill.PaymentTypeReceipt),
			),
		),
	)
}
