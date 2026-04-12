package uy_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/uy"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	r := uy.New()
	assert.NotNil(t, r)
	assert.NoError(t, r.Validate())
}
