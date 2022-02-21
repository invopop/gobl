package num

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("num"), Amount{})
	schema.Register(schema.GOBL.Add("num"), Percentage{})
}
