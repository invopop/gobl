package org

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/uuid"
)

// Inboxes defines an array of inboxes.
type Inboxes []*Inbox

// InboxKey determines the type of inbox data being provided. Pre-defined
// keys are available in the regions package.
type InboxKey string

var (
	inboxKeyValidationRegexp = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
)

// Inbox is used to store data about a connection with a service that is responsible
// for potentially receiving copies of GOBL envelopes or other document formats
// defined locally.
type Inbox struct {
	// Unique ID. Useful if inbox is stored in a database.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`

	// Type of inbox being defined.
	Key InboxKey `json:"key" jsonschema:"title=Key"`

	// Role assigned to this inbox that may be relevant for the consumer.
	Role InboxKey `json:"role,omitempty" jsonschema:"title=Role"`

	// Human name for the inbox.
	Name string `json:"name,omitempty" jsonschema:"title=Name"`

	// Actual Code or ID that identifies the Inbox.
	Code string `json:"code"`
}

// Validate ensures the inbox's fields look good.
func (i *Inbox) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.UUID),
		validation.Field(&i.Key, validation.Required),
		validation.Field(&i.Role),
		validation.Field(&i.Code, validation.Required),
	)
}

// Validate ensures the key complies with the basic syntax
// requirements.
func (k InboxKey) Validate() error {
	return validation.Validate(string(k),
		validation.Match(inboxKeyValidationRegexp),
	)
}

// String provides string representation of key
func (k InboxKey) String() string {
	return string(k)
}
