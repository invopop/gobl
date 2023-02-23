package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceType(t *testing.T) {
	c := bill.InvoiceTypeSelfBilled
	assert.Equal(t, bill.InvoiceType("self-billed"), c)
	assert.Equal(t, cbc.Code("389"), c.UNTDID1001(), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.InvoiceTypeCorrective
	assert.Equal(t, cbc.Code("384"), c.UNTDID1001(), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.InvoiceType("foo")
	assert.Equal(t, cbc.CodeEmpty, c.UNTDID1001(), "unexpected UNTDID result")
	assert.Error(t, c.Validate())

	assert.True(t, c.In("bar", "foo"))
	assert.False(t, c.In("bar", "dom"))

	var d bill.InvoiceType
	assert.Equal(t, bill.InvoiceTypeDefault, d)

}
