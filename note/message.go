package note

import (
	"encoding/json"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

// Message represents a simple message object with a title and some
// content meant.
type Message struct {
	uuid.Identify
	// Summary of the message content
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// Details of what exactly this message wants to communicate.
	Content string `json:"content" jsonschema:"title=Content"`
	// Any additional semi-structured data that might be useful.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta Data"`
}

func messageRules() *rules.Set {
	return rules.For(new(Message),
		rules.Field("content",
			rules.Assert("01", "message content is required", rules.Present),
		),
	)
}

// JSONSchemaExtend adds examples to the JSON Schema.
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
