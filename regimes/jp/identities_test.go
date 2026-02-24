package jp_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeRegistrationNumber(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{name: "already valid", code: "T5010401067252", expected: "T5010401067252"},
		{name: "lowercase t", code: "t5010401067252", expected: "T5010401067252"},
		{name: "with spaces", code: " T5010401067252 ", expected: "T5010401067252"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: jp.IdentityKeyRegistrationNumber, Code: cbc.Code(tt.code)}
			jp.Normalize(id)
			assert.Equal(t, tt.expected, id.Code.String())
		})
	}
}

func TestValidateRegistrationNumber(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  string
	}{
		{name: "valid sony", code: "T5010401067252"},
		{name: "valid nintendo", code: "T1130001011420"},
		{
			name: "missing T prefix",
			code: "5010401067252",
			err:  "code",
		},
		{
			name: "too short",
			code: "T12345",
			err:  "code",
		},
		{
			name: "too long",
			code: "T12345678901234",
			err:  "code",
		},
		{
			name: "empty",
			code: "",
			err:  "code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: jp.IdentityKeyRegistrationNumber, Code: cbc.Code(tt.code)}
			err := jp.Validate(id)
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
