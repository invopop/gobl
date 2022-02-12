package note

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/region"
)

// Message represents the minimum possible contents for a GoBL document type. This is
// mainly meant to be used for testing purposes.
type Message struct {
	// Summary of the message content
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// Details of what exactly this message wants to communicate
	Content string `json:"content" jsonschema:"title=Content"`
	// Any additional semi-structured data that might be useful.
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta Data"`
}

// Validate ensures the message contains everything it should.
func (m *Message) Validate(r region.Region) error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Content, validation.Required),
	)
}
