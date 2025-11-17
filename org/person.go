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
	// Key used to identify the role of the person inside the context of the object.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Complete details on the name of the person.
	Name *Name `json:"name" jsonschema:"title=Name"`
	// Role or job title of the responsibilities of the person within an organization.
	Role string `json:"role,omitempty" jsonschema:"title=Role"`
	// Set of codes used to identify the person, such as ID numbers, social security,
	// driving licenses, etc. that can be attributed to the individual.
	Identities []*Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Regular post addresses for where information should be sent if needed.
	Addresses []*Address `json:"addresses,omitempty" jsonschema:"title=Postal Addresses"`
	// Electronic mail addresses that belong to the person.
	Emails []*Email `json:"emails,omitempty" jsonschema:"title=Email Addresses"`
	// Regular phone or mobile numbers
	Telephones []*Telephone `json:"telephones,omitempty" jsonschema:"title=Telephone Numbers"`
	// Avatars provider links to images or photos or the person.
	Avatars []*Image `json:"avatars,omitempty" jsonschema:"title=Avatars"`
	// Data about the data.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize will try to normalize the person's data.
func (p *Person) Normalize(normalizers tax.Normalizers) {
	if p == nil {
		return
	}

	uuid.Normalize(&p.UUID)
	p.Label = cbc.NormalizeString(p.Label)
	p.Role = cbc.NormalizeString(p.Role)

	tax.Normalize(normalizers, p.Name)
	tax.Normalize(normalizers, p.Identities)
	tax.Normalize(normalizers, p.Addresses)
	tax.Normalize(normalizers, p.Emails)
	tax.Normalize(normalizers, p.Telephones)
	tax.Normalize(normalizers, p.Avatars)
	normalizers.Each(p)
}

// Validate validates the person.
func (p *Person) Validate() error {
	return p.ValidateWithContext(context.Background())
}

// ValidateWithContext validates the person with the given context.
func (p *Person) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, p,
		validation.Field(&p.UUID),
		validation.Field(&p.Label),
		validation.Field(&p.Key),
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Identities),
		validation.Field(&p.Addresses),
		validation.Field(&p.Emails),
		validation.Field(&p.Telephones),
		validation.Field(&p.Avatars),
		validation.Field(&p.Meta),
	)
}
