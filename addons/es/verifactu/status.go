package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billStatusRules() *rules.Set {
	return rules.For(new(bill.Status),
		// Supplier required
		rules.Field("supplier",
			rules.Assert("01", "supplier is required", is.Present),
			rules.Field("tax_id",
				rules.Assert("02", "supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("03", "supplier tax ID code is required", is.Present),
				),
			),
		),
		// Event type extension required
		rules.Field("ext",
			rules.Assert("04", "event type extension is required",
				tax.ExtensionsRequire(ExtKeyEventType),
			),
		),
	)
}
