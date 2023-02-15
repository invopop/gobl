package note

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/internal/here"
	"github.com/invopop/jsonschema"
)

// Message represents a simple message object with a title and some
// content meant.
type Message struct {
	// Summary of the message content
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// Details of what exactly this message wants to communicate.
	Content string `json:"content" jsonschema:"title=Content"`
	// Any additional semi-structured data that might be useful.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta Data"`
}

// Validate ensures the message contains everything it should.
func (m *Message) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Content, validation.Required),
		validation.Field(&m.Meta),
	)
}

func (Message) JSONSchemaExtend(s *jsonschema.Schema) {
	exs := here.Bytes(`
		[
			{
				"title": "Example Title",
				"content": "This is an example message."
			}
		]`)
	if err := json.Unmarshal(exs, &s.Examples); err != nil {
		panic(err)
	}
}
