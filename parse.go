package gobl

import (
	"encoding/json"

	"github.com/invopop/gobl/schema"
)

// Parse unmarshals the provided data and uses the schema ID
// to determine what type of object we're dealing with. As long as the
// provided data contains a schema registered in GOBL, a new
// object instance will be returned.
func Parse(data []byte) (interface{}, error) {
	id, err := schema.Extract(data)
	if err != nil {
		return nil, ErrUnmarshal.WithCause(err)
	}
	if id == schema.UnknownID {
		return nil, ErrUnknownSchema
	}

	obj := id.Interface()
	if err := json.Unmarshal(data, obj); err != nil {
		return nil, ErrUnmarshal.WithCause(err)
	}

	return obj, nil
}
