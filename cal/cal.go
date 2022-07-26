package cal

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("cal"),
		Date{},
		Period{},
	)
}
