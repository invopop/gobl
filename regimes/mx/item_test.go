package mx_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/stretchr/testify/assert"
)

func TestItemValidation(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		err  string
	}{
		{
			name: "valid item",
			item: &org.Item{
				Unit: "kg",
				Identities: []*org.Identity{
					{Code: "12345678", Type: "SAT"},
				},
			},
		},
		{
			name: "missing unit",
			item: &org.Item{
				Identities: []*org.Identity{
					{Code: "12345678", Type: "SAT"},
				},
			},
			err: "unit: cannot be blank",
		},
		{
			name: "missing identities",
			item: &org.Item{
				Unit: "kg",
			},
			err: "identities: SAT code must be present",
		},
		{
			name: "empty identities",
			item: &org.Item{
				Unit:       "kg",
				Identities: []*org.Identity{},
			},
			err: "identities: SAT code must be present",
		},
		{
			name: "missing SAT identity",
			item: &org.Item{
				Unit: "kg",
				Identities: []*org.Identity{
					{Type: "GTIN", Code: "12345678"},
				},
			},
			err: "identities: SAT code must be present",
		},
		{
			name: "SAT in invalid format",
			item: &org.Item{
				Unit: "kg",
				Identities: []*org.Identity{
					{Type: "SAT", Code: "ABC2"},
				},
			},
			err: "identities: SAT code must have 8 digits",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := mx.Validate(ts.item)
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}
