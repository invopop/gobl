// Package num provides support for dealing with amounts and percentages without
// rounding errors.
package num

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("num"), Amount{})
	schema.Register(schema.GOBL.Add("num"), Percentage{})
}
