package org

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		Address{},
		Coordinates{},
		Item{},
		Note{},
		Party{},
		Person{},
		Name{},
		Email{},
		Telephone{},
		Registration{},
		TaxIdentity{},
	}
	schema.RegisterAll(schema.GOBL.Add("org"), objs)
}
