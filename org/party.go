package org

import (
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Party represents a person or business entity.
type Party struct {
	UUID       *uuid.UUID   `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code."`
	TaxID      *TaxID       `json:"tax_id,omitempty" jsonschema:"title=Tax Identity,description=The entity's legal ID code used for tax purposes. They may have other numbers, but we're only interested in those valid for tax pruposes."`
	Name       string       `json:"name" jsonschema:"title=Name,description=Legal name or representation of the organization."`
	Alias      string       `json:"alias,omitempty" jsonschema:"title=Alias,description=Alternate short name."`
	People     []*Person    `json:"people,omitempty" jsonschema:"title=People,description=Details of physical people who represent the party."`
	Addresses  []*Address   `json:"addresses,omitempty" jsonschema:"title=Postal Addresses,description=Regular post addresses for where information should be sent if needed."`
	Emails     []*Email     `json:"emails,omitempty" jsonschema:"title=Email Addresses"`
	Telephones []*Telephone `json:"telephones,omitempty" jsonschema:"title=Telephone Numbers"`
	Meta       Meta         `json:"meta,omitempty" jsonschema:"title=Meta,description=Any additional non-structure information that does not fit into the rest of the document."`
}

// Person represents a human, and how to contact them.
type Person struct {
	UUID   *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code"`
	Name   Name       `json:"name" jsonschema:"title=Name,description=Complete details on the name of the person"`
	Role   string     `json:"role,omitempty" jsonschema:"title=Role,description=What they do within an organization"`
	Emails []Email    `json:"emails,omitempty" jsonschema:"title=Email Addresses,description=Electronic mail addresses that belong to the person."`
	Meta   Meta       `json:"meta,omitempty" jsonschema:"title=Meta,description=Data about the data."`
}

// Name represents what a human is called. This is a complex subject, see this
// w3 article for some insights:
// https://www.w3.org/International/questions/qa-personal-names
type Name struct {
	UUID     *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code"`
	Alias    string     `json:"alias,omitempty" jsonschema:"title=Alias,description=What the person would like to be called"`
	Prefix   string     `json:"prefix,omitempty" jsonschema:"title=Prefix"`
	Given    string     `json:"given" jsonschema:"title=Given,description=The person's given name"`
	Middle   string     `json:"middle,omitempty" jsonschema:"title=Middle,description=Middle names or initials"`
	Surname  string     `json:"surname" jsonschema:"title=Surname"`
	Surname2 string     `json:"surname2,omitempty" jsonschema:"title=Second Surname"`
	Suffix   string     `json:"suffix,omitempty" jsonschema:"title=Suffix"`
	Meta     Meta       `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Email describes the electronic mailing details.
type Email struct {
	UUID    *uuid.UUID `json:"uuid,omitempty"`
	Label   string     `json:"label,omitempty" jsonschema:"title=Label,description=Identifier for the email."`
	Address string     `json:"addr" jsonschema:"title=Address,description=Electronic mailing address."`
	Meta    Meta       `json:"meta,omitempty" jsonschema:"title=Meta,description=Additional fields."`
}

// Telephone describes what is expected for a telephone number.
type Telephone struct {
	UUID   *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	Label  string     `json:"label,omitempty" jsonschema:"title=Label,description=Identifier for this number."`
	Number string     `json:"num" jsonschema:"title=Number,description=The number to be dialed in ITU E.164 international format."`
}

// Validate is used to check the party's data meets minimum expectations.
func (p *Party) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.TaxID),
		validation.Field(&p.People),
		validation.Field(&p.Emails),
		validation.Field(&p.Telephones),
	)
}

// Validate ensures email address looks valid.
func (e *Email) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Address, validation.Required, is.Email),
	)
}

// Validate checks the telephone objects number to ensure it looks correct.
func (t *Telephone) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Number, validation.Required, is.E164),
	)
}
