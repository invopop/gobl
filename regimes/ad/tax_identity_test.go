package ad

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  string
	}{
		{
			name: "valid NRT resident individual",
			code: "F123456A",
		},
		{
			name: "valid NRT SL",
			code: "L123456B",
		},
		{
			name: "valid NRT SA",
			code: "A123456C",
		},
		{
			name: "invalid - too short",
			code: "L12345A",
			err:  "code: must be in a valid format",
		},
		{
			name: "invalid - too long",
			code: "L1234567A",
			err:  "code: must be in a valid format",
		},
		{
			name: "invalid - missing leading letter",
			code: "1234567A",
			err:  "code: must be in a valid format",
		},
		{
			name: "invalid - missing trailing letter",
			code: "L1234567",
			err:  "code: must be in a valid format",
		},
		{
			name: "invalid - spaces",
			code: "L 123456 A",
			err:  "code: must be in a valid format",
		},
		{
			name: "invalid first letter",
			code: "X123456A",
			err:  "code: must be in a valid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: l10n.AD.Tax(),
				Code:    cbc.Code(tt.code),
			}
			err := validateTaxIdentity(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}
