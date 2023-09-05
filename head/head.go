// Package base contains the main structural elements of GOBL in a single
// place so that they can be re-used inside the other packages.
package head

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("head"),
		Header{},
		Stamp{},
	)
}
