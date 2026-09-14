package sg_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	t.Run("normalize GST number", func(t *testing.T) {
		id := &org.Identity{
			Type: sg.IdentityTypeUEN,
			Code: "SG.M91234 56 .7X",
		}
		norm.Normalize(id, tax.RegimeContext("SG"))
		assert.Equal(t, "M91234567X", id.Code.String())
	})
}

func TestValidateUENIdentity(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var id *org.Identity
		err := rules.Validate(id, tax.RegimeContext(sg.CountryCode))
		assert.NoError(t, err)
	})

	tests := []struct {
		name string
		code string
		err  bool
	}{
		{name: "UEN (ROC)", code: "199912345R", err: false},
		{name: "UEN (ROC) real", code: "199307558M", err: false}, // Defence Collective Singapore Ltd.
		{name: "UEN (ROC) real 2", code: "201189853M", err: false},
		{name: "UEN (ROB)", code: "12345678M", err: false},
		{name: "UEN (ROB) real", code: "00192200M", err: false},
		{name: "UEN (Others)", code: "T12LL1234C", err: false},
		{name: "UEN (Others) real", code: "T08GB0020K", err: false}, // IRAS
		{name: "UEN (Others) real 2", code: "S16FC0121D", err: false},
		{name: "NIRC/FIN", code: "S1234567A", err: true},
		{name: "Invalid short", code: "1234567A", err: true},
		{name: "Invalid UEN (ROC)", code: "2199123456", err: true},
		{name: "Invalid UEN (ROB)", code: "1234567A", err: true},
		{name: "Invalid UEN (Others)", code: "T12A1234A", err: true},
		{name: "Bad check char (ROC)", code: "199912345A", err: true},
		{name: "Bad check char (ROB)", code: "12345678A", err: true},
		{name: "Bad check char (Others)", code: "T12LL1234A", err: true},
		{name: "Empty code", code: "", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Type: sg.IdentityTypeUEN,
				Code: cbc.Code(tt.code),
			}
			err := rules.Validate(id, tax.RegimeContext(sg.CountryCode))
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
