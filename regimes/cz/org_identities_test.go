package cz_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/cz"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestOrgIdentityRules(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		{
			name:     "valid IČO - Škoda Auto",
			identity: &org.Identity{Key: cz.IdentityKeyICO, Code: "00177041"},
		},
		{
			name:     "valid IČO - ČEZ",
			identity: &org.Identity{Key: cz.IdentityKeyICO, Code: "45274649"},
		},
		{
			name:     "valid IČO - Komerční banka",
			identity: &org.Identity{Key: cz.IdentityKeyICO, Code: "45317054"},
		},
		{
			name:     "invalid IČO checksum",
			identity: &org.Identity{Key: cz.IdentityKeyICO, Code: "00177042"},
			err:      "[GOBL-CZ-ORG-IDENTITY-01]",
		},
		{
			name:     "invalid IČO too short",
			identity: &org.Identity{Key: cz.IdentityKeyICO, Code: "0017704"},
			err:      "[GOBL-CZ-ORG-IDENTITY-01]",
		},
		{
			name:     "non-IČO identity ignored",
			identity: &org.Identity{Key: "other", Code: "anything"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.identity)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}
