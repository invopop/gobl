package org

import (
	"context"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Identity is used to define a code for a specific context.
type Identity struct {
	// Unique identity for this identity object.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Optional label useful for non-standard identities to give a bit more context.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// The type of Code being represented and usually specific for
	// a particular context, country, or tax regime, and cannot be used
	// alongside the key.
	Type cbc.Code `json:"type,omitempty" jsonschema:"title=Type"`
	// The actual value of the identity code.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
	// Description adds details about what the code could mean or imply
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
}

// Validate ensures the identity looks valid.
func (i *Identity) Validate() error {
	return i.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the identity looks valid inside the provided context.
func (i *Identity) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, i,
		validation.Field(&i.Label),
		validation.Field(&i.Type),
		validation.Field(&i.Code,
			validation.Required,
		),
	)
}

// HasIdentityType provides a validation rule that will determine if at least one
// of the identities defined includes one with the defined type.
func HasIdentityType(typ cbc.Code) validation.Rule {
	return validateIdentitySet{typ: typ}
}

type validateIdentitySet struct {
	typ cbc.Code
}

func (v validateIdentitySet) Validate(value interface{}) error {
	ids, ok := value.([]*Identity)
	if !ok {
		return nil
	}
	for _, row := range ids {
		if row.Type == v.typ {
			return nil
		}
	}
	return fmt.Errorf("missing %s", v.typ)
}

// IdentityForType helps return the identity with a matching type code.
func IdentityForType(in []*Identity, typ cbc.Code) *Identity {
	for _, v := range in {
		if v.Type == typ {
			return v
		}
	}
	return nil
}

// AddIdentity makes it easier to add a new identity to a list and replace an
// existing value with a matching type.
func AddIdentity(in []*Identity, i *Identity) []*Identity {
	if in == nil {
		return []*Identity{i}
	}
	for _, v := range in {
		if v.Type == i.Type {
			*v = *i // copy in place
			return in
		}
	}
	return append(in, i)
}
