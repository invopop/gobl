package org

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("org"),
		Code(""),
		Key(""),
		Unit(""),
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
		Note{},
		Inbox{},
	)
}
