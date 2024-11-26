package in_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/in"
	"github.com/stretchr/testify/assert"
)

func TestItemValidation(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		err  string
	}{
		{
			name: "valid HSN code",
			item: &org.Item{
				Identities: []*org.Identity{
					{
						Type: "HSN",
						Code: "12345678",
					},
				},
			},
			err: "",
		},
		{
			name: "valid HSN code with 4 digits",
			item: &org.Item{
				Identities: []*org.Identity{
					{
						Type: "HSN",
						Code: "1234",
					},
				},
			},
			err: "",
		},
		{
			name: "invalid HSN code format",
			item: &org.Item{
				Identities: []*org.Identity{
					{
						Type: "HSN",
						Code: "12A456",
					},
				},
			},
			err: "must be a 4, 6, or 8-digit number",
		},
		{
			name: "missing HSN identity",
			item: &org.Item{
				Identities: []*org.Identity{},
			},
			err: "", // No error expected since it's not mandatory in some specific cases
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := in.Validate(tc.item)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tc.err)
				}
			}
		})
	}
}
