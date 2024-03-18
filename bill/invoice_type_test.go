package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceType(t *testing.T) {
	c := bill.InvoiceTypeCreditNote
	kd := cbc.GetKeyDefinition(c, bill.InvoiceTypes)
	assert.Equal(t, cbc.Key("credit-note"), c)
	assert.Equal(t, cbc.Code("381"), kd.Map[bill.UNTDID1001Key], "unexpected UNTDID code")
	assert.NoError(t, c.Validate())

	c = bill.InvoiceTypeCorrective
	assert.NoError(t, c.Validate())

	c = cbc.Key("BAD_KEY")
	assert.Error(t, c.Validate())

	c = cbc.Key("foo")
	assert.True(t, c.In("bar", "foo"))
	assert.False(t, c.In("bar", "dom"))
}

func TestInvoiceUNTDID1001(t *testing.T) {
	inv := testInvoiceESForCorrection(t)
	assert.Equal(t, cbc.CodeEmpty, inv.UNTDID1001())
	inv.Type = bill.InvoiceTypeStandard
	assert.Equal(t, cbc.Code("380"), inv.UNTDID1001())
}
