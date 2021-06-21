package pay

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/invopop/gobl/uuid"
)

// URL provides the details required to make a payment online using a website
type URL struct {
	UUID    *uuid.UUID        `json:"uuid,omitempty" jsonschema:"title=UUID"`
	Address string            `json:"addr" jsonschema:"title=Address,description=Full URL to be used for payment."`
	Notes   string            `json:"notes,omitempty" jsonschema:"title=Notes,description=Additional details which may be useful to the end-user."`
	Meta    map[string]string `json:"meta,omitempty" jsonschema:"title=Meta,description=Semi-structured additional data."`
}

// Validate ensures the URL's details look correct.
func (u *URL) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(u.Address, validation.Required, is.URL),
	)
}
