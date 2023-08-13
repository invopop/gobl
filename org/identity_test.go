package org_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/stretchr/testify/assert"
)

func TestIdentityValidation(t *testing.T) {
	// Use Mexico for these tests
	ctx := mx.New().WithContext(context.Background())
	id := &org.Identity{
		Type: "FOO",
		Code: "WAT",
	}
	err := id.ValidateWithContext(ctx)
	assert.NoError(t, err)

	id = &org.Identity{
		Type: "FOO",
		Key:  mx.IdentityKeyCFDIUse,
		Code: "G01",
	}
	err = id.ValidateWithContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type: must be blank when key set.")

	id = &org.Identity{
		Type: "FOO",
		Code: "", // empty
	}
	err = id.ValidateWithContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code: cannot be blank.")

	id = &org.Identity{
		Key:  mx.IdentityKeyCFDIFiscalRegime,
		Code: "606",
	}
	err = id.ValidateWithContext(ctx)
	assert.NoError(t, err)

	id.Code = "500"
	err = id.ValidateWithContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code: invalid mx-cfdi-fiscal-regime.")
}

func TestAddIdentity(t *testing.T) {
	foo := cbc.Key("foo")
	st := struct {
		Identities []*org.Identity
	}{
		Identities: []*org.Identity{
			{
				Key:  foo,
				Code: "BAR",
			},
		},
	}
	st.Identities = org.AddIdentity(st.Identities, &org.Identity{
		Key:  foo,
		Code: "BARDOM",
	})
	assert.Len(t, st.Identities, 1)
	assert.Equal(t, "BARDOM", st.Identities[0].Code.String())
}
