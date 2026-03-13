package sg_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "company GST", code: "M91234567X"},
		{name: "sole proprietorship GST", code: "MR2345678A"},
		{name: "overseas vendor GST", code: "MB2345678A"},
		{name: "overseas vendor GST 2", code: "MX2345678A"},
		{name: "invalid GST short", code: "M91234567", err: "IDENTITY-01"},
		{name: "invalid GST long", code: "M91234567XA", err: "IDENTITY-01"},
		{name: "invalid GST no M", code: "912345678X", err: "IDENTITY-01"},
		{name: "invalid GST no end letter", code: "M912345678", err: "IDENTITY-01"},
		// UEN identities also
		{name: "UEN (ROC)", code: "199912345A"},
		{name: "UEN (ROB)", code: "12345678A"},
		{name: "UEN (Others)", code: "T12AB1234A"},
		{name: "NIRC/FIN", code: "S1234567A", err: "IDENTITY-01"},
		{name: "Invalid short", code: "1234567A", err: "IDENTITY-01"},
		{name: "Invalid UEN (ROC)", code: "2199123456", err: "IDENTITY-01"},
		{name: "Invalid UEN (ROB)", code: "1234567A", err: "IDENTITY-01"},
		{name: "Invalid UEN (Others)", code: "T12A1234A", err: "IDENTITY-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "SG",
				Code:    tt.code,
			}
			err := rules.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		var id *tax.Identity
		err := rules.Validate(id)
		assert.NoError(t, err)
	})
}
