package org_test

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
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

func TestIdentityValidate(t *testing.T) {
	t.Run("with basics", func(t *testing.T) {
		id := &org.Identity{
			Code: "BAR",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0004",
			},
		}
		err := id.Validate()
		assert.NoError(t, err)
	})
	t.Run("with both key and type", func(t *testing.T) {
		id := &org.Identity{
			Key:  "fiscal-code",
			Type: "NIF",
			Code: "1234567890",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "type: must be empty when key is set")
	})
}

func TestIdentitySetValidators(t *testing.T) {
	t.Run("ignore other types", func(t *testing.T) {
		var idents *org.Identity
		err := validation.Validate(idents, org.RequireIdentityType("FOO"))
		assert.NoError(t, err)
	})
	t.Run("require identity type", func(t *testing.T) {
		idents := []*org.Identity{
			{
				Type: "BAR",
				Code: "FOO",
			},
		}
		err := validation.Validate(idents, org.RequireIdentityType("BAR"))
		assert.NoError(t, err)

		err = validation.Validate(idents, org.RequireIdentityType("FOO"))
		assert.ErrorContains(t, err, "missing type 'FOO'")
	})

	t.Run("require identity key", func(t *testing.T) {
		idents := []*org.Identity{
			{
				Type: "BAR",
				Code: "FOO",
			},
			{
				Key:  "fiscal-code",
				Code: "12345",
			},
		}
		err := validation.Validate(idents, org.RequireIdentityKey("fiscal-code"))
		assert.NoError(t, err)

		err = validation.Validate(idents, org.RequireIdentityKey("code"))
		assert.ErrorContains(t, err, "missing key 'code'")
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
