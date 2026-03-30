package nl_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/nl"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		// KVK tests
		{
			name: "valid KVK",
			identity: &org.Identity{
				Type: nl.IdentityTypeKVK,
				Code: "12345678",
			},
		},
		{
			name: "KVK too short",
			identity: &org.Identity{
				Type: nl.IdentityTypeKVK,
				Code: "1234567",
			},
			err: "code: the length must be exactly 8",
		},
		{
			name: "KVK too long",
			identity: &org.Identity{
				Type: nl.IdentityTypeKVK,
				Code: "123456789",
			},
			err: "code: the length must be exactly 8",
		},
		// OIN tests
		{
			name: "valid OIN with register 01",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000001123456789000",
			},
		},
		{
			name: "valid OIN with register 10",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000010123456789000",
			},
		},
		{
			name: "valid OIN with register 99",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000099123456789000",
			},
		},
		{
			name: "OIN invalid register code 00",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000000123456789000",
			},
			err: "must be in a valid format",
		},
		{
			name: "OIN invalid register code 11",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000011123456789000",
			},
			err: "must be in a valid format",
		},
		{
			name: "OIN invalid prefix",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "12345601123456789000",
			},
			err: "must be in a valid format",
		},
		{
			name: "OIN invalid suffix",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000001123456789123",
			},
			err: "must be in a valid format",
		},
		{
			name: "OIN non-numeric",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000001ABCDEFGHI000",
			},
			err: "must be in a valid format",
		},
		// Other identity type
		{
			name: "non-NL identity",
			identity: &org.Identity{
				Type: "other",
				Code: "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := nl.Validate(tt.identity)
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
