package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestTypeKey(t *testing.T) {
	c := bill.TypeKeyCommercial
	assert.Equal(t, bill.TypeKey("commercial"), c)
	assert.Equal(t, org.Code("380"), c.UNTDID1001(), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.TypeKeyCorrected
	assert.Equal(t, org.Code("384"), c.UNTDID1001(), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.TypeKey("foo")
	assert.Equal(t, org.CodeEmpty, c.UNTDID1001(), "unexpected UNTDID result")
	assert.Error(t, c.Validate())

	assert.True(t, c.In("bar", "foo"))
	assert.False(t, c.In("bar", "dom"))

	var d bill.TypeKey
	assert.Equal(t, bill.TypeKeyNone, d)

}
