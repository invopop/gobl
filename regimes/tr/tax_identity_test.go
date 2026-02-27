package tr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/tr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		code     cbc.Code
		expected cbc.Code
	}{
		{code: "1234567890", expected: "1234567890"},
		{code: "TR1234567890", expected: "1234567890"},
		{code: "123 456 7890", expected: "1234567890"},
	}
	for _, tt := range tests {
		tID := &tax.Identity{Country: "TR", Code: tt.code}
		tr.Normalize(tID)
		assert.Equal(t, tt.expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid VKN", code: "4540536920"},
		{name: "empty code", code: ""},
		{name: "too short", code: "123456789", err: "invalid format"},
		{name: "too long", code: "12345678901", err: "invalid format"},
		{name: "non-numeric", code: "123456789A", err: "invalid format"},
		{name: "VKN wrong checksum", code: "1234567891", err: "invalid check digit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "TR", Code: tt.code}
			err := tr.Validate(tID)
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
