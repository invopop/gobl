// Package bill provides models for dealing with Billing and specifically invoicing.
package bill

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("bill"),
		Invoice{},
		Receipt{},
		CorrectionOptions{},
	)
}

// Constants used to help identify document schemas
const (
	ShortSchemaInvoice = "bill/invoice"
	ShortSchemaPayment = "bill/payment"
)
