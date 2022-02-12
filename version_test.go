package gobl_test

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	var v gobl.Version = "v0.10.2"
	sv := v.Semver()
	assert.EqualValues(t, 0, sv.Major())
	assert.EqualValues(t, 10, sv.Minor())
	assert.EqualValues(t, 2, sv.Patch())
}
