package mcp

import "encoding/json"

func marshalJSON(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}
