package sg_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  bool
	}{
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
			tID := &tax.Identity{
				Code: cbc.Code(tt.code),
			}
			err := sg.Validate(tID)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
