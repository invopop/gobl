package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
)

func TestTypeKey(t *testing.T) {
	c := bill.TypeKeyCommercial
	assert.Equal(t, "", string(c))
	assert.Equal(t, "380", c.UNTDID1001(), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.TypeKeyCorrected
	assert.Equal(t, "384", c.UNTDID1001(), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.TypeKey("foo")
	assert.Equal(t, "na", c.UNTDID1001(), "unexpected UNTDID result")
	assert.Error(t, c.Validate())
}
