package dsig

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("dsig"),
		&Digest{},
		&Signature{},
	)
}
