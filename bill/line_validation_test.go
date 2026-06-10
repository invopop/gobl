package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestLineRequireTaxValidation(t *testing.T) {
	r := bill.RequireLineTaxCategory("VAT")
	assert.True(t, r.Check(nil))

	l := &bill.Line{}
	r2 := bill.RequireLineTaxCategory("")
	assert.True(t, r2.Check(l))

	assert.False(t, r.Check(l))

	l = &bill.Line{
		Taxes: tax.Set{
			{
				Category: "VAT",
				Percent:  num.NewPercentage(20, 3),
			},
		},
	}
	assert.True(t, r.Check(l))

	r3 := bill.RequireLineTaxCategory("IRPEF")
	assert.False(t, r3.Check(l))
}
