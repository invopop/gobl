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
	// Key is used to classify the identity for a specific tax regime.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
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
		validation.Field(&i.Key),
		validation.Field(&i.Type,
			validation.When(i.Key != "", validation.Empty.Error("must be blank when key set")),
		),
		validation.Field(&i.Code,
			validation.Required,
			validation.When(i.Key != cbc.KeyEmpty,
				validation.By(validateIdentityCodeForKeyInRegime(ctx, i.Key)),
			),
		),
	)
}

func validateIdentityCodeForKeyInRegime(ctx context.Context, key cbc.Key) validation.RuleFunc {
	return func(value interface{}) error {
		code, ok := value.(cbc.Code)
		if !ok || code == "" {
			return nil
		}
		r := tax.RegimeFromContext(ctx)
		if r == nil {
			return nil // nothing to do without regime
		}
		var codes []*tax.CodeDefinition
		for _, kd := range r.Identities {
			if kd.Key == key {
				codes = kd.Codes
				break
			}
		}
		if len(codes) == 0 {
			return nil
		}
		for _, cd := range codes {
			if cd.Code == code {
				return nil
			}
		}
		return fmt.Errorf("invalid %s", key)
	}
}

// IdentityKeyIn returns a validation rule to help determine if the identity contains
// the expected code.
func IdentityKeyIn(keys ...cbc.Key) validation.Rule {
	out := make([]interface{}, len(keys))
	for i, l := range keys {
		out[i] = l
	}
	return validateIdentity{keyIn: out}
}

type validateIdentity struct {
	keyIn []interface{}
}

func (v validateIdentity) Validate(value interface{}) error {
	id, ok := value.(*Identity)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Key,
			validation.When(len(v.keyIn) > 0, validation.In(v.keyIn...)),
		),
	)
}

// HasIdentityKey provides a validation rule that will determine if at least one
// of the identities defined includes one with the defined key.
func HasIdentityKey(key cbc.Key) validation.Rule {
	return validateIdentitySet{key: key}
}

type validateIdentitySet struct {
	key cbc.Key
}

func (v validateIdentitySet) Validate(value interface{}) error {
	ids, ok := value.([]*Identity)
	if !ok {
		return nil
	}
	for _, row := range ids {
		if row.Key == v.key {
			return nil
		}
	}
	return fmt.Errorf("missing %s", v.key)
}

// IdentityForKey helps return the identity with a matching key.
func IdentityForKey(in []*Identity, key cbc.Key) *Identity {
	for _, v := range in {
		if v.Key == key {
			return v
		}
	}
	return nil
}

// AddIdentity makes it easier to add a new identity to a list and replace an
// existing value with a matching key.
func AddIdentity(in []*Identity, i *Identity) []*Identity {
	if in == nil {
		return []*Identity{i}
	}
	for _, v := range in {
		if v.Key == i.Key {
			*v = *i // copy in place
			return in
		}
	}
	return append(in, i)
}
