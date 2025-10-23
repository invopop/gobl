package sg_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	r := tax.RegimeDefFor("SG")

	t.Run("normalize GST number", func(t *testing.T) {
		id := &org.Identity{
			Type: sg.IdentityTypeUEN,
			Code: "SG.M91234 56 .7X",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "M91234567X", id.Code.String())
	})
}

func TestValidateUENIdentity(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var id *org.Identity
		err := sg.Validate(id)
		assert.NoError(t, err)
	})

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
		{name: "Empty code", code: "", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Type: sg.IdentityTypeUEN,
				Code: cbc.Code(tt.code),
			}
			err := sg.Validate(id)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
