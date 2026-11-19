package ad_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ad"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrgIdentityNormalisation(t *testing.T) {
	r := tax.RegimeDefFor("AD")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "hyphens stripped", input: "L-132950-X", expected: "L132950X"},
		{name: "lowercased", input: "l132950x", expected: "L132950X"},
		{name: "AD prefix stripped", input: "ADL132950X", expected: "L132950X"},
		{name: "NRT label prefix stripped", input: "NRT L132950X", expected: "L132950X"},
		{name: "NRT then AD prefix", input: "NRT AD L132950X", expected: "L132950X"},
		{name: "already clean", input: "L132950X", expected: "L132950X"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Type: ad.IdentityTypeNRT,
				Code: cbc.Code(tt.input),
			}
			r.NormalizeObject(id)
			assert.Equal(t, tt.expected, id.Code.String())
		})
	}
}
