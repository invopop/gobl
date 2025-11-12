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
		Note{},
		Party{},
		Person{},
		Registration{},
		Telephone{},
		Unit(""),
		Website{},
		Attachment{},
	)
}

// ShortSchemaParty is the short schema name for Party
const (
	ShortSchemaParty = "org/party"
)
