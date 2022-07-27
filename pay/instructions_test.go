package pay_test

import (
	"testing"

	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
)

func TestInstructionsKey(t *testing.T) {
	i := new(pay.Instructions)
	i.Key = pay.MethodKey("foo")
	err := i.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key: must be a valid value")

	i.Key = pay.MethodKeyCard
	err = i.Validate()
	assert.NoError(t, err)
}
