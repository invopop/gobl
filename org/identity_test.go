package org_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddIdentity(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		foo := cbc.Code("FOO")
		st := struct {
			Identities []*org.Identity
		}{
			Identities: []*org.Identity{
				{
					Type: foo,
					Code: "BAR",
				},
			},
		}
		st.Identities = org.AddIdentity(st.Identities, &org.Identity{
			Type: foo,
			Code: "BARDOM",
		})
		assert.Len(t, st.Identities, 1)
		assert.Equal(t, "BARDOM", st.Identities[0].Code.String())
	})
	t.Run("nil array", func(t *testing.T) {
		var st struct {
			Identities []*org.Identity
		}
		st.Identities = org.AddIdentity(st.Identities, &org.Identity{
			Type: "FOO",
			Code: "BAR",
		})
		assert.Len(t, st.Identities, 1)
		assert.Equal(t, "BAR", st.Identities[0].Code.String())
	})
	t.Run("append type", func(t *testing.T) {
		foo := cbc.Code("FOO")
		st := struct {
			Identities []*org.Identity
		}{
			Identities: []*org.Identity{
				{
					Type: foo,
					Code: "BAR",
				},
			},
		}
		st.Identities = org.AddIdentity(st.Identities, &org.Identity{
			Type: "FOO2",
			Code: "BARDOM",
		})
		assert.Len(t, st.Identities, 2)
		assert.Equal(t, "BARDOM", st.Identities[1].Code.String())
	})
	t.Run("append key", func(t *testing.T) {
		st := struct {
			Identities []*org.Identity
		}{
			Identities: []*org.Identity{
				{
					Key:  "foo",
					Code: "BAR",
				},
			},
		}
		st.Identities = org.AddIdentity(st.Identities, &org.Identity{
			Key:  "foo-second",
			Code: "BARDOM",
		})
		assert.Len(t, st.Identities, 2)
		assert.Equal(t, "BARDOM", st.Identities[1].Code.String())
	})
}

func TestIdentityNormalize(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		var id *org.Identity
		assert.NotPanics(t, func() {
			id.Normalize()
		})
	})
	t.Run("missing extensions", func(t *testing.T) {
		id := &org.Identity{
			Type: cbc.Code("FOO"),
			Code: "BAR",
			Ext:  tax.Extensions{},
		}
		id.Normalize()
		assert.Equal(t, "FOO", id.Type.String())
		assert.Nil(t, id.Ext)
	})
	t.Run("with extension", func(t *testing.T) {
		id := &org.Identity{
			Code: "BAR",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0004",
			},
		}
		id.Normalize()
		assert.Equal(t, "BAR", id.Code.String())
		assert.Equal(t, "0004", id.Ext[iso.ExtKeySchemeID].String())
	})
}

func TestIdentityRules(t *testing.T) {
	t.Run("with basics", func(t *testing.T) {
		id := &org.Identity{
			Code: "BAR",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0004",
			},
		}
		assert.NoError(t, rules.Validate(id))
	})
	t.Run("missing code", func(t *testing.T) {
		id := &org.Identity{}
		faults := rules.Validate(id)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-ORG-IDENTITY-01"))
		assert.Contains(t, faults.Error(), "identity code must be provided")
	})
	t.Run("with both key and type", func(t *testing.T) {
		id := &org.Identity{
			Key:  "fiscal-code",
			Type: "NIF",
			Code: "1234567890",
		}
		faults := rules.Validate(id)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-ORG-IDENTITY-03"))
		assert.Contains(t, faults.Error(), "identity must have either a key or type defined, but not both")
	})
	t.Run("with valid tax scope", func(t *testing.T) {
		id := &org.Identity{
			Scope: org.IdentityScopeTax,
			Code:  "1234567890",
		}
		assert.NoError(t, rules.Validate(id))
	})
	t.Run("with valid legal scope", func(t *testing.T) {
		id := &org.Identity{
			Scope: org.IdentityScopeLegal,
			Code:  "1234567890",
		}
		assert.NoError(t, rules.Validate(id))
	})
	t.Run("with invalid scope", func(t *testing.T) {
		id := &org.Identity{
			Scope: "INVALID",
			Code:  "1234567890",
		}
		faults := rules.Validate(id)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-ORG-IDENTITY-02"))
		assert.Contains(t, faults.Error(), "identity scope when provided must be either 'tax' or 'legal'")
	})
	t.Run("with no scope", func(t *testing.T) {
		id := &org.Identity{
			Code: "1234567890",
		}
		assert.NoError(t, rules.Validate(id))
	})
	t.Run("with key only", func(t *testing.T) {
		id := &org.Identity{
			Key:  "fiscal-code",
			Code: "1234567890",
		}
		assert.NoError(t, rules.Validate(id))
	})
	t.Run("with type only", func(t *testing.T) {
		id := &org.Identity{
			Type: "NIF",
			Code: "1234567890",
		}
		assert.NoError(t, rules.Validate(id))
	})
}

func TestIdentityTests(t *testing.T) {
	idents := []*org.Identity{
		{Type: "BAR", Code: "FOO"},
		{Key: "fiscal-code", Code: "12345"},
	}

	t.Run("IdentitiesTypeIn matches", func(t *testing.T) {
		assert.True(t, org.IdentitiesTypeIn("BAR").Check(idents))
	})
	t.Run("IdentitiesTypeIn no match", func(t *testing.T) {
		assert.False(t, org.IdentitiesTypeIn("FOO").Check(idents))
	})
	t.Run("IdentitiesTypeIn multiple types", func(t *testing.T) {
		assert.True(t, org.IdentitiesTypeIn("FOO", "BAR").Check(idents))
		assert.False(t, org.IdentitiesTypeIn("FOO", "FUZ").Check(idents))
	})
	t.Run("IdentitiesTypeIn string", func(t *testing.T) {
		assert.Equal(t, "has a type in [BAR, FOO]", org.IdentitiesTypeIn("BAR", "FOO").String())
	})

	t.Run("IdentitiesKeyIn matches", func(t *testing.T) {
		assert.True(t, org.IdentitiesKeyIn("fiscal-code").Check(idents))
	})
	t.Run("IdentitiesKeyIn no match", func(t *testing.T) {
		assert.False(t, org.IdentitiesKeyIn("missing-key").Check(idents))
	})
	t.Run("IdentitiesKeyIn multiple keys", func(t *testing.T) {
		assert.True(t, org.IdentitiesKeyIn("missing", "fiscal-code").Check(idents))
		assert.False(t, org.IdentitiesKeyIn("one", "two").Check(idents))
	})
	t.Run("IdentitiesKeyIn string", func(t *testing.T) {
		assert.Equal(t, "has a key in [fiscal-code]", org.IdentitiesKeyIn("fiscal-code").String())
	})

	t.Run("IdentityTypeIn matches single", func(t *testing.T) {
		id := &org.Identity{Type: "BAR", Code: "FOO"}
		assert.True(t, org.IdentityTypeIn("BAR").Check(id))
		assert.False(t, org.IdentityTypeIn("FOO").Check(id))
	})
	t.Run("IdentityTypeIn string", func(t *testing.T) {
		assert.Equal(t, "type in [BAR]", org.IdentityTypeIn("BAR").String())
	})

	t.Run("IdentityKeyIn matches single", func(t *testing.T) {
		id := &org.Identity{Key: "fiscal-code", Code: "12345"}
		assert.True(t, org.IdentityKeyIn("fiscal-code").Check(id))
		assert.False(t, org.IdentityKeyIn("other").Check(id))
	})
	t.Run("IdentityKeyIn string", func(t *testing.T) {
		assert.Equal(t, "key in [fiscal-code]", org.IdentityKeyIn("fiscal-code").String())
	})

	t.Run("non-identity type returns false", func(t *testing.T) {
		assert.False(t, org.IdentitiesTypeIn("FOO").Check("not-an-identity"))
		assert.False(t, org.IdentitiesTypeIn("FOO").Check(nil))
	})
}

func TestIdentityForType(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		idents := []*org.Identity{
			{
				Type: "BAR",
				Code: "FOO",
			},
			{
				Type: "FOO",
				Code: "BAR",
			},
		}
		id := org.IdentityForType(idents, "FOO")
		assert.Equal(t, "BAR", id.Code.String())
		assert.Nil(t, org.IdentityForType(idents, "BAZ"))
	})
}

func TestIdentityForKey(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		idents := []*org.Identity{
			{
				Key:  "bar",
				Code: "FOO",
			},
			{
				Key:  "foo",
				Code: "BAR",
			},
		}
		id := org.IdentityForKey(idents, "foo")
		assert.Equal(t, "BAR", id.Code.String())
		assert.Nil(t, org.IdentityForKey(idents, "baz"))
	})
}

func TestIdentityForExtKey(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		idents := []*org.Identity{
			{
				Ext: tax.Extensions{
					cbc.Key("foo"): "bar",
				},
			},
			{
				Ext: tax.Extensions{
					cbc.Key("baz"): "qux",
				},
			},
		}
		id := org.IdentityForExtKey(idents, "baz")
		assert.Equal(t, "qux", id.Ext["baz"].String())
		assert.Nil(t, org.IdentityForExtKey(idents, "nonexistent"))
	})
	t.Run("nil extensions", func(t *testing.T) {
		idents := []*org.Identity{
			{
				Code: "1234",
			},
			{
				Code: "5678",
				Ext: tax.Extensions{
					cbc.Key("baz"): "qux",
				},
			},
		}
		id := org.IdentityForExtKey(idents, "baz")
		assert.Equal(t, "qux", id.Ext["baz"].String())
		assert.Nil(t, org.IdentityForExtKey(idents, "nonexistent"))
	})
}

func TestIdentityJSONSchema(t *testing.T) {
	base := here.Doc(`
		{
			"properties": {
				"scope": {
					"$ref": "https://gobl.org/draft-0/cbc/key",
					"title": "Scope"
				}
			}
		}
	`)
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(base), js))
	org.Identity{}.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("scope")
	assert.True(t, ok)
	assert.Len(t, prop.OneOf, 2)
	assert.Equal(t, org.IdentityScopeTax, prop.OneOf[0].Const)
	assert.Equal(t, "Tax", prop.OneOf[0].Title)
	assert.Equal(t, org.IdentityScopeLegal, prop.OneOf[1].Const)
	assert.Equal(t, "Legal", prop.OneOf[1].Title)
}
