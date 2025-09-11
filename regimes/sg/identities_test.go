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
			Key:  sg.IdentityKeyGSTNumber,
			Code: "SG.M91234 56 .7X",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "M91234567X", id.Code.String())
	})
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Key:  sg.IdentityKeyGSTNumber,
				Code: tt.code,
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
