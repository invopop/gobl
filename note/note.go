package note

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.RegisterIn(schema.GOBL.Add("note"), &Message{})
}
