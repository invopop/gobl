package org

import (
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Party represents a person or business entity.
type Party struct {
	ID           string        `json:"id,omitempty" jsonschema:"title=ID,description=Internal ID used to identify the party inside a document."`
	UUID         *uuid.UUID    `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code."`
	TaxID        *TaxID        `json:"tax_id,omitempty" jsonschema:"title=Tax Identity,description=The entity's legal ID code used for tax purposes. They may have other numbers, but we're only interested in those valid for tax pruposes."`
	Name         string        `json:"name" jsonschema:"title=Name,description=Legal name or representation of the organization."`
	Alias        string        `json:"alias,omitempty" jsonschema:"title=Alias,description=Alternate short name."`
	People       []*Person     `json:"people,omitempty" jsonschema:"title=People,description=Details of physical people who represent the party."`
	Addresses    []*Address    `json:"addresses,omitempty" jsonschema:"title=Postal Addresses,description=Regular post addresses for where information should be sent if needed."`
	Emails       []*Email      `json:"emails,omitempty" jsonschema:"title=Email Addresses"`
	Telephones   []*Telephone  `json:"telephones,omitempty" jsonschema:"title=Telephone Numbers"`
	Registration *Registration `json:"registration,omitempty" jsonschema:"title=Registration,description=Additional registration details about the company that may need to be included in a document."`
	Meta         Meta          `json:"meta,omitempty" jsonschema:"title=Meta,description=Any additional semi-structured information that does not fit into the rest of the party."`
}

// Person represents a human, and how to contact them electronically.
type Person struct {
	ID         string       `json:"id,omitempty" jsonschema:"title=ID,description=Internal ID used to identify the person inside a document."`
	UUID       *uuid.UUID   `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code"`
	Name       Name         `json:"name" jsonschema:"title=Name,description=Complete details on the name of the person"`
	Role       string       `json:"role,omitempty" jsonschema:"title=Role,description=What they do within an organization"`
	Emails     []*Email     `json:"emails,omitempty" jsonschema:"title=Email Addresses,description=Electronic mail addresses that belong to the person."`
	Telephones []*Telephone `json:"telephones,omitempty" jsonschema:"title=Telephone Numbers"`
	Meta       Meta         `json:"meta,omitempty" jsonschema:"title=Meta,description=Data about the data."`
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

// Registration is used in countries that require additional information to be associated
// with a company usually related to a specific registration office.
// The definition found here is based on the details required for spain.
// If your country requires additional fields, please let us know.
type Registration struct {
	UUID    *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	Office  string     `json:"office,omitempty" jsonschema:"title=Office,description=Office where the company is registered."`
	Book    string     `json:"book,omitempty" jsonschema:"title=Book"`
	Volume  string     `json:"volume,omitempty" jsonschema:"title=Volume"`
	Sheet   string     `json:"sheet,omitempty" jsonschema:"title=Sheet"`
	Section string     `json:"section,omitempty" jsonschema:"title=Section"`
	Page    string     `json:"page,omitempty" jsonschema:"title=Page"`
	Entry   string     `json:"entry,omitempty" jsonschema:"title=Entry"`
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
