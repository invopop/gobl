// Package bill provides models for dealing with Billing and specifically invoicing.
package bill

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("bill"),
		// Primary schemas
		CorrectionOptions{},
		Delivery{},
		Invoice{},
		Order{},
		Payment{},
		Status{},
		// Sub-schemas - used by primaries
		Charge{},
		Discount{},
		Line{},
		Ordering{},
		PaymentDetails{},
		DeliveryDetails{},
		Tax{},
		Totals{},
	)
	rules.Register(
		"bill",
		rules.GOBL.Add("BILL"),
		invoiceRules(),
		deliveryRules(),
		orderRules(),
		paymentRules(),
		lineRules(),
		subLineRules(),
		lineDiscountRules(),
		lineChargeRules(),
		discountRules(),
		chargeRules(),
		paymentLineRules(),
		taxRules(),
		totalsRules(),
	)
}

// Constants used to help identify document schemas
const (
	ShortSchemaOrder    = "bill/order"
	ShortSchemaDelivery = "bill/delivery"
	ShortSchemaInvoice  = "bill/invoice"
	ShortSchemaPayment  = "bill/payment"
	ShortSchemaStatus   = "bill/status"
)
