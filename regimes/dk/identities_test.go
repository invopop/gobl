package dk_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/dk"
	"github.com/stretchr/testify/assert"
)

func TestIdentityTypeDefinitions(t *testing.T) {
	// Test that the CVR identity type constant is correctly defined
	assert.Equal(t, "CVR", string(dk.IdentityTypeCVR))
}

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		{
			name: "valid CVR 1",
			identity: &org.Identity{
				Type: "CVR",
				Code: "13585628",
			},
		},
		{
			name: "valid CVR 2",
			identity: &org.Identity{
				Type: "CVR",
				Code: "88146328",
			},
		},
		{
			name: "valid CVR 3",
			identity: &org.Identity{
				Type: "CVR",
				Code: "25063864",
			},
		},
		{
			name: "empty code",
			identity: &org.Identity{
				Type: "CVR",
				Code: "",
			},
			err: "",
		},
		{
			name: "too short",
			identity: &org.Identity{
				Type: "CVR",
				Code: "1234567",
			},
			err: "invalid format",
		},
		{
			name: "too long",
			identity: &org.Identity{
				Type: "CVR",
				Code: "123456789",
			},
			err: "invalid format",
		},
		{
			name: "contains letters",
			identity: &org.Identity{
				Type: "CVR",
				Code: "1234567A",
			},
			err: "invalid format",
		},
		{
			name: "bad checksum",
			identity: &org.Identity{
				Type: "CVR",
				Code: "13585627",
			},
			err: "checksum mismatch",
		},
		{
			name: "non-CVR identity",
			identity: &org.Identity{
				Type: "other",
				Code: "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dk.Validate(tt.identity)
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
