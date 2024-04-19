package org

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Inbox is used to store data about a connection with a service that is responsible
// for potentially receiving copies of GOBL envelopes or other document formats
// defined locally.
type Inbox struct {
	// Unique ID. Useful if inbox is stored in a database.
	UUID uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`

	// Type of inbox being defined.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Role assigned to this inbox that may be relevant for the consumer.
	Role cbc.Key `json:"role,omitempty" jsonschema:"title=Role"`

	// Human name for the inbox.
	Name string `json:"name,omitempty" jsonschema:"title=Name"`

	// Actual Code or ID that identifies the Inbox.
	Code string `json:"code"`
}

// Validate ensures the inbox's fields look good.
func (i *Inbox) Validate() error {
	return i.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the inbox's fields look good inside the provided context.
func (i *Inbox) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, i,
		validation.Field(&i.UUID),
		validation.Field(&i.Key, validation.Required),
		validation.Field(&i.Role),
		validation.Field(&i.Code, validation.Required),
	)
}
