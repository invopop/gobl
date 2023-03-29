package gb_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/gb"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "000472631"},
		{name: "good 2", code: "844281425"},
		{name: "good 3", code: "350983637"},
		{name: "good 7", code: "100000132"}, // 100 start
		{name: "good 4", code: "GD103"},
		{name: "good 5", code: "HA503"},
		{name: "good 6", code: "350983637001"}, // last 3 ignored
		{name: "good 6", code: "350983637002"},
		{
			name: "zeros",
			code: "000000000",
			err:  "zeros",
		},
		{
			name: "bad mid length",
			code: "12345678910",
			err:  "invalid format",
		},
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
			name: "not normalized",
			code: "12.449.965-4",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "999999991",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.GB, Code: tt.code}
			err := gb.Validate(tID)
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
