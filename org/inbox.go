package org

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Inbox is used to store data about a connection with a service that is responsible
// for potentially receiving copies of GOBL envelopes or other document formats
// defined locally.
type Inbox struct {
	uuid.Identify
	// Label for the inbox.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Type of inbox being defined.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Role assigned to this inbox that may be relevant for the consumer.
	Role cbc.Key `json:"role,omitempty" jsonschema:"title=Role"`
	// Code or ID that identifies the Inbox.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// URL of the inbox that includes the protocol, server, and path. May
	// be used instead of the Code to identify the inbox.
	URL string `json:"url,omitempty" jsonschema:"title=URL"`
	// Extension code map for any additional regime or addon specific codes that may be required.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Normalize will try to clean the inbox's data.
func (i *Inbox) Normalize(normalizers tax.Normalizers) {
	if i == nil {
		return
	}
	uuid.Normalize(&i.UUID)
	i.Code = cbc.NormalizeCode(i.Code)
	i.Ext = tax.CleanExtensions(i.Ext)
	normalizers.Each(i)
}

// Validate ensures the inbox's fields look good.
func (i *Inbox) Validate() error {
	return i.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the inbox's fields look good inside the provided context.
func (i *Inbox) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, i,
		validation.Field(&i.UUID),
		validation.Field(&i.Key),
		validation.Field(&i.Role),
		validation.Field(&i.Code,
			validation.When(
				i.URL == "",
				validation.Required.Error("cannot be blank without url"),
			),
		),
		validation.Field(&i.URL,
			is.URL,
			validation.When(
				i.Code != "",
				validation.Empty.Error("mutually exclusive with code"),
			),
		),
		validation.Field(&i.Ext),
	)
}

// AddInbox makes it easier to add a new inbox to a list and replace an
// existing value with a matching key.
func AddInbox(in []*Inbox, i *Inbox) []*Inbox {
	if in == nil {
		return []*Inbox{i}
	}
	for _, v := range in {
		if v.Key == i.Key {
			*v = *i // copy in place
			return in
		}
	}
	return append(in, i)
}
