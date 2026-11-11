package ae

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeIdentityNil(t *testing.T) {
	assert.NotPanics(t, func() { normalizeIdentity(nil) })
}
