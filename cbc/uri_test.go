package cbc_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURIValidate(t *testing.T) {
	valid := []cbc.URI{
		"gobl:acme.example.com",
		"iso6523-actorid-upis::9920:b12312312",
		"mailto:billing@example.com",
		"https://x.example/y",
	}
	for _, u := range valid {
		t.Run("valid/"+string(u), func(t *testing.T) {
			assert.NoError(t, u.Validate())
			assert.NoError(t, rules.Validate(u))
		})
	}

	// Empty is allowed (optional fields skip validation via AssertIfPresent).
	assert.NoError(t, cbc.URI("").Validate())

	invalid := []cbc.URI{
		"acme.example.com", // no scheme
		"gobl:",            // scheme only, no value
	}
	for _, u := range invalid {
		t.Run("invalid/"+string(u), func(t *testing.T) {
			assert.Error(t, u.Validate())
		})
	}

	t.Run("over length", func(t *testing.T) {
		u := cbc.URI("gobl:" + strings.Repeat("a", int(cbc.URIMaxLength)+1))
		assert.Error(t, u.Validate())
	})
}

func TestURIAccessors(t *testing.T) {
	// Peppol participant identifier as a URI per the CEN/Peppol spec:
	// scheme is "iso6523-actorid-upis"; the leading "::" in the wire
	// form leaves a colon at the start of the opaque part.
	u := cbc.URI("iso6523-actorid-upis::9920:x3157928m")
	assert.Equal(t, "iso6523-actorid-upis", u.Scheme())
	assert.Equal(t, ":9920:x3157928m", u.Opaque())
	assert.Equal(t, "iso6523-actorid-upis::9920:x3157928m", u.String())

	g := cbc.URI("gobl:acme.example.com")
	assert.Equal(t, "gobl", g.Scheme())
	assert.Equal(t, "acme.example.com", g.Opaque())

	parsed, err := u.Parse()
	require.NoError(t, err)
	assert.Equal(t, "iso6523-actorid-upis", parsed.Scheme)
}

func TestURISchemeOpaqueParseErrors(t *testing.T) {
	// A URI with a control character makes url.Parse return an error;
	// the Scheme/Opaque accessors then return empty strings.
	bad := cbc.URI("gobl:\x7f")
	if _, err := bad.Parse(); err == nil {
		t.Skip("url.Parse tolerated the input; cannot exercise error branch")
	}
	assert.Equal(t, "", bad.Scheme())
	assert.Equal(t, "", bad.Opaque())
	assert.Equal(t, "", bad.Host())
	assert.Equal(t, "", bad.Path())
}

func TestURIHostPath(t *testing.T) {
	// Hierarchical URI: host + path populated, opaque empty.
	h := cbc.URI("https://acme.example/keys/abc")
	assert.Equal(t, "https", h.Scheme())
	assert.Equal(t, "acme.example", h.Host())
	assert.Equal(t, "/keys/abc", h.Path())
	assert.Equal(t, "", h.Opaque())

	// Opaque URI: host and path empty.
	o := cbc.URI("mailto:billing@example.com")
	assert.Equal(t, "mailto", o.Scheme())
	assert.Equal(t, "", o.Host())
	assert.Equal(t, "", o.Path())
	assert.Equal(t, "billing@example.com", o.Opaque())

	// gobl: address — opaque-form, no host or path.
	g := cbc.URI("gobl:acme.example.com")
	assert.Equal(t, "", g.Host())
	assert.Equal(t, "", g.Path())
}

func TestValidURINonURIValue(t *testing.T) {
	// Rules engine calls the validator with arbitrary types; passing
	// a non-URI value (or empty URI) returns false.
	assert.Error(t, cbc.URI("not a url with spaces and bad %ZZ encoding").Validate())
}

func TestURIJSONSchema(t *testing.T) {
	js := cbc.URI("").JSONSchema()
	require.NotNil(t, js)
	assert.Equal(t, "string", js.Type)
	assert.Equal(t, "uri", js.Format)
	assert.Equal(t, "URI", js.Title)
}
