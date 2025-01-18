package sg_test

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{name: "UEN (ROC)", code: "199912345A", expected: true},
		{name: "UEN (ROB)", code: "123456789A", expected: true},
		{name: "UEN (Others)", code: "T12AB1234A", expected: true},
		{name: "NIRC/FIN", code: "S1234567A", expected: true},
		{name: "GST", code: "AB1234567A", expected: true},
		{name: "Invalid short", code: "1234567A", expected: false},
		{name: "Invalid long", code: "A123456789", expected: false},
		{name: "Invalid UEN (ROC)", code: "2199123456", expected: false},
		{name: "Invalid UEN (ROB)", code: "12345678A", expected: false},
		{name: "Invalid UEN (Others)", code: "T12A1234A", expected: false},
		{name: "Invalid NIRC/FIN", code: "S123456A", expected: false},
		{name: "Invalid GST", code: "A1234567A", expected: false},
	}

	for _, tt := range tests {
		fmt.Println(tt.name)
		tID := &tax.Identity{
			Code: cbc.Code(tt.code),
		}
		err := sg.Validate(tID)
		if tt.expected {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
