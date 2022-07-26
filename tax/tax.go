package tax

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("tax"),
		Set{},
		Total{},
		Region{},
	)
}
