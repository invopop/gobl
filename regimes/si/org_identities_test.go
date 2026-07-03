package si_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/si"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateOrgIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		{
			// Krka d.d., Novo mesto's publicly listed registration number.
			name: "valid",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "5043611000",
			},
		},
		{
			name: "valid 2",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "1234579000",
			},
		},
		{
			// Weighted sum leaves a remainder of 1, so the check digit is 0.
			name: "valid with zero check digit",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "1234510000",
			},
		},
		{
			name: "too short",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "5043611",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "non numeric suffix",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "1234579A00",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "checksum mismatch",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "1234571000",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			// The weighted sum leaves a remainder of 0, which is never
			// assigned: no valid registration number produces it.
			name: "remainder zero",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "1234567000",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "empty code",
			identity: &org.Identity{
				Type: si.IdentityTypeMaticna,
				Code: "",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "other identity type",
			identity: &org.Identity{
				Type: "OTHER",
				Code: "not checked",
			},
		},
	}

	opts := []rules.WithContext{
		tax.RegimeContext(si.CountryCode),
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
