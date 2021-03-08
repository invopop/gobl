package gobl

// Party represents a person or business entity.
type Party struct {
	UUID      string    `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code"`
	TaxID     string    `json:"tax_id,omitempty" jsonschema:"title=Tax ID"`
	TaxScheme string    `json:"tax_scheme,omitempty" jsonschema:"title=Tax Scheme"`
	Name      string    `json:"name" jsonschema:"title=Name"`
	Alias     string    `json:"alias,omitempty" jsonschema:"title=Alias"`
	Contacts  []Person  `json:"contacts,omitempty"`
	Addresses []Address `json:"addresses,omitempty" jsonschema:"title=Postal Addresses"`
	Meta      Meta      `json:"meta,omitempty"`
}

// Person represents a human. This is a complex subject, see this article
// for some insights: https://www.w3.org/International/questions/qa-personal-names
type Person struct {
	UUID     string `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identity code"`
	Role     string `json:"role,omitempty" jsonschema:"title=Role,description=Role within an organization"`
	Prefix   string `json:"prefix,omitempty"`
	Given    string `json:"given"`
	Middle   string `json:"middle,omitempty"`
	Surname  string `json:"surname"`
	Surname2 string `json:"surname2,omitempty"`
	Suffix   string `json:"suffix,omitempty"`
	Meta     Meta   `json:"meta,omitempty"`
}
