package hr_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/hr"
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
			err: "[GOBL-ORG-IDENTITY-01]",
		},
		{
			name: "too short",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "1234567890",
			},
			err: "[GOBL-HR-ORG-IDENTITY-01]",
		},
		{
			name: "too long",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "123456789012",
			},
			err: "[GOBL-HR-ORG-IDENTITY-01]",
		},
		{
			name: "contains letters",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "1234567890A",
			},
			err: "[GOBL-HR-ORG-IDENTITY-01]",
		},
		{
			name: "invalid checksum",
			identity: &org.Identity{
				Type: hr.IdentityTypeOIB,
				Code: "12345678900",
			},
			err: "[GOBL-HR-ORG-IDENTITY-02]",
		},
		{
			name: "non-OIB identity type is not validated",
			identity: &org.Identity{
				Type: "OTHER",
				Code: "invalid_code",
			},
		},
	}

	opts := []rules.WithContext{
		tax.RegimeContext(hr.CountryCode),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.identity, opts...)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}
