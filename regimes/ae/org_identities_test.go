package ae_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ae"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTINValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid", code: "1234567897"},
		{name: "first 10 digits of a TRN", code: "1345678906"},

		{name: "too short", code: "123456789", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "too long", code: "12345678901", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "a full TRN", code: "134567890123003", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "non-numeric", code: "12345ABCDE", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "not starting with 1", code: "2234567895", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "wrong check digit", code: "1234567890", err: "[GOBL-AE-ORG-IDENTITY-02]"},
	}

	ctx := tax.RegimeContext(ae.CountryCode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: ae.IdentityTypeTIN, Code: tt.code}
			err := rules.Validate(id, ctx)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestEmiratesIDValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid", code: "784199012345671"},

		{name: "too short", code: "78419901234567", err: "[GOBL-AE-ORG-IDENTITY-03]"},
		{name: "too long", code: "7841990123456712", err: "[GOBL-AE-ORG-IDENTITY-03]"},
		{name: "non-numeric", code: "784ABCD12345671", err: "[GOBL-AE-ORG-IDENTITY-03]"},
		{name: "not starting with 784", code: "123199012345671", err: "[GOBL-AE-ORG-IDENTITY-03]"},
	}

	ctx := tax.RegimeContext(ae.CountryCode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: ae.IdentityTypeEmiratesID, Code: tt.code}
			err := rules.Validate(id, ctx)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}
