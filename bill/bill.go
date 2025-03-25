// Package bill provides models for dealing with Billing and specifically invoicing.
package bill

import (
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
}

// Constants used to help identify document schemas
const (
	ShortSchemaOrder    = "bill/order"
	ShortSchemaDelivery = "bill/delivery"
	ShortSchemaInvoice  = "bill/invoice"
	ShortSchemaPayment  = "bill/payment"
)
