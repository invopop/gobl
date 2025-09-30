package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNameNormalize(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var n *org.Name
		assert.NotPanics(t, func() {
			n.Normalize()
		})
	})
	n := &org.Name{
		Alias:    "  john doe  ",
		Prefix:   "  Mr.  ",
		Given:    "  John  ",
		Middle:   "  Quincy  ",
		Surname:  "  Doe  ",
		Surname2: "  Smith  ",
		Suffix:   "  Jr.  ",
	}
	n.Normalize()
	require.Equal(t, "john doe", n.Alias)
	require.Equal(t, "Mr.", n.Prefix)
	require.Equal(t, "John", n.Given)
	require.Equal(t, "Quincy", n.Middle)
	require.Equal(t, "Doe", n.Surname)
	require.Equal(t, "Smith", n.Surname2)
	require.Equal(t, "Jr.", n.Suffix)
}

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
