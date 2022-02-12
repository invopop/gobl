package org

import "github.com/invopop/gobl/schema"

func init() {
	objs := []interface{}{
		Address{},
		Coordinates{},
		Date{},
		Item{},
		ItemCode{},
		Notes{},
		Party{},
		Person{},
		Name{},
		Email{},
		Telephone{},
		Registration{},
		Period{},
		TaxID{},
	}
	schema.RegisterAllIn(schema.GOBL.Add("org"), objs)
}
