package gr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tID *tax.Identity
		assert.NotPanics(t, func() {
			gr.Normalize(tID)
		})
	})

	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{name: "no change", code: "123456789", want: "123456789"},
		{name: "remove spaces", code: " 123 456 789 ", want: "123456789"},
		{name: "remove dashes", code: "123-456-789", want: "123456789"},
		{name: "remove dots", code: "123.456.789", want: "123456789"},
		{name: "remove slashes", code: "123/456/789", want: "123456789"},
		{name: "remove EL", code: "EL123456789", want: "123456789"},
		{name: "remove GR", code: "GR123456789", want: "123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "EL", Code: tt.code}
			gr.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code)
		})
	}

	tID := &tax.Identity{Country: "GR", Code: "123456789"}
	gr.Normalize(tID)
	assert.Equal(t, "EL123456789", tID.String())
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "older than 1999", code: "064677095"},
		{name: "legal entity", code: "925667500"},
		{name: "natural person", code: "320806520"},
		{name: "digit repetition", code: "691063615"},
		{
			name: "too long",
			code: "1234567890123",
			err:  "invalid format",
		},
		{
			name: "too short",
			code: "123456",
			err:  "invalid format",
		},
		{
			name: "non-numeric",
			code: "12345678A",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "925667503",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "EL", Code: tt.code}
			err := gr.Validate(tID)
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
