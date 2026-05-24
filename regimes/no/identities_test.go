package no_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/no"
	"github.com/stretchr/testify/assert"
)

func TestIdentityTypeDefinitions(t *testing.T) {
	assert.Equal(t, "ORG", string(no.IdentityTypeORG))
}

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		{
			name: "valid ORG 1",
			identity: &org.Identity{
				Type: "ORG",
				Code: "923609016",
			},
		},
		{
			name: "valid ORG 2",
			identity: &org.Identity{
				Type: "ORG",
				Code: "982463718",
			},
		},
		{
			name: "valid ORG 3",
			identity: &org.Identity{
				Type: "ORG",
				Code: "889640782",
			},
		},
		{
			name: "empty code",
			identity: &org.Identity{
				Type: "ORG",
				Code: "",
			},
		},
		{
			name: "too short",
			identity: &org.Identity{
				Type: "ORG",
				Code: "12345678",
			},
			err: "invalid format",
		},
		{
			name: "too long",
			identity: &org.Identity{
				Type: "ORG",
				Code: "1234567890",
			},
			err: "invalid format",
		},
		{
			name: "contains letters",
			identity: &org.Identity{
				Type: "ORG",
				Code: "12345678A",
			},
			err: "invalid format",
		},
		{
			name: "bad checksum",
			identity: &org.Identity{
				Type: "ORG",
				Code: "923609017",
			},
			err: "checksum mismatch",
		},
		{
			name: "non-ORG identity",
			identity: &org.Identity{
				Type: "other",
				Code: "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := no.Validate(tt.identity)
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
