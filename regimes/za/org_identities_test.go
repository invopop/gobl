package za_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/za"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityTypeDefinitions(t *testing.T) {
	assert.Equal(t, "CIPC", string(za.IdentityTypeCIPC))
}

func TestCIPCIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid private company", code: "2020/123456/07"},
		{name: "valid public company", code: "1998/654321/06"},
		{name: "valid close corporation", code: "2005/123456/23"},
		{name: "too short year", code: "202/123456/07", err: "[GOBL-ZA-ORG-IDENTITY-01]"},
		{name: "too short sequence", code: "2020/12345/07", err: "[GOBL-ZA-ORG-IDENTITY-01]"},
		{name: "missing suffix", code: "2020/123456", err: "[GOBL-ZA-ORG-IDENTITY-01]"},
		{name: "wrong separators", code: "2020-123456-07", err: "[GOBL-ZA-ORG-IDENTITY-01]"},
		{name: "non-numeric", code: "2020/12345A/07", err: "[GOBL-ZA-ORG-IDENTITY-01]"},
	}

	ctx := tax.RegimeContext(za.CountryCode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: za.IdentityTypeCIPC, Code: tt.code}
			err := rules.Validate(id, ctx)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestNonCIPCIdentityIgnored(t *testing.T) {
	ctx := tax.RegimeContext(za.CountryCode)
	id := &org.Identity{Type: "other", Code: "not-a-cipc-number"}
	assert.NoError(t, rules.Validate(id, ctx))
}
