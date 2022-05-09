package tax

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		Set{},
		Total{},
		Region{},
	}
	schema.RegisterAll(schema.GOBL.Add("tax"), objs)
}
