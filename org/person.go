package org

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Person represents a human, and how to contact them electronically.
type Person struct {
	uuid.Identify
	// Label can be used to identify the person in a given context in a single
	// language, for example "Attn", "Contact", "Responsible", etc.
	Label string `json:"label,omitempty" jsonschema:"title=Label,example=Attn"`
	// Complete details on the name of the person
	Name *Name `json:"name" jsonschema:"title=Name"`
	// What they do within an organization
	Role string `json:"role,omitempty" jsonschema:"title=Role"`
	// Electronic mail addresses that belong to the person.
	Emails []*Email `json:"emails,omitempty" jsonschema:"title=Email Addresses"`
	// Regular phone or mobile numbers
	Telephones []*Telephone `json:"telephones,omitempty" jsonschema:"title=Telephone Numbers"`
	// Avatars provider links to images or photos or the person.
	Avatars []*Image `json:"avatars,omitempty" jsonschema:"title=Avatars"`
	// Data about the data.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate validates the person.
func (p *Person) Validate() error {
	return p.ValidateWithContext(context.Background())
}

// ValidateWithContext validates the person with the given context.
func (p *Person) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, p,
		validation.Field(&p.UUID),
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Label),
		validation.Field(&p.Emails),
		validation.Field(&p.Telephones),
		validation.Field(&p.Avatars),
		validation.Field(&p.Meta),
	)
}
