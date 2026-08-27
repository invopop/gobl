package mt

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

// isValidTaxIdentityCode guards against an empty code and a non-cbc.Code value.
// The rules pipeline never triggers either branch (AssertIfPresent skips empty
// values and the "code" field is always a cbc.Code), so they are exercised here
// by calling the function directly to lock in its contract.
func TestIsValidTaxIdentityCodeGuards(t *testing.T) {
	assert.False(t, isValidTaxIdentityCode(cbc.Code("")), "empty code is not valid")
	assert.False(t, isValidTaxIdentityCode("12701906"), "a plain string is not a cbc.Code")
	assert.False(t, isValidTaxIdentityCode(12701906), "a non-string value is not a cbc.Code")
	assert.True(t, isValidTaxIdentityCode(cbc.Code("12701906")), "sanity: a valid code still passes")
}
