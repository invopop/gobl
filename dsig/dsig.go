// Package dsig provides models for dealing with digital signatures.
package dsig

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("dsig"),
		&Digest{},
		&Signature{},
	)
	rules.Register(
		"dsig",
		rules.GOBL.Add("DSIG"),
		digestRules(),
	)
}

func digestRules() *rules.Set {
	return rules.For(new(Digest),
		rules.Field("alg",
			rules.Assert("01", "algorithm is required", is.Present),
		),
		rules.Field("val",
			rules.Assert("02", "value is required", is.Present),
		),
	)
}
