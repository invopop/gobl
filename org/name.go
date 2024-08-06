package org

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Name represents what a human is called. This is a complex subject, see this
// w3 article for some insights:
// https://www.w3.org/International/questions/qa-personal-names
type Name struct {
	uuid.Identify
	// What the person would like to be called
	Alias string `json:"alias,omitempty" jsonschema:"title=Alias"`
	// Additional prefix to add to name, like Mrs. or Mr.
	Prefix string `json:"prefix,omitempty" jsonschema:"title=Prefix"`
	// Person's given or first name
	Given string `json:"given,omitempty" jsonschema:"title=Given"`
	// Middle names or initials
	Middle string `json:"middle,omitempty" jsonschema:"title=Middle"`
	// Second or Family name.
	Surname string `json:"surname,omitempty" jsonschema:"title=Surname"`
	// Additional second of family name.
	Surname2 string `json:"surname2,omitempty" jsonschema:"title=Second Surname"`
	// Titles to include after the name.
	Suffix string `json:"suffix,omitempty" jsonschema:"title=Suffix"`
	// Any additional useful data.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the name looks valid.
func (n *Name) Validate() error {
	return n.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the name looks valid inside the provided context.
func (n *Name) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, n,
		validation.Field(&n.UUID),
		validation.Field(&n.Given,
			validation.When(n.Surname == "",
				validation.Required,
			),
		),
		validation.Field(&n.Surname,
			validation.When(n.Given == "",
				validation.Required,
			),
		),
		validation.Field(&n.Meta),
	)
}
