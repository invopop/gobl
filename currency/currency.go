// Package currency provides models for dealing with currencies.
package currency

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("currency"),
		Code(""),
		ExchangeRate{},
	)
}
