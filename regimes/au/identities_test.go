package au_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/au"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "ACN good 1", code: "010499966"},
		{name: "ACN good 2", code: "813283831"},
		{name: "ACN good 3", code: "419673715"},
		{
			name: "zeros",
			code: "00000000000",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "123456789012",
			err:  "invalid format",
		},
		{
			name: "too short",
			code: "12345678",
			err:  "invalid format",
		},
		{
			name: "not normalized",
			code: "12.449.95-4",
			err:  "invalid format",
		},
		{
			name: "letter",
			code: "12a44939544",
			err:  "invalid format",
		},
		{
			name: "bad ACN checksum 1",
			code: "419673716",
			err:  "checksum mismatch",
		},
		{
			name: "bad ACN checksum 2",
			code: "678381888",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &org.Identity{Key: au.IdentityCompanyNumber, Code: tt.code}
			err := au.Validate(tID)
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
