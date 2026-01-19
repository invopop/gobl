package sg_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  bool
	}{
		{name: "company GST", code: "M91234567X", err: false},
		{name: "sole proprietorship GST", code: "MR2345678A", err: false},
		{name: "overseas vendor GST", code: "MB2345678A", err: false},
		{name: "overseas vendor GST 2", code: "MX2345678A", err: false},
		{name: "invalid GST short", code: "M91234567", err: true},
		{name: "invalid GST long", code: "M91234567XA", err: true},
		{name: "invalid GST no M", code: "912345678X", err: true},
		{name: "invalid GST no end letter", code: "M912345678", err: true},
		// UEN identities also
		{name: "UEN (ROC)", code: "199912345A", err: false},
		{name: "UEN (ROB)", code: "12345678A", err: false},
		{name: "UEN (Others)", code: "T12AB1234A", err: false},
		{name: "NIRC/FIN", code: "S1234567A", err: true},
		{name: "Invalid short", code: "1234567A", err: true},
		{name: "Invalid UEN (ROC)", code: "2199123456", err: true},
		{name: "Invalid UEN (ROB)", code: "1234567A", err: true},
		{name: "Invalid UEN (Others)", code: "T12A1234A", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "SG",
				Code:    tt.code,
			}
			err := sg.Validate(id)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		var id *tax.Identity
		err := sg.Validate(id)
		assert.NoError(t, err)
	})
}
