package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Party represents a person or business entity.
type Party struct {
	// Internal ID used to identify the party inside a document.
	ID string `json:"id,omitempty" jsonschema:"title=ID"`
	// Unique identity code
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// The entity's legal ID code used for tax purposes. They may have other numbers, but we're only interested in those valid for tax purposes.
	TaxID *tax.Identity `json:"tax_id,omitempty" jsonschema:"title=Tax Identity"`
	// Legal name or representation of the organization.
	Name string `json:"name" jsonschema:"title=Name"`
	// Alternate short name.
	Alias string `json:"alias,omitempty" jsonschema:"title=Alias"`
	// Details of physical people who represent the party.
	People []*Person `json:"people,omitempty" jsonschema:"title=People"`
	// Digital inboxes used for forwarding electronic versions of documents
	Inboxes []*Inbox `json:"inboxes,omitempty" jsonschema:"title=Inboxes"`
	// Regular post addresses for where information should be sent if needed.
	Addresses []*Address `json:"addresses,omitempty" jsonschema:"title=Postal Addresses"`
	// Electronic mail addresses
	Emails []*Email `json:"emails,omitempty" jsonschema:"title=Email Addresses"`
	// Regular telephone numbers
	Telephones []*Telephone `json:"telephones,omitempty" jsonschema:"title=Telephone Numbers"`
	// Additional registration details about the company that may need to be included in a document.
	Registration *Registration `json:"registration,omitempty" jsonschema:"title=Registration"`
	// Any additional semi-structured information that does not fit into the rest of the party.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Person represents a human, and how to contact them electronically.
type Person struct {
	// Internal ID used to identify the person inside a document.
	ID string `json:"id,omitempty" jsonschema:"title=ID"`
	// Unique identity code
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Complete details on the name of the person
	Name Name `json:"name" jsonschema:"title=Name"`
	// What they do within an organization
	Role string `json:"role,omitempty" jsonschema:"title=Role"`
	// Electronic mail addresses that belong to the person.
	Emails []*Email `json:"emails,omitempty" jsonschema:"title=Email Addresses"`
	// Regular phone or mobile numbers
	Telephones []*Telephone `json:"telephones,omitempty" jsonschema:"title=Telephone Numbers"`
	// Data about the data.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Name represents what a human is called. This is a complex subject, see this
// w3 article for some insights:
// https://www.w3.org/International/questions/qa-personal-names
type Name struct {
	// Unique identity code
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// What the person would like to be called
	Alias string `json:"alias,omitempty" jsonschema:"title=Alias"`
	// Additional prefix to add to name, like Mrs. or Mr.
	Prefix string `json:"prefix,omitempty" jsonschema:"title=Prefix"`
	// Person's given or first name
	Given string `json:"given" jsonschema:"title=Given"`
	// Middle names or initials
	Middle string `json:"middle,omitempty" jsonschema:"title=Middle"`
	// Second or Family name.
	Surname string `json:"surname" jsonschema:"title=Surname"`
	// Additional second of family name.
	Surname2 string `json:"surname2,omitempty" jsonschema:"title=Second Surname"`
	// Titles to include after the name.
	Suffix string `json:"suffix,omitempty" jsonschema:"title=Suffix"`
	// Any additional useful data.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Email describes the electronic mailing details.
type Email struct {
	// Unique identity code
	UUID *uuid.UUID `json:"uuid,omitempty"`
	// Identifier for the email.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Electronic mailing address.
	Address string `json:"addr" jsonschema:"title=Address"`
	// Additional fields.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Telephone describes what is expected for a telephone number.
type Telephone struct {
	// Unique identity code
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Identifier for this number.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// The number to be dialed in ITU E.164 international format.
	Number string `json:"num" jsonschema:"title=Number"`
}

// Registration is used in countries that require additional information to be associated
// with a company usually related to a specific registration office.
// The definition found here is based on the details required for spain.
// If your country requires additional fields, please let us know.
type Registration struct {
	UUID    *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	Office  string     `json:"office,omitempty" jsonschema:"title=Office"`
	Book    string     `json:"book,omitempty" jsonschema:"title=Book"`
	Volume  string     `json:"volume,omitempty" jsonschema:"title=Volume"`
	Sheet   string     `json:"sheet,omitempty" jsonschema:"title=Sheet"`
	Section string     `json:"section,omitempty" jsonschema:"title=Section"`
	Page    string     `json:"page,omitempty" jsonschema:"title=Page"`
	Entry   string     `json:"entry,omitempty" jsonschema:"title=Entry"`
}

// Calculate performs any calculations required on the Party or
// it's properties, like the tax identity.
func (p *Party) Calculate() error {
	if p.TaxID != nil {
		if err := p.TaxID.Calculate(); err != nil {
			return err
		}
		r := p.TaxID.Regime()
		return r.CalculateDocument(p)
	}
	return nil
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
		validation.Field(&e.Address, validation.Required, is.EmailFormat),
	)
}

// Validate checks the telephone objects number to ensure it looks correct.
func (t *Telephone) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Number, validation.Required, is.E164),
	)
}
