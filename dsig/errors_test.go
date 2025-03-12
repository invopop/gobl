package dsig_test

import (
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert.Equal(t, "key is not valid", dsig.ErrKeyInvalid.Error())
}
