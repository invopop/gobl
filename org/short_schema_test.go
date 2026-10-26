package org_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

func TestShortSchemas(t *testing.T) {
	tests := map[string]any{
		org.ShortSchemaAttribute:    org.Attribute{},
		org.ShortSchemaDocumentRef:  org.DocumentRef{},
		org.ShortSchemaIdentity:     org.Identity{},
		org.ShortSchemaItem:         org.Item{},
		org.ShortSchemaNote:         org.Note{},
		org.ShortSchemaParty:        org.Party{},
		org.ShortSchemaRegistration: org.Registration{},
	}
	for short, obj := range tests {
		t.Run(short, func(t *testing.T) {
			id := schema.Lookup(obj)
			assert.True(t, strings.HasSuffix(id.String(), "/"+short), "schema ID %s", id)
		})
	}
}
