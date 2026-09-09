package dk_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/dk"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityTypeDefinitions(t *testing.T) {
	// Test that the CVR identity type constant is correctly defined
	assert.Equal(t, "CVR", string(dk.IdentityTypeCVR))
	// Test that the CPR identity type constant is correctly defined
	assert.Equal(t, "CPR", string(dk.IdentityTypeCPR))
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
			err: "[GOBL-ORG-IDENTITY-01]",
		},
		{
			name: "too short",
			identity: &org.Identity{
				Type: "CVR",
				Code: "1234567",
			},
			err: "[GOBL-DK-ORG-IDENTITY-01]",
		},
		{
			name: "too long",
			identity: &org.Identity{
				Type: "CVR",
				Code: "123456789",
			},
			err: "[GOBL-DK-ORG-IDENTITY-01]",
		},
		{
			name: "contains letters",
			identity: &org.Identity{
				Type: "CVR",
				Code: "1234567A",
			},
			err: "[GOBL-DK-ORG-IDENTITY-01]",
		},
		{
			name: "bad checksum",
			identity: &org.Identity{
				Type: "CVR",
				Code: "13585627",
			},
			err: "[GOBL-DK-ORG-IDENTITY-01]",
		},
		{
			name: "non-CVR identity",
			identity: &org.Identity{
				Type: "other",
				Code: "invalid",
			},
		},
		{
			name: "valid CPR",
			identity: &org.Identity{
				Type: "CPR",
				Code: "1111111118",
			},
		},
		{
			name: "CPR too short",
			identity: &org.Identity{
				Type: "CPR",
				Code: "123456789",
			},
			err: "[GOBL-DK-ORG-IDENTITY-02]",
		},
		{
			name: "CPR too long",
			identity: &org.Identity{
				Type: "CPR",
				Code: "12345678901",
			},
			err: "[GOBL-DK-ORG-IDENTITY-02]",
		},
		{
			name: "CPR contains letters",
			identity: &org.Identity{
				Type: "CPR",
				Code: "123456789A",
			},
			err: "[GOBL-DK-ORG-IDENTITY-02]",
		},
	}

	opts := []rules.WithContext{
		tax.RegimeContext(dk.CountryCode),
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

func TestNormalizeOrgIdentity(t *testing.T) {
	t.Run("CPR with hyphen", func(t *testing.T) {
		id := &org.Identity{Type: dk.IdentityTypeCPR, Code: "150605-4321"}
		norm.Normalize(id, tax.RegimeContext(dk.CountryCode))
		assert.Equal(t, cbc.Code("1506054321"), id.Code)
	})

	t.Run("unknown type ignored", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "150605-4321"}
		norm.Normalize(id, tax.RegimeContext(dk.CountryCode))
		assert.Equal(t, cbc.Code("150605-4321"), id.Code)
	})
}
