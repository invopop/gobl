package org

import (
	"github.com/asaskevich/govalidator"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

const (
	// InboxKeyPeppol is the key used to identify a Peppol inbox and apply
	// normalization rules for participant IDs.
	InboxKeyPeppol cbc.Key = "peppol"
)

// Inbox is used to store data about a connection with a service that is responsible
// for automatically receiving copies of GOBL envelopes or other document formats.
//
// In other formats, this may also be known as an "endpoint".
type Inbox struct {
	uuid.Identify

	// Label for the inbox.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Type of inbox being defined if required for clarification between multiple
	// inboxes.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Scheme ID of the code used to identify the inbox. This is context specific
	// and usually an ISO 6523 code or CEF (Connecting Europe Facility) code.
	Scheme cbc.Code `json:"scheme,omitempty" jsonschema:"title=Scheme"`

	// Code or ID that identifies the Inbox. Mutually exclusive with URL and Email.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// URL of the inbox that includes the protocol, server, and path. May
	// be used instead of the Code to identify the inbox. Mutually exclusive with
	// Code and Email.
	URL string `json:"url,omitempty" jsonschema:"title=URL"`
	// Email address for the inbox. Mutually exclusive with Code and URL.
	Email string `json:"email,omitempty" jsonschema:"title=Email"`
}

// Normalize will try to clean the inbox's data.
func (i *Inbox) Normalize() {
	if i == nil {
		return
	}
	uuid.Normalize(&i.UUID)
	code := i.Code.String()
	if govalidator.IsEmail(code) {
		i.Email = code
		i.Code = ""
	} else if govalidator.IsURL(code) {
		i.URL = code
		i.Code = ""
	}
	i.Label = cbc.NormalizeString(i.Label)
	i.Scheme = cbc.NormalizeAlphanumericalCode(i.Scheme)
	i.Code = cbc.NormalizeCode(i.Code)

	// Custom normalizations
	switch i.Key {
	case InboxKeyPeppol:
		if i.Scheme == "" {
			if len(i.Code) >= 5 && i.Code[4] == ':' {
				numbers := i.Code[:4]
				i.Scheme = numbers
				i.Code = i.Code[5:]
			}
		}
	}
}

func inboxRules() *rules.Set {
	return rules.For(new(Inbox),
		rules.Assert("01", "inbox requires a code, url, or email",
			is.Func("one of code, url, email required", func(val any) bool {
				i, ok := val.(*Inbox)
				return ok && i != nil && (i.Code != "" || i.URL != "" || i.Email != "")
			}),
		),
		rules.Assert("02", "inbox url must be blank when code or email is provided",
			is.Func("url exclusive", func(val any) bool {
				i, ok := val.(*Inbox)
				return ok && i != nil && (i.URL == "" || (i.Code == "" && i.Email == ""))
			}),
		),
		rules.Assert("03", "inbox email must be blank when code or url is provided",
			is.Func("email exclusive", func(val any) bool {
				i, ok := val.(*Inbox)
				return ok && i != nil && (i.Email == "" || (i.Code == "" && i.URL == ""))
			}),
		),
		rules.Field("url",
			rules.AssertIfPresent("04", "inbox url must be valid", is.URL),
		),
		rules.Field("email",
			rules.AssertIfPresent("05", "inbox email must be valid", is.EmailFormat),
		),
	)
}

// AddInbox makes it easier to add a new inbox to a list and replace an
// existing value with a matching key.
func AddInbox(in []*Inbox, i *Inbox) []*Inbox {
	if i == nil {
		return in
	}
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
