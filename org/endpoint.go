package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// Endpoint identifies an address to which electronic documents may be
// sent, expressed as a single URI. The URI scheme identifies the
// network or namespace, e.g. "gobl:acme.example.com" (GOBL Net),
// "peppol:9920:x3157928m" (Peppol), or "mailto:billing@example.com".
type Endpoint struct {
	uuid.Identify

	// Label for the endpoint.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// URI that identifies the endpoint.
	URI cbc.URI `json:"uri" jsonschema:"title=URI"`
}

// Normalize will try to clean the endpoint's data.
func (e *Endpoint) Normalize() {
	if e == nil {
		return
	}
	uuid.Normalize(&e.UUID)
	e.Label = cbc.NormalizeString(e.Label)
}

func endpointRules() *rules.Set {
	return rules.For(new(Endpoint),
		rules.Field("uri",
			rules.Assert("01", "endpoint uri is required", is.Present),
		),
	)
}

// Endpoint returns the party's first endpoint whose URI uses the given
// scheme, or nil if none is present.
func (p *Party) Endpoint(scheme string) *Endpoint {
	for _, e := range p.Endpoints {
		if e != nil && e.URI.Scheme() == scheme {
			return e
		}
	}
	return nil
}
