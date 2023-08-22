package org_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
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
