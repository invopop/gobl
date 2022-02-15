package pay

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		Advance{},
		Instructions{},
		Terms{},
	}
	schema.RegisterAll(schema.GOBL.Add("pay"), objs)
}
