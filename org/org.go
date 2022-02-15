package org

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		Address{},
		Coordinates{},
		Date{},
		Item{},
		Note{},
		Party{},
		Person{},
		Name{},
		Email{},
		Telephone{},
		Registration{},
		Period{},
		TaxID{},
	}
	schema.RegisterAll(schema.GOBL.Add("org"), objs)
}
