package in_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/in"
	"github.com/stretchr/testify/assert"
)

func TestNormalizePAN(t *testing.T) {
	tests := []struct {
		name     string
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "already normalized", input: "ABCDE1234F", expected: "ABCDE1234F"},
		{name: "lowercase input", input: "abcde1234f", expected: "ABCDE1234F"},
		{name: "mixed case input", input: "AbCdE1234f", expected: "ABCDE1234F"},
		{name: "extra spaces", input: "  ABCDE1234F  ", expected: "ABCDE1234F"},
		{name: "special characters", input: "AB-CDE1234F", expected: "ABCDE1234F"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: in.IdentityKeyPAN, Code: tt.input}
			in.Normalize(id)
			assert.Equal(t, tt.expected, id.Code)
		})
	}
}

func TestValidatePAN(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid PAN 1", code: "BAJPC4350M"},
		{name: "valid PAN 2", code: "DAJPC4150P"},
		{name: "valid PAN 3", code: "XGZFE7225A"},
		{name: "valid PAN 4", code: "CTUGE1616Y"},

		{
			name: "too short",
			code: "ABC1234F",
			err:  "code: must be in a valid format.",
		},
		{
			name: "contains spaces",
			code: "ABCDE 1234F",
			err:  "code: must be in a valid format.",
		},
		{
			name: "extra characters",
			code: "ABCDE1234F12",
			err:  "code: must be in a valid format.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: in.IdentityKeyPAN, Code: tt.code}
			err := in.Validate(id)

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
