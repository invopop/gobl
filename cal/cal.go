package cal

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		Date{},
		Period{},
	}
	schema.RegisterAll(schema.GOBL.Add("cal"), objs)
}
