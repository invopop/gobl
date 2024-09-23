// Package org contains structures related to organization.
package org

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("org"),
		Address{},
		Coordinates{},
		DocumentRef{},
		Email{},
		Identity{},
		Image{},
		Inbox{},
		Item{},
		Name{},
		Party{},
		Person{},
		Registration{},
		Telephone{},
		Unit(""),
		Website{},
	)
}
