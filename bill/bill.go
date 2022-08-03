package bill

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("bill"),
		// None of bill's sub-models are meant to be used outside an invoice.
		Invoice{},
	)
}
