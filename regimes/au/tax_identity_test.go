package au

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  bool
	}{
		{name: "valid ABN", code: "51824753556", err: false},
		{name: "valid ABN 2", code: "53004085616", err: false},
		{name: "invalid ABN short", code: "5182475355", err: true},
		{name: "invalid ABN long", code: "518247535566", err: true},
		{name: "invalid ABN leading zero", code: "01824753556", err: true},
		{name: "invalid ABN letters", code: "51A24753556", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "AU",
				Code:    tt.code,
			}
			err := Validate(id)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		var id *tax.Identity
		err := Validate(id)
		assert.NoError(t, err)
	})

	t.Run("non tax identity", func(t *testing.T) {
		err := Validate(struct{}{})
		assert.NoError(t, err)
	})

	t.Run("checksum validation", func(t *testing.T) {
		assert.NoError(t, validateTaxCode(cbc.Code("51824753556")))
		assert.Error(t, validateTaxCode(cbc.Code("51824753555")))
	})

	t.Run("tax code validator ignores unsupported value type", func(t *testing.T) {
		assert.NoError(t, validateTaxCode("not-a-cbc-code"))
	})

	t.Run("tax code validator ignores empty code", func(t *testing.T) {
		assert.NoError(t, validateTaxCode(cbc.Code("")))
	})

	t.Run("checksum rejects invalid length", func(t *testing.T) {
		assert.False(t, validABNChecksum("5182475355"))
	})
}
