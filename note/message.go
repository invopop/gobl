package note

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/schema"
)

// Message represents the minimum possible contents for a GoBL document type. This is
// mainly meant to be used for testing purposes.
type Message struct {
	schema.Def
	Title   string   `json:"title,omitempty" jsonschema:"title=Title,description=Summary of the message content."`
	Content string   `json:"content" jsonschema:"title=Content,description=Details of what exactly this message wants to communicate."`
	Meta    org.Meta `json:"meta,omitempty" jsonschema:"title=Meta Data,description=Any additional semi-structured data that might be useful."`
}

// NewMessage instantiates a new message with the correct base data.
func NewMessage() *Message {
	m := new(Message)
	m.Schema = MessageType.ID()
	return m
}

// Validate ensures the message contains everything it should.
func (m *Message) Validate(r region.Region) error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Schema, validation.Required),
		validation.Field(&m.Content, validation.Required),
	)
}
