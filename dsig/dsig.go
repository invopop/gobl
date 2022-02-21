package dsig

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	objs := []interface{}{
		&Digest{},
		&Signature{},
	}
	schema.RegisterAll(schema.GOBL.Add("dsig"), objs)
}
