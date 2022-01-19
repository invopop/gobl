package gobl_test

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	var v gobl.Version = "gobl.org/0.10.0"
	assert.Equal(t, "0.10.0", v.Semver())
}
