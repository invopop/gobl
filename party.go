package gobl

// Party represents a person or business entity.
type Party struct {
	TaxID     string   `json:"tax_id,omitempty" jsonschema:"title=Tax ID"`
	TaxScheme string   `json:"tax_scheme,omitempty" jsonschema:"title=Tax Scheme"`
	Name      string   `json:"name" jsonschema:"title=Name"`
	Alias     string   `json:"alias,omitempty" jsonschema:"title=Alias"`
	Address   *Address `json:"address,omitempty" jsonschema:"title=Postal Address"`
}

// Person represents a human. This is a complex subject, see this article
// for some insights: https://www.w3.org/International/questions/qa-personal-names
type Person struct {
	Prefix   string `json:"prefix,omitempty"`
	Given    string `json:"given"`
	Middle   string `json:"middle,omitempty"`
	Surname  string `json:"surname"`
	Surname2 string `json:"surname2,omitempty"`
	Suffix   string `json:"suffix,omitempty"`
	Alias    string `json:"alias,omitempty" jsonschema:"title=Alias,description=Nick or alternative name"`
}

// Address represents a physical location.
type Address struct {
	Name        string `json:"name,omitempty" jsonschema:"title=Name,description=Building or house name"`
	Number      string `json:"number,omitempty" jsonschema:"title=Number"`
	Interior    string `json:"interior,omitempty" jsonschema:"title=Interior"`
	StreetName  string `json:"street_name,omitempty" jsonschema:"title=Street Name"`
	StreetExtra string `json:"street_extra,omitempty" jsonschema:"title=Street Extra"`
	City        string `json:"city,omitempty" jsonschema:"title=City"`
	District    string `json:"district,omitempty" jsonschema:"title=District"`
	State       string `json:"state,omitempty" jsonschema:"title=State"`
	Country     string `json:"country" jsonschema:"title=Country Code,description=ISO 3166-1 alpha-2 two-letter country code"`
	PostCode    string `json:"post_code,omitempty" jsonschema:"title=Post Code"`
}
