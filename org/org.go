package org

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("org"),
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
		Inbox{},
	)
}
