// Package org contains structures related to organization.
package org

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("org"),
		Address{},
		Attribute{},
		Coordinates{},
		DocumentRef{},
		Email{},
		Endpoint{},
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
	rules.Register(
		"org",
		rules.GOBL.Add("ORG"),
		attachmentRules(),
		attributeRules(),
		coordinatesRules(),
		documentRefRules(),
		emailRules(),
		endpointRules(),
		identityRules(),
		imageRules(),
		inboxRules(),
		itemRules(),
		nameRules(),
		noteRules(),
		personRules(),
		registrationRules(),
		telephoneRules(),
		unitRules(),
		websiteRules(),
	)
}

// ShortSchemaParty is the short schema name for Party
const (
	ShortSchemaParty = "org/party"
	ShortSchemaItem  = "org/item"
)
