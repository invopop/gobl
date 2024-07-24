// Package head defines the contents to be used in envelope
// headers.
package head

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("head"),
		Header{},
		Stamp{},
		Link{},
	)
}
