package note

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
)

// Message represents the minimum possible contents for a GoBL document type. This is
// mainly meant to be used for testing purposes.
type Message struct {
	Title   string   `json:"title,omitempty" jsonschema:"title=Title,description=Summary of the message content."`
	Content string   `json:"content" jsonschema:"title=Content,description=Details of what exactly this message wants to communicate."`
	Meta    org.Meta `json:"meta,omitempty" jsonschema:"title=Meta Data,description=Any additional semi-structured data that might be useful."`
}

// Type provides the document type used for mapping.
func (Message) Type() string {
	return MessageType
}

// Validate ensures the message contains everything it should.
func (m *Message) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Content, validation.Required),
	)
}
