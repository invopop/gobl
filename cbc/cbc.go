package cbc

import "github.com/invopop/gobl/schema"

/*
 * cbc for Common Basic Components
 *
 * Name is taken from the similar namespace in UBL.
 */

func init() {
	schema.Register(schema.GOBL.Add("cbc"),
		Code(""),
		Key(""),
		Meta{},
		Note{},
		Stamp{},
	)
}
