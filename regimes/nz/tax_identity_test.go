package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 8-digit IRD", code: "99114622"},
		{name: "valid 9-digit IRD using secondary weights", code: "136410132"},
		{name: "empty code", code: "", err: ""},
		{name: "invalid check digit", code: "99114623", err: "IDENTITY-02"},
		{name: "below valid range", code: "00000002", err: "IDENTITY-02"},
		{name: "too short", code: "9911462", err: "IDENTITY-01"},
		{name: "too long", code: "9911462200", err: "IDENTITY-01"},
		{name: "non-numeric", code: "9911462A", err: "IDENTITY-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "NZ",
				Code:    tt.code,
			}

			err := rules.Validate(tID)
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
