package gobl_test

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	var v gobl.Version = "gobl.org/v0.10.0"
	assert.Equal(t, "v0.10.0", v.Semver())
	assert.Equal(t, "gobl.org", v.Domain())
}
