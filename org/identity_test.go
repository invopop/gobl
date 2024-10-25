package org_test

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestAddIdentity(t *testing.T) {
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
}

func TestIdentityNormalize(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		var id *org.Identity
		assert.NotPanics(t, func() {
			id.Normalize(nil)
		})
	})
	t.Run("missing extensions", func(t *testing.T) {
		id := &org.Identity{
			Type: cbc.Code("FOO"),
			Code: "BAR",
			Ext:  tax.Extensions{},
		}
		id.Normalize(nil)
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
		id.Normalize(nil)
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
