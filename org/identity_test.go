package org_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestIdentityKeyExtraction(t *testing.T) {
	data := `{"key":"example-key","code":"BAR"}`
	i := new(org.Identity)
	err := json.Unmarshal([]byte(data), i)
	require.NoError(t, err)
	assert.Equal(t, "example-key", i.Key.String())
	assert.Equal(t, "BAR", i.Code.String())
}
