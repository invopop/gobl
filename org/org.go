package org

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		Code(""),
		Key(""),
		Address{},
		Coordinates{},
		Item{},
		Party{},
		Person{},
		Name{},
		Email{},
		Telephone{},
		Registration{},
		TaxIdentity{},
		Meta{},
		Notes{},
		Inbox{},
	}
	schema.RegisterAll(schema.GOBL.Add("org"), objs)
}
