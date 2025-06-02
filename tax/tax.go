// Package tax encapsulates models related to taxation.
package tax

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("tax"),
		Identity{},
		Set{},
		Extensions{},
		Total{},
		RegimeDef{},
		AddonDef{},
		CatalogueDef{},
	)
}
