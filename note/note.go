// Package note provides models for generating simple messages.
package note

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("note"), &Message{})
}
