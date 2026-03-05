package tr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/tr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	r := tax.RegimeDefFor("TR")

	t.Run("normalize TCKN", func(t *testing.T) {
		id := &org.Identity{
			Type: tr.IdentityTypeTCKN,
			Code: " 125 903 265 14 ",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "12590326514", id.Code.String())
	})
}

func TestValidateTCKNIdentity(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var id *org.Identity
		err := tr.Validate(id)
		assert.NoError(t, err)
	})

	tests := []struct {
		name string
		code string
		err  string
	}{
		{name: "valid TCKN", code: "12590326514"},
		{name: "valid TCKN negative modulo", code: "19191919190"},
		{name: "too short", code: "1234567890", err: "invalid format"},
		{name: "too long", code: "123456789012", err: "invalid format"},
		{name: "non-numeric", code: "1234567890A", err: "invalid format"},
		{name: "starts with zero", code: "01234567890", err: "invalid format"},
		{name: "wrong check digit 10", code: "12345678900", err: "invalid check digit"},
		{name: "wrong check digit 11", code: "19191919191", err: "invalid check digit"},
		{name: "empty code", code: "", err: "code: cannot be blank"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Type: tr.IdentityTypeTCKN,
				Code: cbc.Code(tt.code),
			}
			err := tr.Validate(id)
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
