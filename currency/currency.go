// Package currency provides models for dealing with currencies.
package currency

import (
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/schema"
)

func init() {
	definitions = new(defs)
	if err := definitions.load(data.Content, "currency"); err != nil {
		panic(err)
	}
	schema.Register(schema.GOBL.Add("currency"),
		Code(""),
		ExchangeRate{},
	)
}
