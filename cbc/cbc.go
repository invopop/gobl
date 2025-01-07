// Package cbc provides a set of Common Basic Components.
//
// Name is take from the similar namespace used in UBL.
package cbc

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("cbc"),
		Code(""),
		CodeMap{},
		Key(""),
		Definition{},
		Meta{},
	)
}
