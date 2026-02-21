package hr_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/hr"
	"github.com/stretchr/testify/assert"
)

func TestIdentityTypeDefinitions(t *testing.T) {
	// Test that the OIB identity type constant is correctly defined
	assert.Equal(t, "OIB", string(hr.IdentityTypeOIB))
}

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		{
			name: "valid OIB 1",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "12345678903",
			},
		},
		{
			name: "valid OIB 2",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "11111111119",
			},
		},
		{
			name: "empty code",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "",
			},
		},
		{
			name: "too short",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "1234567890",
			},
			err: "invalid format",
		},
		{
			name: "too long",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "123456789012",
			},
			err: "invalid format",
		},
		{
			name: "contains letters",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "1234567890A",
			},
			err: "invalid format",
		},
		{
			name: "invalid checksum",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "12345678900",
			},
			err: "invalid checksum",
		},
		{
			name: "non-OIB identity type is not validated",
			identity: &org.Identity{
				Type: "OTHER",
				Code: "invalid_code",
			},
		},
		{
			name:     "nil identity",
			identity: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hr.Validate(tt.identity)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
