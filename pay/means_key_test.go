package pay_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
)

func TestMeansKey(t *testing.T) {
	i := new(pay.Instructions)
	i.Key = cbc.Key("foo")
	err := i.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key: must be or start with a valid key")

	i.Key = pay.MeansKeyCard
	err = i.Validate()
	assert.NoError(t, err)

	// Key with an extension
	i.Key = pay.MeansKeyCard.With("visa")
	err = i.Validate()
	assert.NoError(t, err)
}
