package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceType(t *testing.T) {
	c := bill.InvoiceTypeCreditNote
	assert.Equal(t, cbc.Key("credit-note"), c)
	assert.Equal(t, cbc.Code("381"), bill.InvoiceTypes.UNTDID1001(c), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.InvoiceTypeCorrective
	assert.Equal(t, cbc.Code("384"), bill.InvoiceTypes.UNTDID1001(c), "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = cbc.Key("BAD_KEY")
	assert.Error(t, c.Validate())

	c = cbc.Key("foo")
	assert.Equal(t, cbc.CodeEmpty, bill.InvoiceTypes.UNTDID1001(c), "unexpected UNTDID result")

	assert.True(t, c.In("bar", "foo"))
	assert.False(t, c.In("bar", "dom"))

}
