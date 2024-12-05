package in_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/in"
	"github.com/stretchr/testify/assert"
)

func TestOrgIdentityNormalize(t *testing.T) {
	tests := []struct {
		name     string
		typ      cbc.Code
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "PAN already normalized", input: "ABCDE1234F", typ: in.IdentityTypePAN, expected: "ABCDE1234F"},
		{name: "PAN lowercase input", input: "abcde1234f", typ: in.IdentityTypePAN, expected: "ABCDE1234F"},
		{name: "PAN mixed case input", typ: in.IdentityTypePAN, input: "AbCdE1234f", expected: "ABCDE1234F"},
		{name: "PAN extra spaces", typ: in.IdentityTypePAN, input: "  ABCDE1234F  ", expected: "ABCDE1234F"},
		{name: "PAN special characters", typ: in.IdentityTypePAN, input: "AB-CDE1234F", expected: "ABCDE1234F"},

		{name: "HSN already normalized", input: "12345678", typ: in.IdentityTypeHSN, expected: "12345678"},
		{name: "HSN symbols", input: "1234-5678", typ: in.IdentityTypeHSN, expected: "12345678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: in.IdentityTypePAN, Code: tt.input}
			in.Normalize(id)
			assert.Equal(t, tt.expected, id.Code)
		})
	}
}

func TestOrgIdentityValidate(t *testing.T) {
	tests := []struct {
		name string
		typ  cbc.Code
		code cbc.Code
		err  string
	}{
		{name: "not HSN nor PAN", typ: "XYZ", code: "1234"},

		{name: "valid PAN 1", typ: in.IdentityTypePAN, code: "BAJPC4350M"},
		{name: "valid PAN 2", typ: in.IdentityTypePAN, code: "DAJPC4150P"},
		{name: "valid PAN 3", typ: in.IdentityTypePAN, code: "XGZFE7225A"},
		{name: "valid PAN 4", typ: in.IdentityTypePAN, code: "CTUGE1616Y"},

		{name: "valid HSN 1", typ: in.IdentityTypeHSN, code: "1234"},
		{name: "valid HSN 2", typ: in.IdentityTypeHSN, code: "123456"},
		{name: "valid HSN 3", typ: in.IdentityTypeHSN, code: "12345678"},

		{
			name: "PAN too short",
			typ:  in.IdentityTypePAN,
			code: "ABC1234F",
			err:  "code: must be in a valid format.",
		},
		{
			name: "PAN contains spaces",
			typ:  in.IdentityTypePAN,
			code: "ABCDE 1234F",
			err:  "code: must be in a valid format.",
		},
		{
			name: "PAN extra characters",
			typ:  in.IdentityTypePAN,
			code: "ABCDE1234F12",
			err:  "code: must be in a valid format.",
		},
		{
			name: "HSN too short",
			typ:  in.IdentityTypeHSN,
			code: "123",
			err:  "code: must be a 4, 6, or 8 digit number",
		},
		{
			name: "HSN mid",
			typ:  in.IdentityTypeHSN,
			code: "12345",
			err:  "code: must be a 4, 6, or 8 digit number",
		},
		{
			name: "HSN mid 2",
			typ:  in.IdentityTypeHSN,
			code: "1234567",
			err:  "code: must be a 4, 6, or 8 digit number",
		},
		{
			name: "HSN long",
			typ:  in.IdentityTypeHSN,
			code: "123456789",
			err:  "code: must be a 4, 6, or 8 digit number",
		},
		{
			name: "HSN contains letters",
			typ:  in.IdentityTypeHSN,
			code: "1234A6",
			err:  "code: must be a 4, 6, or 8 digit number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: tt.typ, Code: tt.code}
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
