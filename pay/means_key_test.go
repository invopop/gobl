package pay_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestMeansKey(t *testing.T) {
	i := new(pay.Instructions)
	i.Key = cbc.Key("foo")
	err := rules.Validate(i)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key must be valid")

	i.Key = pay.MeansKeyCard
	err = rules.Validate(i)
	assert.NoError(t, err)

	// Key with an extension
	i.Key = pay.MeansKeyCard.With("visa")
	err = rules.Validate(i)
	assert.NoError(t, err)
}
