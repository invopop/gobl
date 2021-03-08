package gobl

// Party represents a person or business entity.
type Party struct {
	UUID      string    `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code"`
	TaxID     string    `json:"tax_id,omitempty" jsonschema:"title=Tax ID"`
	TaxScheme string    `json:"tax_scheme,omitempty" jsonschema:"title=Tax Scheme"`
	Name      string    `json:"name" jsonschema:"title=Name"`
	Alias     string    `json:"alias,omitempty" jsonschema:"title=Alias"`
	People    []Person  `json:"people,omitempty"`
	Addresses []Address `json:"addresses,omitempty" jsonschema:"title=Postal Addresses"`
	Emails    []Email   `json:"emails,omitempty" jsonschema:"title=Email Addresses"`
	Meta      Meta      `json:"meta,omitempty"`
}

// Person represents a human, and how to contact them.
type Person struct {
	UUID   string  `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code"`
	Role   string  `json:"role,omitempty" jsonschema:"title=Role,description=Role within an organization"`
	Name   Name    `json:"name"`
	Emails []Email `json:"emails,omitempty"`
}

// Name represents what a human is called. This is a complex subject, see this
// w3 article for some insights:
// https://www.w3.org/International/questions/qa-personal-names
type Name struct {
	UUID     string `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code."`
	Alias    string `json:"alias,omitempty" jsonschema:"title=Alias,description=What the person would like to be called."`
	Prefix   string `json:"prefix,omitempty"`
	Given    string `json:"given"`
	Middle   string `json:"middle,omitempty"`
	Surname  string `json:"surname"`
	Surname2 string `json:"surname2,omitempty"`
	Suffix   string `json:"suffix,omitempty"`
	Meta     Meta   `json:"meta,omitempty"`
}

// Email describes the electronic mailing details.
type Email struct {
	UUID    string `json:"uuid,omitempty"`
	Label   string `json:"label,omitempty" jsonschema:"title=Label,description=Identifier for the email."`
	Address string `json:"addr,omitempty" jsonschema:"title=Address,description=Electronic mailing address."`
}
