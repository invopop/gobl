// Package bill provides models for dealing with Billing and specifically invoicing.
package bill

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("bill"),
		Invoice{},
		Payment{},
		CorrectionOptions{},
	)
}

// Constants used to help identify document schemas
const (
	ShortSchemaOrder    = "bill/order"
	ShortSchemaDelivery = "bill/delivery"
	ShortSchemaInvoice  = "bill/invoice"
	ShortSchemaPayment  = "bill/payment"
)
