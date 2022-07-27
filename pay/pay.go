package pay

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("pay"),
		Advance{},
		Instructions{},
		Terms{},
	)
}
