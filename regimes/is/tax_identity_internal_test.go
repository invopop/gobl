package is

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidKennitalaGuards covers the defensive branches that are unreachable
// through the rules pipeline (which regex-guards inputs) but exist so the
// helper remains safe when called directly.
func TestValidKennitalaGuards(t *testing.T) {
	// Digit guard: non-digit character among the first eight positions.
	assert.False(t, validKennitala("12345A7890"))
	// Century marker: valid checksum arithmetic on positions 0-8 but the
	// tenth digit is 5, which the registry never issues (only 0/8/9).
	assert.False(t, validKennitala("1201743395"))
}

// TestIsValidKennitalaCodeTypeGuard covers the defensive `return false` branch
// taken when the rules framework hands the validator a non-cbc.Code value.
func TestIsValidKennitalaCodeTypeGuard(t *testing.T) {
	assert.False(t, isValidKennitalaCode("1201743399"))
	assert.False(t, isValidKennitalaCode(nil))
}
