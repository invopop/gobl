package au_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		code        cbc.Code
		expectedErr string
	}{
		{name: "valid ABN from ATO", code: "51824753556"},
		{name: "valid ABN with spaces", code: "51 824 753 556"},
		{name: "second valid ABN", code: "53004085616"},
		{name: "empty ABN", code: ""},
		{name: "too short", code: "1234567890", expectedErr: "invalid length"},
		{name: "too long", code: "123456789012", expectedErr: "invalid length"},
		{name: "invalid checksum", code: "11111111111", expectedErr: "invalid checksum"},
		{name: "contains non-numeric characters", code: "5182475355A", expectedErr: "invalid characters, expected numeric"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AU", Code: tt.code}

			err := au.Validate(tID)

			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		var tID *tax.Identity
		assert.NoError(t, au.Validate(tID))
	})
}

func TestNormalizeTaxIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "spaces stripped", input: "51 824 753 556", expected: "51824753556"},
		{name: "country prefix stripped", input: "AU51824753556", expected: "51824753556"},
		{name: "already normalized", input: "51824753556", expected: "51824753556"},
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AU", Code: tt.input}

			au.Normalize(tID)

			assert.Equal(t, tt.expected, tID.Code)
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var tID *tax.Identity
			au.Normalize(tID)
		})
	})
}
