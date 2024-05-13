package be_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/be"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good enterprise 1", code: "0413172884"},
		{name: "good enterprise 2", code: "0414445663"},
		{name: "good enterprise 3", code: "0897223571"},
		{name: "good personal 1", code: "897222383"},
		{name: "good personal 2", code: "897231984"},
		{
			name: "zeros",
			code: "000000000",
			err:  "invalid format",
		},
		{
			name: "bad mid length",
			code: "82238333",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "01234567890123",
			err:  "invalid format",
		},
		{
			name: "too short",
			code: "0123456",
			err:  "invalid format",
		},
		{
			name: "not normalized",
			code: "0.413.172-884",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "0413172885",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.BE, Code: tt.code}
			err := be.Validate(tID)
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
