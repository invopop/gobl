// Package bill provides models for dealing with Billing and specifically invoicing.
package bill

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("bill"),
		Delivery{},
		Invoice{},
		Order{},
		Payment{},
		CorrectionOptions{},
		Line{},
	)
}

// Constants used to help identify document schemas
const (
	ShortSchemaOrder    = "bill/order"
	ShortSchemaDelivery = "bill/delivery"
	ShortSchemaInvoice  = "bill/invoice"
	ShortSchemaPayment  = "bill/payment"
)
