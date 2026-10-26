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
		partyRules(),
		personRules(),
		registrationRules(),
		telephoneRules(),
		websiteRules(),
	)
}

// Short schema IDs of objects used inside documents
const (
	ShortSchemaAttribute    = "org/attribute"
	ShortSchemaDocumentRef  = "org/document-ref"
	ShortSchemaIdentity     = "org/identity"
	ShortSchemaItem         = "org/item"
	ShortSchemaNote         = "org/note"
	ShortSchemaParty        = "org/party"
	ShortSchemaRegistration = "org/registration"
)
