package gs1_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/gs1"
	"github.com/stretchr/testify/assert"
)

func TestCheckGLN(t *testing.T) {
	t.Parallel()

	valid := []cbc.Code{
		"9429041234563", // NZ prefix 94
		"5060012345679", // UK prefix 50
		"4006381333931", // DE prefix 40
	}

	invalid := []cbc.Code{
		"9429041234560",  // bad check digit
		"942904123456",   // too short
		"94290412345638", // too long
		"942904123456A",  // letters
		"",               // empty
	}

	for _, code := range valid {
		t.Run("Valid:"+code.String(), func(t *testing.T) {
			t.Parallel()
			assert.True(t, gs1.CheckGLN(code), "should be valid: %s", code)
		})
	}

	for _, code := range invalid {
		t.Run("Invalid:"+code.String(), func(t *testing.T) {
			t.Parallel()
			assert.False(t, gs1.CheckGLN(code), "should be invalid: %s", code)
		})
	}
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	assert.True(t, gs1.HasPrefix("9429041234563", "94"))
	assert.False(t, gs1.HasPrefix("9429041234563", "50"))
	assert.False(t, gs1.HasPrefix("94", "942"))
	assert.True(t, gs1.HasPrefix("5060012345670", "506"))
}
