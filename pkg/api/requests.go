package api

import (
	"encoding/json"

	"github.com/invopop/gobl/dsig"
)

// Request types for the HTTP API. These mirror the ops.* request types but
// use json.RawMessage instead of []byte so that JSON payloads are accepted
// as-is rather than requiring base64 encoding.

type buildRequest struct {
	Template json.RawMessage `json:"template,omitempty"`
	Data     json.RawMessage `json:"data"`
	DocType  string          `json:"type,omitempty"`
	Envelop  bool            `json:"envelop,omitempty"`
}

type signRequest struct {
	Template   json.RawMessage  `json:"template,omitempty"`
	Data       json.RawMessage  `json:"data"`
	PrivateKey *dsig.PrivateKey `json:"privatekey,omitempty"`
	DocType    string           `json:"type,omitempty"`
	Envelop    bool             `json:"envelop,omitempty"`
}

type validateRequest struct {
	Data json.RawMessage `json:"data"`
}

type verifyRequest struct {
	Data      json.RawMessage `json:"data"`
	PublicKey *dsig.PublicKey `json:"publickey,omitempty"`
}

type correctRequest struct {
	Data    json.RawMessage `json:"data"`
	Options json.RawMessage `json:"options,omitempty"`
	Schema  bool            `json:"schema,omitempty"`
}

type replicateRequest struct {
	Data json.RawMessage `json:"data"`
}
