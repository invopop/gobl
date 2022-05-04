package tax

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		// Rates{},
		Total{},
		Region{},
	}
	schema.RegisterAll(schema.GOBL.Add("tax"), objs)
}
