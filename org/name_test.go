package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/require"
)

func TestNameValidation(t *testing.T) {
	tests := []struct {
		name string
		n    *org.Name
		err  string
	}{
		{
			name: "empty",
			n:    &org.Name{},

			err: "given: cannot be blank; surname: cannot be blank.",
		},
		{
			name: "given",
			n: &org.Name{
				Given: "John",
			},
		},
		{
			name: "surname",
			n: &org.Name{
				Surname: "Doe",
			},
		},
		{
			name: "both",
			n: &org.Name{
				Given:   "John",
				Surname: "Doe",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.Validate()
			if tt.err == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.err)
		})
	}
}
