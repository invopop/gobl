package nl_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/nl"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
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
			err: "[GOBL-NL-ORG-IDENTITY-01] ($.code) identity code for type KVK must be valid",
		},
		{
			name: "KVK too long",
			identity: &org.Identity{
				Type: nl.IdentityTypeKVK,
				Code: "123456789",
			},
			err: "[GOBL-NL-ORG-IDENTITY-01]",
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
			err: "[GOBL-NL-ORG-IDENTITY-02] ($.code) identity code for type OIN must be valid",
		},
		{
			name: "OIN invalid register code 11",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000011123456789000",
			},
			err: "[GOBL-NL-ORG-IDENTITY-02]",
		},
		{
			name: "OIN invalid prefix",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "12345601123456789000",
			},
			err: "[GOBL-NL-ORG-IDENTITY-02]",
		},
		{
			name: "OIN invalid suffix",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000001123456789123",
			},
			err: "[GOBL-NL-ORG-IDENTITY-02]",
		},
		{
			name: "OIN non-numeric",
			identity: &org.Identity{
				Type: nl.IdentityTypeOIN,
				Code: "00000001ABCDEFGHI000",
			},
			err: "[GOBL-NL-ORG-IDENTITY-02]",
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

	opts := []rules.WithContext{
		tax.RegimeContext(nl.CountryCode),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.identity, opts...)
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
