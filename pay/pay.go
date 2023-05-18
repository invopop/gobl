// Package pay handles models related to payments.
package pay

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("pay"),
		MeansKey(""),
		Advance{},
		Instructions{},
		Terms{},
	)
}
