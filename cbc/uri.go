package cbc

import (
	"fmt"
	"net/url"

	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/jsonschema"
)

// URIMaxLength is the maximum number of characters allowed in a URI.
var URIMaxLength uint64 = 2048

// URI is a Uniform Resource Identifier in `scheme:opaque` form, used to
// identify a resource or address. Examples include "gobl:acme.example.com",
// "iso6523-actorid-upis::9920:b123123123" (a Peppol participant
// identifier), and "mailto:billing@example.com". The scheme names the
// namespace; its interpretation is left to consumers.
type URI string

// String provides the string representation of the URI.
func (u URI) String() string {
	return string(u)
}

// Parse parses the URI using the standard library.
func (u URI) Parse() (*url.URL, error) {
	return url.Parse(string(u))
}

// Scheme returns the URI's scheme (the part before the first colon), or
// an empty string if the URI cannot be parsed.
func (u URI) Scheme() string {
	p, err := url.Parse(string(u))
	if err != nil {
		return ""
	}
	return p.Scheme
}

// Opaque returns the URI's scheme-specific part, excluding any query
// component, or an empty string if the URI cannot be parsed.
func (u URI) Opaque() string {
	p, err := url.Parse(string(u))
	if err != nil {
		return ""
	}
	return p.Opaque
}

// Host returns the URI's host component (for hierarchical URIs like
// "https://acme.example/x"), or an empty string if the URI cannot be
// parsed or carries no host (e.g. opaque URIs such as "mailto:a@b").
func (u URI) Host() string {
	p, err := url.Parse(string(u))
	if err != nil {
		return ""
	}
	return p.Host
}

// Path returns the URI's path component (for hierarchical URIs), or
// an empty string if the URI cannot be parsed or carries no path.
func (u URI) Path() string {
	p, err := url.Parse(string(u))
	if err != nil {
		return ""
	}
	return p.Path
}

// Validate ensures the URI is well-formed.
func (u URI) Validate() error {
	return rules.Validate(u)
}

// JSONSchema provides a representation of the type for usage in schemas.
func (URI) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:      "string",
		Format:    "uri",
		Title:     "URI",
		MaxLength: &URIMaxLength,
		Description: here.Doc(`
			Uniform Resource Identifier in scheme:opaque form used to identify a
			resource or address, e.g. "gobl:acme.example.com" or
			"iso6523-actorid-upis::9920:x3157928m".
		`),
	}
}

func uriRules() *rules.Set {
	return rules.For(URI(""),
		rules.Assert("01", fmt.Sprintf("uri must be no longer than %d characters", URIMaxLength),
			is.Length(0, int(URIMaxLength)),
		),
		rules.AssertIfPresent("02", "uri must be a valid absolute URI with a scheme",
			is.Func("valid uri", validURI),
		),
	)
}

// validURI returns true for the empty URI (optional fields are skipped,
// matching the rest of the cbc validators) and otherwise requires a
// parseable absolute URI with a non-empty scheme and value.
func validURI(val any) bool {
	u, ok := val.(URI)
	if !ok || u == "" {
		return false
	}
	p, err := url.Parse(string(u))
	if err != nil {
		return false
	}
	return p.Scheme != "" && (p.Opaque != "" || p.Host != "" || p.Path != "")
}
